// Code generated by skv2. DO NOT EDIT.

//go:generate mockgen -source ./snapshot.go -destination mocks/snapshot.go

// The Input Snapshot contains the set of all:
// * IssuedCertificates
// * CertificateRequests
// * PodBounceDirectives
// * Secrets
// * Pods
// read from a given cluster or set of clusters, across all namespaces.
//
// A snapshot can be constructed from either a single Manager (for a single cluster)
// or a ClusterWatcher (for multiple clusters) using the SnapshotBuilder.
//
// Resources in a MultiCluster snapshot will have their ClusterName set to the
// name of the cluster from which the resource was read.

package input

import (
	"context"
	"encoding/json"

	"github.com/solo-io/skv2/pkg/verifier"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/hashicorp/go-multierror"

	"github.com/solo-io/skv2/pkg/controllerutils"
	"github.com/solo-io/skv2/pkg/multicluster"
	"sigs.k8s.io/controller-runtime/pkg/client"

	certificates_mesh_gloo_solo_io_v1alpha2 "github.com/solo-io/gloo-mesh/pkg/api/certificates.mesh.gloo.solo.io/v1alpha2"
	certificates_mesh_gloo_solo_io_v1alpha2_sets "github.com/solo-io/gloo-mesh/pkg/api/certificates.mesh.gloo.solo.io/v1alpha2/sets"

	v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	v1_sets "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/sets"
)

// the snapshot of input resources consumed by translation
type Snapshot interface {

	// return the set of input IssuedCertificates
	IssuedCertificates() certificates_mesh_gloo_solo_io_v1alpha2_sets.IssuedCertificateSet
	// return the set of input CertificateRequests
	CertificateRequests() certificates_mesh_gloo_solo_io_v1alpha2_sets.CertificateRequestSet
	// return the set of input PodBounceDirectives
	PodBounceDirectives() certificates_mesh_gloo_solo_io_v1alpha2_sets.PodBounceDirectiveSet

	// return the set of input Secrets
	Secrets() v1_sets.SecretSet
	// return the set of input Pods
	Pods() v1_sets.PodSet
	// update the status of all input objects which support
	// the Status subresource (in the local cluster)
	SyncStatuses(ctx context.Context, c client.Client) error

	// update the status of all input objects which support
	// the Status subresource (across multiple clusters)
	SyncStatusesMultiCluster(ctx context.Context, mcClient multicluster.Client) error
	// serialize the entire snapshot as JSON
	MarshalJSON() ([]byte, error)
}

type snapshot struct {
	name string

	issuedCertificates  certificates_mesh_gloo_solo_io_v1alpha2_sets.IssuedCertificateSet
	certificateRequests certificates_mesh_gloo_solo_io_v1alpha2_sets.CertificateRequestSet
	podBounceDirectives certificates_mesh_gloo_solo_io_v1alpha2_sets.PodBounceDirectiveSet

	secrets v1_sets.SecretSet
	pods    v1_sets.PodSet
}

func NewSnapshot(
	name string,

	issuedCertificates certificates_mesh_gloo_solo_io_v1alpha2_sets.IssuedCertificateSet,
	certificateRequests certificates_mesh_gloo_solo_io_v1alpha2_sets.CertificateRequestSet,
	podBounceDirectives certificates_mesh_gloo_solo_io_v1alpha2_sets.PodBounceDirectiveSet,

	secrets v1_sets.SecretSet,
	pods v1_sets.PodSet,

) Snapshot {
	return &snapshot{
		name: name,

		issuedCertificates:  issuedCertificates,
		certificateRequests: certificateRequests,
		podBounceDirectives: podBounceDirectives,
		secrets:             secrets,
		pods:                pods,
	}
}

func (s snapshot) IssuedCertificates() certificates_mesh_gloo_solo_io_v1alpha2_sets.IssuedCertificateSet {
	return s.issuedCertificates
}

func (s snapshot) CertificateRequests() certificates_mesh_gloo_solo_io_v1alpha2_sets.CertificateRequestSet {
	return s.certificateRequests
}

func (s snapshot) PodBounceDirectives() certificates_mesh_gloo_solo_io_v1alpha2_sets.PodBounceDirectiveSet {
	return s.podBounceDirectives
}

func (s snapshot) Secrets() v1_sets.SecretSet {
	return s.secrets
}

