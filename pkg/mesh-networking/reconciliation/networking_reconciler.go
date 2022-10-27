package reconciliation

import (
	"context"
	"fmt"
	"sync"
	"time"

	ratelimit "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/istio/mesh/mtls"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/output/istio"
	"github.com/solo-io/gloo-mesh/pkg/common/reconciliation"
	"github.com/solo-io/skv2/pkg/resource"

	"github.com/rotisserie/eris"
	commonv1 "github.com/solo-io/gloo-mesh/pkg/api/common.mesh.gloo.solo.io/v1"
	"github.com/solo-io/skv2/pkg/stats"

	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/utils/settingsutils"
	"github.com/solo-io/skv2/contrib/pkg/output"

	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/gloo-mesh/codegen/io"
	"github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/extensions/v1beta1"
	"github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/input"
	settingsv1 "github.com/solo-io/gloo-mesh/pkg/api/settings.mesh.gloo.solo.io/v1"
	"github.com/solo-io/gloo-mesh/pkg/common/defaults"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/apply"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/extensions"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/reporting"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/utils/metautils"
	"github.com/solo-io/go-utils/contextutils"
	skinput "github.com/solo-io/skv2/contrib/pkg/input"
	"github.com/solo-io/skv2/contrib/pkg/sets"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/skv2/pkg/ezkube"
	skv2predicate "github.com/solo-io/skv2/pkg/predicate"
	"github.com/solo-io/skv2/pkg/reconcile"
	"github.com/solo-io/skv2/pkg/verifier"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
)

var recorder = reconciliation.NewRecorder(
	prometheus.CounterOpts{
		Name: "gloo_mesh_networking_reconciles_total",
		Help: "The total number of reconciles (state resyncs) Networking reconciler has performed.",
	},
	"api",
	input.LocalSnapshotGVKs,
	"istio",
	istio.SnapshotGVKs,
)

// function which defines how the Networking reconciler should be registered with internal components.
type RegisterReconcilerFunc func(
	ctx context.Context,
	cfg *rest.Config,
	reconcile skinput.SingleClusterReconcileFunc,
	reconcileOpts input.ReconcileOptions,
) (skinput.InputReconciler, error)

// function which defines how the Networking reconciler should apply its output snapshots.
type SyncOutputsFunc func(
	ctx context.Context,
	inputs input.LocalSnapshot,
	outputSnap *translation.Outputs,
	syncOpts output.OutputOpts,
) error

type networkingReconciler struct {
	ctx                        context.Context
	cfg                        *rest.Config
	localBuilder               input.LocalBuilder
	remoteBuilder              input.RemoteBuilder
	applier                    apply.Applier
	reporter                   reporting.Reporter
	translator                 translation.Translator
	syncOutputs                SyncOutputsFunc
	mgmtClient                 client.Client
	history                    *stats.SnapshotHistory
	totalReconciles            int
	verboseMode                bool
	settingsRef                *v1.ObjectRef
	extensionClients           extensions.Clientset
	reconciler                 skinput.InputReconciler
	remoteResourceVerifier     verifier.ServerResourceVerifier
	disallowIntersectingConfig bool

	// Lock for the lastSnapshot
	snapshotLock sync.RWMutex
	// This snap must be accessed using the lock to ensure no data races
	lastSnapshot input.LocalSnapshot
}

var (
	// pushNotificationId is a special identifier for a reconcile event triggered by an extension server pushing a notification
	pushNotificationId = &v1.ObjectRef{
		Name: "push-notification-event",
	}
)

