// Code generated by skv2. DO NOT EDIT.

package input

import (
	"context"
	"time"

	"github.com/solo-io/skv2/contrib/pkg/input"
	sk_core_v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/skv2/pkg/multicluster"
	"github.com/solo-io/skv2/pkg/reconcile"

	appmesh_k8s_aws_v1beta2 "github.com/aws/aws-app-mesh-controller-for-k8s/apis/appmesh/v1beta2"
	appmesh_k8s_aws_v1beta2_controllers "github.com/solo-io/external-apis/pkg/api/appmesh/appmesh.k8s.aws/v1beta2/controller"
	apps_v1_controllers "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1/controller"
	v1_controllers "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/controller"
	settings_mesh_gloo_solo_io_v1alpha2 "github.com/solo-io/gloo-mesh/pkg/api/settings.mesh.gloo.solo.io/v1alpha2"
	settings_mesh_gloo_solo_io_v1alpha2_controllers "github.com/solo-io/gloo-mesh/pkg/api/settings.mesh.gloo.solo.io/v1alpha2/controller"
	apps_v1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// The Input Reconciler calls a simple func(id) error whenever a
// storage event is received for any of:
// * Meshes
// * VirtualNodes
// * ConfigMaps
// * Services
// * Pods
// * Nodes
// * Deployments
// * ReplicaSets
// * DaemonSets
// * StatefulSets
// from a remote cluster.
// * Settings
// from the local cluster.

type ReconcileOptions struct {
	Remote RemoteReconcileOptions
	Local  LocalReconcileOptions

	// the ReconcileInterval, if greater than 0, will limit the number of reconciles
	// to one per interval.
	ReconcileInterval time.Duration
}

// register the given multi cluster reconcile func with the cluster watcher
// register the given single cluster reconcile func with the local manager
func RegisterInputReconciler(
	ctx context.Context,
	clusters multicluster.ClusterWatcher,
	multiClusterReconcileFunc input.MultiClusterReconcileFunc,
	mgr manager.Manager,
	singleClusterReconcileFunc input.SingleClusterReconcileFunc,
	options ReconcileOptions,
) (input.InputReconciler, error) {
	// [appmesh.k8s.aws/v1beta2 v1 apps/v1] false 3
	// [settings.mesh.gloo.solo.io/v1alpha2]

	base := input.NewInputReconciler(
		ctx,
		multiClusterReconcileFunc,
		singleClusterReconcileFunc,
		options.ReconcileInterval,
	)

	// initialize reconcile loops

	// initialize Meshes reconcile loop for remote clusters
	appmesh_k8s_aws_v1beta2_controllers.NewMulticlusterMeshReconcileLoop("Mesh", clusters, options.Remote.Meshes).AddMulticlusterMeshReconciler(ctx, &remoteInputReconciler{base: base}, options.Remote.Predicates...)
	// initialize VirtualNodes reconcile loop for remote clusters
	appmesh_k8s_aws_v1beta2_controllers.NewMulticlusterVirtualNodeReconcileLoop("VirtualNode", clusters, options.Remote.VirtualNodes).AddMulticlusterVirtualNodeReconciler(ctx, &remoteInputReconciler{base: base}, options.Remote.Predicates...)

	// initialize ConfigMaps reconcile loop for remote clusters
	v1_controllers.NewMulticlusterConfigMapReconcileLoop("ConfigMap", clusters, options.Remote.ConfigMaps).AddMulticlusterConfigMapReconciler(ctx, &remoteInputReconciler{base: base}, options.Remote.Predicates...)
	// initialize Services reconcile loop for remote clusters
	v1_controllers.NewMulticlusterServiceReconcileLoop("Service", clusters, options.Remote.Services).AddMulticlusterServiceReconciler(ctx, &remoteInputReconciler{base: base}, options.Remote.Predicates...)
	// initialize Pods reconcile loop for remote clusters
	v1_controllers.NewMulticlusterPodReconcileLoop("Pod", clusters, options.Remote.Pods).AddMulticlusterPodReconciler(ctx, &remoteInputReconciler{base: base}, options.Remote.Predicates...)
	// initialize Nodes reconcile loop for remote clusters
	v1_controllers.NewMulticlusterNodeReconcileLoop("Node", clusters, options.Remote.Nodes).AddMulticlusterNodeReconciler(ctx, &remoteInputReconciler{base: base}, options.Remote.Predicates...)

	// initialize Deployments reconcile loop for remote clusters
	apps_v1_controllers.NewMulticlusterDeploymentReconcileLoop("Deployment", clusters, options.Remote.Deployments).AddMulticlusterDeploymentReconciler(ctx, &remoteInputReconciler{base: base}, options.Remote.Predicates...)
	// initialize ReplicaSets reconcile loop for remote clusters
	apps_v1_controllers.NewMulticlusterReplicaSetReconcileLoop("ReplicaSet", clusters, options.Remote.ReplicaSets).AddMulticlusterReplicaSetReconciler(ctx, &remoteInputReconciler{base: base}, options.Remote.Predicates...)
	// initialize DaemonSets reconcile loop for remote clusters
	apps_v1_controllers.NewMulticlusterDaemonSetReconcileLoop("DaemonSet", clusters, options.Remote.DaemonSets).AddMulticlusterDaemonSetReconciler(ctx, &remoteInputReconciler{base: base}, options.Remote.Predicates...)
	// initialize StatefulSets reconcile loop for remote clusters
	apps_v1_controllers.NewMulticlusterStatefulSetReconcileLoop("StatefulSet", clusters, options.Remote.StatefulSets).AddMulticlusterStatefulSetReconciler(ctx, &remoteInputReconciler{base: base}, options.Remote.Predicates...)

	// initialize Settings reconcile loop for local cluster
	if err := settings_mesh_gloo_solo_io_v1alpha2_controllers.NewSettingsReconcileLoop("Settings", mgr, options.Local.Settings).RunSettingsReconciler(ctx, &localInputReconciler{base: base}, options.Local.Predicates...); err != nil {
		return nil, err
	}

	return base, nil
}

// Options for reconciling a snapshot in remote clusters
type RemoteReconcileOptions struct {

	// Options for reconciling Meshes
	Meshes reconcile.Options
	// Options for reconciling VirtualNodes
	VirtualNodes reconcile.Options

	// Options for reconciling ConfigMaps
	ConfigMaps reconcile.Options
	// Options for reconciling Services
	Services reconcile.Options
	// Options for reconciling Pods
	Pods reconcile.Options
	// Options for reconciling Nodes
	Nodes reconcile.Options

	// Options for reconciling Deployments
	Deployments reconcile.Options
	// Options for reconciling ReplicaSets
	ReplicaSets reconcile.Options
	// Options for reconciling DaemonSets
	DaemonSets reconcile.Options
	// Options for reconciling StatefulSets
	StatefulSets reconcile.Options

	// optional predicates for filtering remote events
	Predicates []predicate.Predicate
}

type remoteInputReconciler struct {
	base input.InputReconciler
}

func (r *remoteInputReconciler) ReconcileMesh(clusterName string, obj *appmesh_k8s_aws_v1beta2.Mesh) (reconcile.Result, error) {
	obj.ClusterName = clusterName
	return r.base.ReconcileRemoteGeneric(obj)
}

func (r *remoteInputReconciler) ReconcileMeshDeletion(clusterName string, obj reconcile.Request) error {
	ref := &sk_core_v1.ClusterObjectRef{
		Name:        obj.Name,
		Namespace:   obj.Namespace,
		ClusterName: clusterName,
	}
	_, err := r.base.ReconcileRemoteGeneric(ref)
	return err
}

func (r *remoteInputReconciler) ReconcileVirtualNode(clusterName string, obj *appmesh_k8s_aws_v1beta2.VirtualNode) (reconcile.Result, error) {
	obj.ClusterName = clusterName
	return r.base.ReconcileRemoteGeneric(obj)
}

func (r *remoteInputReconciler) ReconcileVirtualNodeDeletion(clusterName string, obj reconcile.Request) error {
	ref := &sk_core_v1.ClusterObjectRef{
		Name:        obj.Name,
		Namespace:   obj.Namespace,
		ClusterName: clusterName,
	}
	_, err := r.base.ReconcileRemoteGeneric(ref)
	return err
}

func (r *remoteInputReconciler) ReconcileConfigMap(clusterName string, obj *v1.ConfigMap) (reconcile.Result, error) {
	obj.ClusterName = clusterName
	return r.base.ReconcileRemoteGeneric(obj)
}

func (r *remoteInputReconciler) ReconcileConfigMapDeletion(clusterName string, obj reconcile.Request) error {
	ref := &sk_core_v1.ClusterObjectRef{
		Name:        obj.Name,
		Namespace:   obj.Namespace,
		ClusterName: clusterName,
	}
	_, err := r.base.ReconcileRemoteGeneric(ref)
	return err
}

func (r *remoteInputReconciler) ReconcileService(clusterName string, obj *v1.Service) (reconcile.Result, error) {
	obj.ClusterName = clusterName
	return r.base.ReconcileRemoteGeneric(obj)
}

func (r *remoteInputReconciler) ReconcileServiceDeletion(clusterName string, obj reconcile.Request) error {
	ref := &sk_core_v1.ClusterObjectRef{
		Name:        obj.Name,
		Namespace:   obj.Namespace,
		ClusterName: clusterName,
	}
	_, err := r.base.ReconcileRemoteGeneric(ref)
	return err
}

func (r *remoteInputReconciler) ReconcilePod(clusterName string, obj *v1.Pod) (reconcile.Result, error) {
	obj.ClusterName = clusterName
	return r.base.ReconcileRemoteGeneric(obj)
}

func (r *remoteInputReconciler) ReconcilePodDeletion(clusterName string, obj reconcile.Request) error {
	ref := &sk_core_v1.ClusterObjectRef{
		Name:        obj.Name,
		Namespace:   obj.Namespace,
		ClusterName: clusterName,
	}
	_, err := r.base.ReconcileRemoteGeneric(ref)
	return err
}

func (r *remoteInputReconciler) ReconcileNode(clusterName string, obj *v1.Node) (reconcile.Result, error) {
	obj.ClusterName = clusterName
	return r.base.ReconcileRemoteGeneric(obj)
}

func (r *remoteInputReconciler) ReconcileNodeDeletion(clusterName string, obj reconcile.Request) error {
	ref := &sk_core_v1.ClusterObjectRef{
		Name:        obj.Name,
		Namespace:   obj.Namespace,
		ClusterName: clusterName,
	}
	_, err := r.base.ReconcileRemoteGeneric(ref)
	return err
}

func (r *remoteInputReconciler) ReconcileDeployment(clusterName string, obj *apps_v1.Deployment) (reconcile.Result, error) {
	obj.ClusterName = clusterName
	return r.base.ReconcileRemoteGeneric(obj)
}

func (r *remoteInputReconciler) ReconcileDeploymentDeletion(clusterName string, obj reconcile.Request) error {
	ref := &sk_core_v1.ClusterObjectRef{
		Name:        obj.Name,
		Namespace:   obj.Namespace,
		ClusterName: clusterName,
	}
	_, err := r.base.ReconcileRemoteGeneric(ref)
	return err
}

func (r *remoteInputReconciler) ReconcileReplicaSet(clusterName string, obj *apps_v1.ReplicaSet) (reconcile.Result, error) {
	obj.ClusterName = clusterName
	return r.base.ReconcileRemoteGeneric(obj)
}

func (r *remoteInputReconciler) ReconcileReplicaSetDeletion(clusterName string, obj reconcile.Request) error {
	ref := &sk_core_v1.ClusterObjectRef{
		Name:        obj.Name,
		Namespace:   obj.Namespace,
		ClusterName: clusterName,
	}
	_, err := r.base.ReconcileRemoteGeneric(ref)
	return err
}

func (r *remoteInputReconciler) ReconcileDaemonSet(clusterName string, obj *apps_v1.DaemonSet) (reconcile.Result, error) {
	obj.ClusterName = clusterName
	return r.base.ReconcileRemoteGeneric(obj)
}

func (r *remoteInputReconciler) ReconcileDaemonSetDeletion(clusterName string, obj reconcile.Request) error {
	ref := &sk_core_v1.ClusterObjectRef{
		Name:        obj.Name,
		Namespace:   obj.Namespace,
		ClusterName: clusterName,
	}
	_, err := r.base.ReconcileRemoteGeneric(ref)
	return err
}

func (r *remoteInputReconciler) ReconcileStatefulSet(clusterName string, obj *apps_v1.StatefulSet) (reconcile.Result, error) {
	obj.ClusterName = clusterName
	return r.base.ReconcileRemoteGeneric(obj)
}

func (r *remoteInputReconciler) ReconcileStatefulSetDeletion(clusterName string, obj reconcile.Request) error {
	ref := &sk_core_v1.ClusterObjectRef{
		Name:        obj.Name,
		Namespace:   obj.Namespace,
		ClusterName: clusterName,
	}
	_, err := r.base.ReconcileRemoteGeneric(ref)
	return err
}

// Options for reconciling a snapshot in remote clusters
type LocalReconcileOptions struct {

	// Options for reconciling Settings
	Settings reconcile.Options

	// optional predicates for filtering local events
	Predicates []predicate.Predicate
}

type localInputReconciler struct {
	base input.InputReconciler
}

func (r *localInputReconciler) ReconcileSettings(obj *settings_mesh_gloo_solo_io_v1alpha2.Settings) (reconcile.Result, error) {
	return r.base.ReconcileLocalGeneric(obj)
}

func (r *localInputReconciler) ReconcileSettingsDeletion(obj reconcile.Request) error {
	ref := &sk_core_v1.ObjectRef{
		Name:      obj.Name,
		Namespace: obj.Namespace,
	}
	_, err := r.base.ReconcileLocalGeneric(ref)
	return err
}
