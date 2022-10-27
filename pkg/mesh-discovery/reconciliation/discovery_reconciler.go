package reconciliation

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/solo-io/gloo-mesh/pkg/api/discovery.mesh.gloo.solo.io/output/discovery"
	"github.com/solo-io/gloo-mesh/pkg/common/reconciliation"

	"github.com/prometheus/client_golang/prometheus"
	settingsv1 "github.com/solo-io/gloo-mesh/pkg/api/settings.mesh.gloo.solo.io/v1"
	"sigs.k8s.io/controller-runtime/pkg/metrics"

	"github.com/solo-io/gloo-mesh/pkg/common/defaults"

	"github.com/solo-io/skv2/pkg/stats"

	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	appmeshv1beta2 "github.com/aws/aws-app-mesh-controller-for-k8s/apis/appmesh/v1beta2"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/skv2/pkg/reconcile"
	"github.com/solo-io/skv2/pkg/verifier"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo-mesh/pkg/api/discovery.mesh.gloo.solo.io/input"
	"github.com/solo-io/gloo-mesh/pkg/mesh-discovery/translation"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/skv2/contrib/pkg/output"
	"github.com/solo-io/skv2/contrib/pkg/sets"
	"github.com/solo-io/skv2/pkg/ezkube"
	"github.com/solo-io/skv2/pkg/multicluster"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// describes the states the reconciler can be in
type reconcilerState = float64

const (
	// the set of states the reconciler can be in.
	// these are used as metric values
	stateBooting reconcilerState = iota
	stateWaitingForSettings
	stateRunning
)

var (
	reconcilerStateGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gloo_mesh_discovery_reconciling_state",
		Help: "Indicates whether the Discovery reconciler is currently running (settings CR is detected). controller_namespace indicates the namespace to which the controller is deployed. A value of 0 indicates the reconciler is still initializing. A value of 1 indicates the reconciler is waiting for the Settings object to be created. A value of 2 indicates Settings object has been detected and the reconciler is running.",
	})

	recorder = reconciliation.NewRecorder(
		prometheus.CounterOpts{
			Name: "gloo_mesh_discovery_reconciles_total",
			Help: "The total number of reconciles (state resyncs) Discovery reconciler has performed.",
		},
		"detected",
		input.DiscoveryInputSnapshotGVKs,
		"discovery",
		discovery.SnapshotGVKs,
	)
)

func init() {
	metrics.Registry.MustRegister(reconcilerStateGauge)
}

type discoveryReconciler struct {
	ctx                   context.Context
	discoveryInputBuilder input.DiscoveryInputBuilder
	settingsBuilder       input.SettingsBuilder
	translator            translation.Translator
	localClient           client.Client
	history               *stats.SnapshotHistory
	verboseMode           bool
	settingsRef           *v1.ObjectRef
	verifier              verifier.ServerResourceVerifier
	totalReconciles       int
	state                 reconcilerState // track whether we are waiting for the settings object
}

