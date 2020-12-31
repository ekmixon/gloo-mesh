// Code generated by skv2. DO NOT EDIT.

//go:generate mockgen -source ./multicluster_reconcilers.go -destination mocks/multicluster_reconcilers.go

// Definitions for the multicluster Kubernetes Controllers
package controller

import (
	"context"

	rbac_enterprise_mesh_gloo_solo_io_v1alpha1 "github.com/solo-io/gloo-mesh/pkg/api/rbac.enterprise.mesh.gloo.solo.io/v1alpha1"

	"github.com/pkg/errors"
	"github.com/solo-io/skv2/pkg/ezkube"
	"github.com/solo-io/skv2/pkg/multicluster"
	mc_reconcile "github.com/solo-io/skv2/pkg/multicluster/reconcile"
	"github.com/solo-io/skv2/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// Reconcile Upsert events for the Role Resource across clusters.
// implemented by the user
type MulticlusterRoleReconciler interface {
	ReconcileRole(clusterName string, obj *rbac_enterprise_mesh_gloo_solo_io_v1alpha1.Role) (reconcile.Result, error)
}

// Reconcile deletion events for the Role Resource across clusters.
// Deletion receives a reconcile.Request as we cannot guarantee the last state of the object
// before being deleted.
// implemented by the user
type MulticlusterRoleDeletionReconciler interface {
	ReconcileRoleDeletion(clusterName string, req reconcile.Request) error
}

type MulticlusterRoleReconcilerFuncs struct {
	OnReconcileRole         func(clusterName string, obj *rbac_enterprise_mesh_gloo_solo_io_v1alpha1.Role) (reconcile.Result, error)
	OnReconcileRoleDeletion func(clusterName string, req reconcile.Request) error
}

func (f *MulticlusterRoleReconcilerFuncs) ReconcileRole(clusterName string, obj *rbac_enterprise_mesh_gloo_solo_io_v1alpha1.Role) (reconcile.Result, error) {
	if f.OnReconcileRole == nil {
		return reconcile.Result{}, nil
	}
	return f.OnReconcileRole(clusterName, obj)
}

func (f *MulticlusterRoleReconcilerFuncs) ReconcileRoleDeletion(clusterName string, req reconcile.Request) error {
	if f.OnReconcileRoleDeletion == nil {
		return nil
	}
	return f.OnReconcileRoleDeletion(clusterName, req)
}

type MulticlusterRoleReconcileLoop interface {
	// AddMulticlusterRoleReconciler adds a MulticlusterRoleReconciler to the MulticlusterRoleReconcileLoop.
	AddMulticlusterRoleReconciler(ctx context.Context, rec MulticlusterRoleReconciler, predicates ...predicate.Predicate)
}

type multiclusterRoleReconcileLoop struct {
	loop multicluster.Loop
}

func (m *multiclusterRoleReconcileLoop) AddMulticlusterRoleReconciler(ctx context.Context, rec MulticlusterRoleReconciler, predicates ...predicate.Predicate) {
	genericReconciler := genericRoleMulticlusterReconciler{reconciler: rec}

	m.loop.AddReconciler(ctx, genericReconciler, predicates...)
}

func NewMulticlusterRoleReconcileLoop(name string, cw multicluster.ClusterWatcher, options reconcile.Options) MulticlusterRoleReconcileLoop {
	return &multiclusterRoleReconcileLoop{loop: mc_reconcile.NewLoop(name, cw, &rbac_enterprise_mesh_gloo_solo_io_v1alpha1.Role{}, options)}
}

type genericRoleMulticlusterReconciler struct {
	reconciler MulticlusterRoleReconciler
}

func (g genericRoleMulticlusterReconciler) ReconcileDeletion(cluster string, req reconcile.Request) error {
	if deletionReconciler, ok := g.reconciler.(MulticlusterRoleDeletionReconciler); ok {
		return deletionReconciler.ReconcileRoleDeletion(cluster, req)
	}
	return nil
}

func (g genericRoleMulticlusterReconciler) Reconcile(cluster string, object ezkube.Object) (reconcile.Result, error) {
	obj, ok := object.(*rbac_enterprise_mesh_gloo_solo_io_v1alpha1.Role)
	if !ok {
		return reconcile.Result{}, errors.Errorf("internal error: Role handler received event for %T", object)
	}
	return g.reconciler.ReconcileRole(cluster, obj)
}

// Reconcile Upsert events for the RoleBinding Resource across clusters.
// implemented by the user
type MulticlusterRoleBindingReconciler interface {
	ReconcileRoleBinding(clusterName string, obj *rbac_enterprise_mesh_gloo_solo_io_v1alpha1.RoleBinding) (reconcile.Result, error)
}

// Reconcile deletion events for the RoleBinding Resource across clusters.
// Deletion receives a reconcile.Request as we cannot guarantee the last state of the object
// before being deleted.
// implemented by the user
type MulticlusterRoleBindingDeletionReconciler interface {
	ReconcileRoleBindingDeletion(clusterName string, req reconcile.Request) error
}

type MulticlusterRoleBindingReconcilerFuncs struct {
	OnReconcileRoleBinding         func(clusterName string, obj *rbac_enterprise_mesh_gloo_solo_io_v1alpha1.RoleBinding) (reconcile.Result, error)
	OnReconcileRoleBindingDeletion func(clusterName string, req reconcile.Request) error
}

func (f *MulticlusterRoleBindingReconcilerFuncs) ReconcileRoleBinding(clusterName string, obj *rbac_enterprise_mesh_gloo_solo_io_v1alpha1.RoleBinding) (reconcile.Result, error) {
	if f.OnReconcileRoleBinding == nil {
		return reconcile.Result{}, nil
	}
	return f.OnReconcileRoleBinding(clusterName, obj)
}

func (f *MulticlusterRoleBindingReconcilerFuncs) ReconcileRoleBindingDeletion(clusterName string, req reconcile.Request) error {
	if f.OnReconcileRoleBindingDeletion == nil {
		return nil
	}
	return f.OnReconcileRoleBindingDeletion(clusterName, req)
}

type MulticlusterRoleBindingReconcileLoop interface {
	// AddMulticlusterRoleBindingReconciler adds a MulticlusterRoleBindingReconciler to the MulticlusterRoleBindingReconcileLoop.
	AddMulticlusterRoleBindingReconciler(ctx context.Context, rec MulticlusterRoleBindingReconciler, predicates ...predicate.Predicate)
}

type multiclusterRoleBindingReconcileLoop struct {
	loop multicluster.Loop
}

func (m *multiclusterRoleBindingReconcileLoop) AddMulticlusterRoleBindingReconciler(ctx context.Context, rec MulticlusterRoleBindingReconciler, predicates ...predicate.Predicate) {
	genericReconciler := genericRoleBindingMulticlusterReconciler{reconciler: rec}

	m.loop.AddReconciler(ctx, genericReconciler, predicates...)
}

func NewMulticlusterRoleBindingReconcileLoop(name string, cw multicluster.ClusterWatcher, options reconcile.Options) MulticlusterRoleBindingReconcileLoop {
	return &multiclusterRoleBindingReconcileLoop{loop: mc_reconcile.NewLoop(name, cw, &rbac_enterprise_mesh_gloo_solo_io_v1alpha1.RoleBinding{}, options)}
}

type genericRoleBindingMulticlusterReconciler struct {
	reconciler MulticlusterRoleBindingReconciler
}

func (g genericRoleBindingMulticlusterReconciler) ReconcileDeletion(cluster string, req reconcile.Request) error {
	if deletionReconciler, ok := g.reconciler.(MulticlusterRoleBindingDeletionReconciler); ok {
		return deletionReconciler.ReconcileRoleBindingDeletion(cluster, req)
	}
	return nil
}

func (g genericRoleBindingMulticlusterReconciler) Reconcile(cluster string, object ezkube.Object) (reconcile.Result, error) {
	obj, ok := object.(*rbac_enterprise_mesh_gloo_solo_io_v1alpha1.RoleBinding)
	if !ok {
		return reconcile.Result{}, errors.Errorf("internal error: RoleBinding handler received event for %T", object)
	}
	return g.reconciler.ReconcileRoleBinding(cluster, obj)
}