func (s snapshot) Pods() v1_sets.PodSet {
	return s.pods
}
func (s snapshot) SyncStatuses(ctx context.Context, c client.Client) error {

	for _, obj := range s.IssuedCertificates().List() {
		if _, err := controllerutils.UpdateStatus(ctx, c, obj); err != nil {
			return err
		}
	}
	for _, obj := range s.CertificateRequests().List() {
		if _, err := controllerutils.UpdateStatus(ctx, c, obj); err != nil {
			return err
		}
	}

	return nil
}

func (s snapshot) SyncStatusesMultiCluster(ctx context.Context, mcClient multicluster.Client) error {

	for _, obj := range s.IssuedCertificates().List() {
		clusterClient, err := mcClient.Cluster(obj.ClusterName)
		if err != nil {
			return err
		}
		if _, err := controllerutils.UpdateStatus(ctx, clusterClient, obj); err != nil {
			return err
		}
	}
	for _, obj := range s.CertificateRequests().List() {
		clusterClient, err := mcClient.Cluster(obj.ClusterName)
		if err != nil {
			return err
		}
		if _, err := controllerutils.UpdateStatus(ctx, clusterClient, obj); err != nil {
			return err
		}
	}

	return nil
}

func (s snapshot) MarshalJSON() ([]byte, error) {
	snapshotMap := map[string]interface{}{"name": s.name}

	snapshotMap["issuedCertificates"] = s.issuedCertificates.List()
	snapshotMap["certificateRequests"] = s.certificateRequests.List()
	snapshotMap["podBounceDirectives"] = s.podBounceDirectives.List()
	snapshotMap["secrets"] = s.secrets.List()
	snapshotMap["pods"] = s.pods.List()
	return json.Marshal(snapshotMap)
}

// builds the input snapshot from API Clients.
// Two types of builders are available:
// a builder for snapshots of resources across multiple clusters
// a builder for snapshots of resources within a single cluster
type Builder interface {
	BuildSnapshot(ctx context.Context, name string, opts BuildOptions) (Snapshot, error)
}

// Options for building a snapshot
type BuildOptions struct {

	// List options for composing a snapshot from IssuedCertificates
	IssuedCertificates ResourceBuildOptions
	// List options for composing a snapshot from CertificateRequests
	CertificateRequests ResourceBuildOptions
	// List options for composing a snapshot from PodBounceDirectives
	PodBounceDirectives ResourceBuildOptions

	// List options for composing a snapshot from Secrets
	Secrets ResourceBuildOptions
	// List options for composing a snapshot from Pods
	Pods ResourceBuildOptions
}

// Options for reading resources of a given type
type ResourceBuildOptions struct {

	// List options for composing a snapshot from a resource type
	ListOptions []client.ListOption

	// If provided, ensure the resource has been verified before adding it to snapshots
	Verifier verifier.ServerResourceVerifier
}

// build a snapshot from resources across multiple clusters
type multiClusterBuilder struct {
	clusters multicluster.Interface
	client   multicluster.Client
}

// Produces snapshots of resources across all clusters defined in the ClusterSet
func NewMultiClusterBuilder(
	clusters multicluster.Interface,
	client multicluster.Client,
) Builder {
	return &multiClusterBuilder{
		clusters: clusters,
		client:   client,
	}
}

func (b *multiClusterBuilder) BuildSnapshot(ctx context.Context, name string, opts BuildOptions) (Snapshot, error) {

	issuedCertificates := certificates_mesh_gloo_solo_io_v1alpha2_sets.NewIssuedCertificateSet()
	certificateRequests := certificates_mesh_gloo_solo_io_v1alpha2_sets.NewCertificateRequestSet()
	podBounceDirectives := certificates_mesh_gloo_solo_io_v1alpha2_sets.NewPodBounceDirectiveSet()

	secrets := v1_sets.NewSecretSet()
	pods := v1_sets.NewPodSet()

	var errs error

	for _, cluster := range b.clusters.ListClusters() {

		if err := b.insertIssuedCertificatesFromCluster(ctx, cluster, issuedCertificates, opts.IssuedCertificates); err != nil {
			errs = multierror.Append(errs, err)
		}
		if err := b.insertCertificateRequestsFromCluster(ctx, cluster, certificateRequests, opts.CertificateRequests); err != nil {
			errs = multierror.Append(errs, err)
		}
		if err := b.insertPodBounceDirectivesFromCluster(ctx, cluster, podBounceDirectives, opts.PodBounceDirectives); err != nil {
			errs = multierror.Append(errs, err)
		}
		if err := b.insertSecretsFromCluster(ctx, cluster, secrets, opts.Secrets); err != nil {
			errs = multierror.Append(errs, err)
		}
		if err := b.insertPodsFromCluster(ctx, cluster, pods, opts.Pods); err != nil {
			errs = multierror.Append(errs, err)
		}

	}

	outputSnap := NewSnapshot(
		name,

		issuedCertificates,
		certificateRequests,
		podBounceDirectives,
		secrets,
		pods,
	)

	return outputSnap, errs
}