func Start(
	ctx context.Context,
	localBuilder input.LocalBuilder,
	remoteBuilder input.RemoteBuilder,
	applier apply.Applier,
	reporter reporting.Reporter,
	translator translation.Translator,
	registerReconciler RegisterReconcilerFunc,
	syncOutputs SyncOutputsFunc,
	mgr manager.Manager,
	history *stats.SnapshotHistory,
	verboseMode bool,
	settingsRef *v1.ObjectRef,
	extensionClients extensions.Clientset,
	disallowIntersectingConfig bool,
	watchOutputTypes bool,
) error {
	mgmtClient := mgr.GetClient()

	ctx = contextutils.WithLogger(ctx, "networking-reconciler")
	ctx = reconciliation.ContextWithRecorder(ctx, recorder)

	remoteResourceVerifier := buildRemoteResourceVerifier(ctx)

	r := &networkingReconciler{
		ctx:                        ctx,
		cfg:                        mgr.GetConfig(),
		localBuilder:               localBuilder,
		remoteBuilder:              remoteBuilder,
		applier:                    applier,
		reporter:                   reporter,
		translator:                 translator,
		mgmtClient:                 mgmtClient,
		history:                    history,
		verboseMode:                verboseMode,
		syncOutputs:                syncOutputs,
		settingsRef:                settingsRef,
		extensionClients:           extensionClients,
		disallowIntersectingConfig: disallowIntersectingConfig,
		remoteResourceVerifier:     remoteResourceVerifier,
	}

	// watch local input types for changes
	// also watch istio output types for changes, including objects managed by Gloo Mesh itself
	// this should eventually reach a steady state since Gloo Mesh performs equality checks before updating existing objects
	remoteReconcileOptions := reconcile.Options{Verifier: remoteResourceVerifier}
	remoteReconcileOpts := input.RemoteReconcileOptions{
		IssuedCertificates:    remoteReconcileOptions,
		PodBounceDirectives:   remoteReconcileOptions,
		XdsConfigs:            remoteReconcileOptions,
		DestinationRules:      remoteReconcileOptions,
		EnvoyFilters:          remoteReconcileOptions,
		Gateways:              remoteReconcileOptions,
		ServiceEntries:        remoteReconcileOptions,
		VirtualServices:       remoteReconcileOptions,
		AuthorizationPolicies: remoteReconcileOptions,
		Sidecars:              remoteReconcileOptions,
		RateLimitConfigs:      remoteReconcileOptions,
		Predicates: []predicate.Predicate{
			skv2predicate.SimplePredicate{
				Filter: skv2predicate.SimpleEventFilterFunc(isIgnoredConfigMap),
			},
		},
	}
	// ignore all events (i.e. don't reconcile) if not watching output types
	if !watchOutputTypes {
		remoteReconcileOpts.Predicates = append(
			remoteReconcileOpts.Predicates,
			skv2predicate.SimplePredicate{
				Filter: skv2predicate.SimpleEventFilterFunc(
					func(obj client.Object) bool {
						return true
					},
				),
			},
		)
	}
	contextutils.LoggerFrom(ctx).Debugw("starting input watches", "watch_gvks", input.LocalSnapshotGVKs)
	reconciler, err := registerReconciler(
		ctx,
		mgr.GetConfig(),
		r.reconcile,
		input.ReconcileOptions{
			Local: input.LocalReconcileOptions{
				Predicates: []predicate.Predicate{
					skv2predicate.SimplePredicate{
						Filter: skv2predicate.SimpleEventFilterFunc(r.isIgnoredSecret),
					},
				},
			},
			Remote:            remoteReconcileOpts,
			ReconcileInterval: time.Second / 2,
		},
	)
	if err != nil {
		return err
	}

	r.reconciler = reconciler

	return nil
}