func Start(
	ctx context.Context,
	agentCluster string,
	localMgr manager.Manager,
	clusters multicluster.Interface,
	mcClient multicluster.Client,
	history *stats.SnapshotHistory,
	verboseMode bool,
	settingsRef *v1.ObjectRef,
) error {

	ctx = contextutils.WithLogger(ctx, "mesh-discovery")
	ctx = reconciliation.ContextWithRecorder(ctx, recorder)

	settingsBuilder := input.NewSingleClusterSettingsBuilder(localMgr)

	var discoveryInputBuilder input.DiscoveryInputBuilder
	if clusters != nil {
		// run in master mode; I/O wired up to local and remote clusters
		discoveryInputBuilder = input.NewMultiClusterDiscoveryInputBuilder(clusters, mcClient)
	} else {
		// run in agent mode;  I/O wired up to local cluster only
		discoveryInputBuilder = input.NewSingleClusterDiscoveryInputBuilderWithClusterName(localMgr, agentCluster)

		// signal to other parts of Discovery that we are running in AGENT_MODE
		if err := os.Setenv(defaults.AgentClusterEnv, agentCluster); err != nil {
			return err
		}
	}

	verifier := verifier.NewVerifier(ctx, map[schema.GroupVersionKind]verifier.ServerVerifyOption{
		// only warn (avoids error) if appmesh Mesh resource is not available on cluster
		schema.GroupVersionKind{
			Group:   appmeshv1beta2.GroupVersion.Group,
			Version: appmeshv1beta2.GroupVersion.Version,
			Kind:    "Mesh",
		}: verifier.ServerVerifyOption_IgnoreIfNotPresent,
	})
	translator := translation.NewTranslator(translation.DefaultDependencyFactory)
	r := &discoveryReconciler{
		ctx:                   ctx,
		discoveryInputBuilder: discoveryInputBuilder,
		settingsBuilder:       settingsBuilder,
		translator:            translator,
		localClient:           localMgr.GetClient(),
		history:               history,
		verboseMode:           verboseMode,
		verifier:              verifier,
		settingsRef:           settingsRef,
	}

	discoveryPredicates := []predicate.Predicate{
		reconciliation.FilterLeaderElectionObject,
	}

	if clusters != nil {
		// running in non-relay mode; our reconciler should watch local and remote resources
		if _, err := input.RegisterInputReconciler(
			ctx,
			clusters,
			r.reconcile,
			localMgr,
			r.reconcileLocal,
			input.ReconcileOptions{
				Remote: input.RemoteReconcileOptions{
					Meshes: reconcile.Options{
						Verifier: verifier,
					},
					Predicates: discoveryPredicates,
				},
				Local:             input.LocalReconcileOptions{},
				ReconcileInterval: time.Second / 2,
			},
		); err != nil {
			return err
		}
	} else {
		// running in agent mode; our reconciler should watch only local resources
		if _, err := input.RegisterSingleClusterAgentReconciler(
			ctx,
			localMgr,
			r.reconcileLocal,
			time.Second/2,
			reconcile.Options{
				Verifier: verifier,
			},
			discoveryPredicates...,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *discoveryReconciler) reconcileLocal(obj ezkube.ResourceId) (bool, error) {
	clusterObj, ok := obj.(ezkube.ClusterResourceId)
	if !ok {
		contextutils.LoggerFrom(r.ctx).Debugf("ignoring event for non cluster type %T %v", obj, sets.Key(obj))
		return false, nil
	}
	return r.reconcile(clusterObj)
}

func (r *discoveryReconciler) reconcile(obj ezkube.ClusterResourceId) (bool, error) {
	r.totalReconciles++
	ctx := contextutils.WithLogger(r.ctx, fmt.Sprintf("reconcile-%v", r.totalReconciles))

	contextutils.LoggerFrom(ctx).Debugf("object triggered resync: %T<%v>", obj, sets.Key(obj))

	remoteInputSnap, err := r.discoveryInputBuilder.BuildSnapshot(ctx, "mesh-discovery-remote", input.DiscoveryInputBuildOptions{
		// ignore NoKindMatchError for AppMesh Mesh CRs
		// (only clusters with AppMesh Controller installed will
		// have this kind registered)
		Meshes: input.ResourceDiscoveryInputBuildOptions{
			Verifier: r.verifier,
		},
	})
	if err != nil {
		// failed to read from cache; should never happen
		return false, err
	}

	localInputSnap, err := r.settingsBuilder.BuildSnapshot(ctx, "mesh-discovery-local", input.SettingsBuildOptions{
		Settings: input.ResourceSettingsBuildOptions{
			// Ensure that only declared Settings object exists in snapshot.
			ListOptions: []client.ListOption{
				client.InNamespace(r.settingsRef.Namespace),
			},
		},
	})
	if err != nil {
		// failed to read from cache; should never happen
		return false, err
	}

	settings := r.getSettingsIfExists(ctx, localInputSnap)
	if settings == nil {
		// wait for another event to re-check for settings
		return false, nil
	}

	outputSnap, err := r.translator.Translate(
		ctx,
		remoteInputSnap,
		settings.Spec.GetDiscovery(),
	)
	if err != nil {
		// internal translator errors should never happen
		return false, err
	}

	var errs error
	errHandler := output.ErrorHandlerFuncs{
		HandleWriteErrorFunc: func(resource ezkube.Object, err error) {
			errs = multierror.Append(errs, eris.Wrapf(err, "writing resource %v failed", sets.Key(resource)))
		},
		HandleDeleteErrorFunc: func(resource ezkube.Object, err error) {
			errs = multierror.Append(errs, eris.Wrapf(err, "deleting resource %v failed", sets.Key(resource)))
		},
		HandleListErrorFunc: func(err error) {
			errs = multierror.Append(errs, eris.Wrapf(err, "listing failed"))
		},
	}
	syncOpts := output.OutputOpts{
		ErrHandler: errHandler,
	}
	outputSnap.ApplyLocalCluster(ctx, r.localClient, syncOpts)

	r.history.SetInput(remoteInputSnap)
	r.history.SetOutput(outputSnap)

	recorder.RecordReconcileResult(
		ctx,
		remoteInputSnap.Generic(),
		outputSnap.Generic(),
		errs == nil,
	)

	return false, errs
}

// if we are waiting for settings, this will log a warning and return nil, nil
// also logs and updates metrics for the reconciler state
func (r *discoveryReconciler) getSettingsIfExists(ctx context.Context, snap input.SettingsSnapshot) *settingsv1.Settings {
	settings, err := snap.Settings().Find(r.settingsRef)
	if err != nil {
		if r.state != stateWaitingForSettings {
			// settings is not ready yet, we switch to waiting state
			r.setReconcileState(ctx, stateWaitingForSettings)
		}
		return nil
	}

	if r.state != stateRunning {
		// settings is ready, we switch to running state
		r.setReconcileState(ctx, stateRunning)
	}

	return settings
}

func (r *discoveryReconciler) setReconcileState(ctx context.Context, state reconcilerState) {
	reconcilerStateGauge.Set(state)
	info := contextutils.LoggerFrom(ctx).Infof
	switch state {
	case stateBooting:
		info("discovery reconciler is booting")
	case stateWaitingForSettings:
		info("discovery reconciler is waiting for settings object %v to be received", r.settingsRef)
	case stateRunning:
		info("discovery reconciler is running using settings object %v", r.settingsRef)
	}
	r.state = state
}