func (b *multiClusterBuilder) insertIssuedCertificatesFromCluster(ctx context.Context, cluster string, issuedCertificates certificates_mesh_gloo_solo_io_v1alpha2_sets.IssuedCertificateSet, opts ResourceBuildOptions) error {
	issuedCertificateClient, err := certificates_mesh_gloo_solo_io_v1alpha2.NewMulticlusterIssuedCertificateClient(b.client).Cluster(cluster)
	if err != nil {
		return err
	}

	if opts.Verifier != nil {
		mgr, err := b.clusters.Cluster(cluster)
		if err != nil {
			return err
		}

		gvk := schema.GroupVersionKind{
			Group:   "certificates.mesh.gloo.solo.io",
			Version: "v1alpha2",
			Kind:    "IssuedCertificate",
		}

		if resourceRegistered, err := opts.Verifier.VerifyServerResource(
			cluster,
			mgr.GetConfig(),
			gvk,
		); err != nil {
			return err
		} else if !resourceRegistered {
			return nil
		}
	}

	issuedCertificateList, err := issuedCertificateClient.ListIssuedCertificate(ctx, opts.ListOptions...)
	if err != nil {
		return err
	}

	for _, item := range issuedCertificateList.Items {
		item := item               // pike
		item.ClusterName = cluster // set cluster for in-memory processing
		issuedCertificates.Insert(&item)
	}

	return nil
}
func (b *multiClusterBuilder) insertCertificateRequestsFromCluster(ctx context.Context, cluster string, certificateRequests certificates_mesh_gloo_solo_io_v1alpha2_sets.CertificateRequestSet, opts ResourceBuildOptions) error {
	certificateRequestClient, err := certificates_mesh_gloo_solo_io_v1alpha2.NewMulticlusterCertificateRequestClient(b.client).Cluster(cluster)
	if err != nil {
		return err
	}

	if opts.Verifier != nil {
		mgr, err := b.clusters.Cluster(cluster)
		if err != nil {
			return err
		}

		gvk := schema.GroupVersionKind{
			Group:   "certificates.mesh.gloo.solo.io",
			Version: "v1alpha2",
			Kind:    "CertificateRequest",
		}

		if resourceRegistered, err := opts.Verifier.VerifyServerResource(
			cluster,
			mgr.GetConfig(),
			gvk,
		); err != nil {
			return err
		} else if !resourceRegistered {
			return nil
		}
	}

	certificateRequestList, err := certificateRequestClient.ListCertificateRequest(ctx, opts.ListOptions...)
	if err != nil {
		return err
	}

	for _, item := range certificateRequestList.Items {
		item := item               // pike
		item.ClusterName = cluster // set cluster for in-memory processing
		certificateRequests.Insert(&item)
	}

	return nil
}
func (b *multiClusterBuilder) insertPodBounceDirectivesFromCluster(ctx context.Context, cluster string, podBounceDirectives certificates_mesh_gloo_solo_io_v1alpha2_sets.PodBounceDirectiveSet, opts ResourceBuildOptions) error {
	podBounceDirectiveClient, err := certificates_mesh_gloo_solo_io_v1alpha2.NewMulticlusterPodBounceDirectiveClient(b.client).Cluster(cluster)
	if err != nil {
		return err
	}

	if opts.Verifier != nil {
		mgr, err := b.clusters.Cluster(cluster)
		if err != nil {
			return err
		}

		gvk := schema.GroupVersionKind{
			Group:   "certificates.mesh.gloo.solo.io",
			Version: "v1alpha2",
			Kind:    "PodBounceDirective",
		}

		if resourceRegistered, err := opts.Verifier.VerifyServerResource(
			cluster,
			mgr.GetConfig(),
			gvk,
		); err != nil {
			return err
		} else if !resourceRegistered {
			return nil
		}
	}

	podBounceDirectiveList, err := podBounceDirectiveClient.ListPodBounceDirective(ctx, opts.ListOptions...)
	if err != nil {
		return err
	}

	for _, item := range podBounceDirectiveList.Items {
		item := item               // pike
		item.ClusterName = cluster // set cluster for in-memory processing
		podBounceDirectives.Insert(&item)
	}

	return nil
}