// reconcile global state
func (r *networkingReconciler) reconcile(obj ezkube.ResourceId) (bool, error) {
	r.totalReconciles++
	ctx := contextutils.WithLogger(r.ctx, fmt.Sprintf("reconcile-%v", r.totalReconciles))

	contextutils.LoggerFrom(ctx).Debugf("object triggered resync: %T<%v>", obj, sets.Key(obj))

	// build the input snapshot from the caches
	inputSnap, err := r.localBuilder.BuildSnapshot(ctx, "mesh-networking", input.LocalBuildOptions{
		// only look at kube clusters in our own namespace
		KubernetesClusters: input.ResourceLocalBuildOptions{
			ListOptions: []client.ListOption{client.InNamespace(defaults.GetPodNamespace())},
		},
	})
	if err != nil {
		// failed to read from cache; should never happen
		return false, eris.Wrapf(err, "failed to build api snapshot from cache")
	}

	// Set last input snapshot for predicates
	r.snapshotLock.Lock()
	r.lastSnapshot = inputSnap
	r.snapshotLock.Unlock()

	if err := r.syncSettings(&ctx, inputSnap); err != nil {
		// fail early if settings failed to sync
		return false, eris.Wrapf(err, "failed to sync settings")
	}

	// nil istioInputSnap signals to downstream translators that intersecting config should not be detected
	var userSupplied input.RemoteSnapshot
	if r.disallowIntersectingConfig {
		selector := labels.NewSelector()
		for k := range metautils.TranslatedObjectLabels() {
			// select objects without the translated object label key
			requirement, err := labels.NewRequirement(k, selection.DoesNotExist, nil)
			if err != nil {
				// shouldn't happen
				return false, err
			}
			selector = selector.Add([]labels.Requirement{*requirement}...)
		}
		resourceBuildOptions := input.ResourceRemoteBuildOptions{
			ListOptions: []client.ListOption{
				&client.ListOptions{LabelSelector: selector},
			},
			Verifier: r.remoteResourceVerifier,
		}
		userSupplied, err = r.remoteBuilder.BuildSnapshot(ctx, "mesh-networking-istio-inputs", input.RemoteBuildOptions{
			IssuedCertificates:    resourceBuildOptions,
			PodBounceDirectives:   resourceBuildOptions,
			XdsConfigs:            resourceBuildOptions,
			DestinationRules:      resourceBuildOptions,
			EnvoyFilters:          resourceBuildOptions,
			Gateways:              resourceBuildOptions,
			ServiceEntries:        resourceBuildOptions,
			VirtualServices:       resourceBuildOptions,
			AuthorizationPolicies: resourceBuildOptions,
			Sidecars:              resourceBuildOptions,
			RateLimitConfigs:      resourceBuildOptions,
		})
		if err != nil {
			// failed to read from cache; should never happen
			return false, eris.Wrapf(err, "failed to build user snapshot from cache")
		}
	}

	// apply policies to the discovery resources they target
	r.applier.Apply(ctx, inputSnap, userSupplied)

	// append errors as we still want to sync statuses if applying translation fails
	var errs error

	// translate and apply outputs
	outputs, err := r.translateAndSyncOutputs(ctx, inputSnap, userSupplied)
	if err != nil {
		errs = multierror.Append(errs, eris.Wrap(err, "translation error"))
	}

	contextutils.LoggerFrom(ctx).Debugf("syncing input object statuses")
	// update statuses of input objects
	if err := inputSnap.SyncStatuses(ctx, r.mgmtClient, input.LocalSyncStatusOptions{
		// keep this list up to date with all networking status outputs
		Settings:                true,
		Destination:             true,
		Workload:                true,
		Mesh:                    true,
		TrafficPolicy:           true,
		AccessPolicy:            true,
		VirtualMesh:             true,
		WasmDeployment:          true,
		AccessLogRecord:         true,
		VirtualDestination:      true,
		VirtualGateway:          true,
		VirtualHost:             true,
		RouteTable:              true,
		ServiceDependency:       true,
		CertificateVerification: true,
	}); err != nil {
		errs = multierror.Append(errs, eris.Wrap(err, "updating input object statuses"))
	}

	var outputSnap resource.ClusterSnapshot
	if outputs != nil && outputs.Istio != nil {
		outputSnap = outputs.Istio.Generic()
	}
	recorder.RecordReconcileResult(
		ctx,
		inputSnap.Generic(),
		outputSnap,
		errs == nil,
	)

	return false, errs
}