func (b *multiClusterBuilder) insertSecretsFromCluster(ctx context.Context, cluster string, secrets v1_sets.SecretSet, opts ResourceBuildOptions) error {
	secretClient, err := v1.NewMulticlusterSecretClient(b.client).Cluster(cluster)
	if err != nil {
		return err
	}

	if opts.Verifier != nil {
		mgr, err := b.clusters.Cluster(cluster)
		if err != nil {
			return err
		}

		gvk := schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Secret",
		}

		if resourceRegistered, err := opts.Verifier.VerifyServerResource(
			cluster,
			mgr.GetConfig(),
			gvk,
		); err != nil {
			return err
		} else if !resourceRegistered {
			return nil
		}
	}

	secretList, err := secretClient.ListSecret(ctx, opts.ListOptions...)
	if err != nil {
		return err
	}

	for _, item := range secretList.Items {
		item := item               // pike
		item.ClusterName = cluster // set cluster for in-memory processing
		secrets.Insert(&item)
	}

	return nil
}
func (b *multiClusterBuilder) insertPodsFromCluster(ctx context.Context, cluster string, pods v1_sets.PodSet, opts ResourceBuildOptions) error {
	podClient, err := v1.NewMulticlusterPodClient(b.client).Cluster(cluster)
	if err != nil {
		return err
	}

	if opts.Verifier != nil {
		mgr, err := b.clusters.Cluster(cluster)
		if err != nil {
			return err
		}

		gvk := schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		}

		if resourceRegistered, err := opts.Verifier.VerifyServerResource(
			cluster,
			mgr.GetConfig(),
			gvk,
		); err != nil {
			return err
		} else if !resourceRegistered {
			return nil
		}
	}

	podList, err := podClient.ListPod(ctx, opts.ListOptions...)
	if err != nil {
		return err
	}

	for _, item := range podList.Items {
		item := item               // pike
		item.ClusterName = cluster // set cluster for in-memory processing
		pods.Insert(&item)
	}

	return nil
}

// build a snapshot from resources in a single cluster
type singleClusterBuilder struct {
	mgr manager.Manager
}

// Produces snapshots of resources across all clusters defined in the ClusterSet
func NewSingleClusterBuilder(
	mgr manager.Manager,
) Builder {
	return &singleClusterBuilder{
		mgr: mgr,
	}
}

func (b *singleClusterBuilder) BuildSnapshot(ctx context.Context, name string, opts BuildOptions) (Snapshot, error) {

	issuedCertificates := certificates_mesh_gloo_solo_io_v1alpha2_sets.NewIssuedCertificateSet()
	certificateRequests := certificates_mesh_gloo_solo_io_v1alpha2_sets.NewCertificateRequestSet()
	podBounceDirectives := certificates_mesh_gloo_solo_io_v1alpha2_sets.NewPodBounceDirectiveSet()

	secrets := v1_sets.NewSecretSet()
	pods := v1_sets.NewPodSet()

	var errs error

	if err := b.insertIssuedCertificates(ctx, issuedCertificates, opts.IssuedCertificates); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := b.insertCertificateRequests(ctx, certificateRequests, opts.CertificateRequests); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := b.insertPodBounceDirectives(ctx, podBounceDirectives, opts.PodBounceDirectives); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := b.insertSecrets(ctx, secrets, opts.Secrets); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := b.insertPods(ctx, pods, opts.Pods); err != nil {
		errs = multierror.Append(errs, err)
	}

	outputSnap := NewSnapshot(
		name,

		issuedCertificates,
		certificateRequests,
		podBounceDirectives,
		secrets,
		pods,
	)

	return outputSnap, errs
}