// returns true if the passed object is a secret which is ignored by gloo mesh
func (r *networkingReconciler) isIgnoredSecret(obj client.Object) bool {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return false
	}

	r.snapshotLock.RLock()
	defer r.snapshotLock.RUnlock()
	if r.lastSnapshot == nil {
		return true
	}

	// Search virtual meshes for secret reference
	for _, v := range r.lastSnapshot.VirtualMeshes().List() {
		rootCaSecret := v.Spec.MtlsConfig.GetShared().GetRootCertificateAuthority().GetSecret()
		// If the secret reference is nil we can break out.
		if rootCaSecret == nil {
			continue
		}
		// If the secret is being referenced, we should handle this event
		if rootCaSecret.GetName() == secret.Name &&
			rootCaSecret.GetNamespace() == secret.Namespace {
			return false
		}
	}
	// Check if generated secret type
	return !mtls.IsSigningCert(secret)
}

func (r *networkingReconciler) translateAndSyncOutputs(ctx context.Context, in input.LocalSnapshot, userSupplied input.RemoteSnapshot) (*translation.Outputs, error) {

	outputSnap, err := r.translator.Translate(ctx, in, userSupplied, r.reporter)
	if err != nil {
		// internal translator errors should never happen
		return nil, err
	}

	r.history.SetInput(in)
	r.history.SetOutput(outputSnap)

	contextutils.LoggerFrom(ctx).Debugf("syncing outputs")
	errHandler := newErrHandler(ctx, in)
	verifier := verifier.NewOutputVerifier(ctx, r.cfg, map[schema.GroupVersionKind]verifier.ServerVerifyOption{
		// ignore if ratelimitconfig resource is not available on cluster
		ratelimit.RateLimitConfigGVK: verifier.ServerVerifyOption_IgnoreIfNotPresent,
	})
	syncOpts := output.OutputOpts{
		Verifier:   verifier,
		ErrHandler: errHandler,
	}
	if err := r.syncOutputs(ctx, in, outputSnap, syncOpts); err != nil {
		return nil, multierror.Append(err, errHandler.Errors())
	}

	return outputSnap, errHandler.Errors()
}

// stores settings inside the context and initiates connections to extension servers.
// processing/validation errors will be reported to the settings status.
func (r *networkingReconciler) syncSettings(ctx *context.Context, in input.LocalSnapshot) error {
	settings, err := in.Settings().Find(r.settingsRef)
	if err != nil {
		return eris.Wrapf(err, "settings object does not exist")
	}

	*ctx = settingsutils.ContextWithSettings(*ctx, settings)

	settings.Status = settingsv1.SettingsStatus{
		ObservedGeneration: settings.Generation,
		State:              commonv1.ApprovalState_ACCEPTED,
	}

	// update configured NetworkExtensionServers for the extension clients which are called inside the translator.
	if err := r.extensionClients.ConfigureServers(settings.Spec.NetworkingExtensionServers, func(_ *v1beta1.PushNotification) {
		// ignore error because underlying impl should never error here
		_, _ = r.reconciler.ReconcileLocalGeneric(pushNotificationId)
	}); err != nil {
		return eris.Wrap(err, "failed to configure networking extensions clients")
	}

	return nil
}

// returns true if the passed object is a configmap which is of a type that is ignored by GlooMesh
// this is necessary because Istio-controlled configmaps update very frequently
func isIgnoredConfigMap(obj client.Object) bool {
	_, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return false
	}
	return !metautils.IsTranslated(obj)
}

// build a verifier that ignores NoKindMatch errors for mesh-specific types
// we expect these errors on clusters on which that mesh is not deployed
func buildRemoteResourceVerifier(ctx context.Context) verifier.ServerResourceVerifier {
	ctx = contextutils.WithLogger(ctx, "resource-verifier")
	options := map[schema.GroupVersionKind]verifier.ServerVerifyOption{}
	for groupVersion, kinds := range io.IstioNetworkingOutputTypes.Snapshot {
		for _, kind := range kinds {
			gvk := schema.GroupVersionKind{
				Group:   groupVersion.Group,
				Version: groupVersion.Version,
				Kind:    kind,
			}
			options[gvk] = verifier.ServerVerifyOption_IgnoreIfNotPresent
		}
	}
	return verifier.NewVerifier(ctx, options)
}