func (b *singleClusterBuilder) insertIssuedCertificates(ctx context.Context, issuedCertificates certificates_mesh_gloo_solo_io_v1alpha2_sets.IssuedCertificateSet, opts ResourceBuildOptions) error {

	if opts.Verifier != nil {
		gvk := schema.GroupVersionKind{
			Group:   "certificates.mesh.gloo.solo.io",
			Version: "v1alpha2",
			Kind:    "IssuedCertificate",
		}

		if resourceRegistered, err := opts.Verifier.VerifyServerResource(
			"", // verify in the local cluster
			b.mgr.GetConfig(),
			gvk,
		); err != nil {
			return err
		} else if !resourceRegistered {
			return nil
		}
	}

	issuedCertificateList, err := certificates_mesh_gloo_solo_io_v1alpha2.NewIssuedCertificateClient(b.mgr.GetClient()).ListIssuedCertificate(ctx, opts.ListOptions...)
	if err != nil {
		return err
	}

	for _, item := range issuedCertificateList.Items {
		item := item // pike
		issuedCertificates.Insert(&item)
	}

	return nil
}
func (b *singleClusterBuilder) insertCertificateRequests(ctx context.Context, certificateRequests certificates_mesh_gloo_solo_io_v1alpha2_sets.CertificateRequestSet, opts ResourceBuildOptions) error {

	if opts.Verifier != nil {
		gvk := schema.GroupVersionKind{
			Group:   "certificates.mesh.gloo.solo.io",
			Version: "v1alpha2",
			Kind:    "CertificateRequest",
		}

		if resourceRegistered, err := opts.Verifier.VerifyServerResource(
			"", // verify in the local cluster
			b.mgr.GetConfig(),
			gvk,
		); err != nil {
			return err
		} else if !resourceRegistered {
			return nil
		}
	}

	certificateRequestList, err := certificates_mesh_gloo_solo_io_v1alpha2.NewCertificateRequestClient(b.mgr.GetClient()).ListCertificateRequest(ctx, opts.ListOptions...)
	if err != nil {
		return err
	}

	for _, item := range certificateRequestList.Items {
		item := item // pike
		certificateRequests.Insert(&item)
	}

	return nil
}
func (b *singleClusterBuilder) insertPodBounceDirectives(ctx context.Context, podBounceDirectives certificates_mesh_gloo_solo_io_v1alpha2_sets.PodBounceDirectiveSet, opts ResourceBuildOptions) error {

	if opts.Verifier != nil {
		gvk := schema.GroupVersionKind{
			Group:   "certificates.mesh.gloo.solo.io",
			Version: "v1alpha2",
			Kind:    "PodBounceDirective",
		}

		if resourceRegistered, err := opts.Verifier.VerifyServerResource(
			"", // verify in the local cluster
			b.mgr.GetConfig(),
			gvk,
		); err != nil {
			return err
		} else if !resourceRegistered {
			return nil
		}
	}

	podBounceDirectiveList, err := certificates_mesh_gloo_solo_io_v1alpha2.NewPodBounceDirectiveClient(b.mgr.GetClient()).ListPodBounceDirective(ctx, opts.ListOptions...)
	if err != nil {
		return err
	}

	for _, item := range podBounceDirectiveList.Items {
		item := item // pike
		podBounceDirectives.Insert(&item)
	}

	return nil
}

func (b *singleClusterBuilder) insertSecrets(ctx context.Context, secrets v1_sets.SecretSet, opts ResourceBuildOptions) error {

	if opts.Verifier != nil {
		gvk := schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Secret",
		}

		if resourceRegistered, err := opts.Verifier.VerifyServerResource(
			"", // verify in the local cluster
			b.mgr.GetConfig(),
			gvk,
		); err != nil {
			return err
		} else if !resourceRegistered {
			return nil
		}
	}

	secretList, err := v1.NewSecretClient(b.mgr.GetClient()).ListSecret(ctx, opts.ListOptions...)
	if err != nil {
		return err
	}

	for _, item := range secretList.Items {
		item := item // pike
		secrets.Insert(&item)
	}

	return nil
}
func (b *singleClusterBuilder) insertPods(ctx context.Context, pods v1_sets.PodSet, opts ResourceBuildOptions) error {

	if opts.Verifier != nil {
		gvk := schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		}

		if resourceRegistered, err := opts.Verifier.VerifyServerResource(
			"", // verify in the local cluster
			b.mgr.GetConfig(),
			gvk,
		); err != nil {
			return err
		} else if !resourceRegistered {
			return nil
		}
	}

	podList, err := v1.NewPodClient(b.mgr.GetClient()).ListPod(ctx, opts.ListOptions...)
	if err != nil {
		return err
	}

	for _, item := range podList.Items {
		item := item // pike
		pods.Insert(&item)
	}

	return nil
}
