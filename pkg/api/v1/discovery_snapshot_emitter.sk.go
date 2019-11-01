// Code generated by solo-kit. DO NOT EDIT.

package v1

import (
	"sync"
	"time"

	gloo_solo_io "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	skstats "github.com/solo-io/solo-kit/pkg/stats"

	"github.com/solo-io/go-utils/errutils"
)

var (
	// Deprecated. See mDiscoveryResourcesIn
	mDiscoverySnapshotIn = stats.Int64("discovery.zephyr.solo.io/emitter/snap_in", "Deprecated. Use discovery.zephyr.solo.io/emitter/resources_in. The number of snapshots in", "1")

	// metrics for emitter
	mDiscoveryResourcesIn    = stats.Int64("discovery.zephyr.solo.io/emitter/resources_in", "The number of resource lists received on open watch channels", "1")
	mDiscoverySnapshotOut    = stats.Int64("discovery.zephyr.solo.io/emitter/snap_out", "The number of snapshots out", "1")
	mDiscoverySnapshotMissed = stats.Int64("discovery.zephyr.solo.io/emitter/snap_missed", "The number of snapshots missed", "1")

	// views for emitter
	// deprecated: see discoveryResourcesInView
	discoverysnapshotInView = &view.View{
		Name:        "discovery.zephyr.solo.io/emitter/snap_in",
		Measure:     mDiscoverySnapshotIn,
		Description: "Deprecated. Use discovery.zephyr.solo.io/emitter/resources_in. The number of snapshots updates coming in.",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}

	discoveryResourcesInView = &view.View{
		Name:        "discovery.zephyr.solo.io/emitter/resources_in",
		Measure:     mDiscoveryResourcesIn,
		Description: "The number of resource lists received on open watch channels",
		Aggregation: view.Count(),
		TagKeys: []tag.Key{
			skstats.NamespaceKey,
			skstats.ResourceKey,
		},
	}
	discoverysnapshotOutView = &view.View{
		Name:        "discovery.zephyr.solo.io/emitter/snap_out",
		Measure:     mDiscoverySnapshotOut,
		Description: "The number of snapshots updates going out",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
	discoverysnapshotMissedView = &view.View{
		Name:        "discovery.zephyr.solo.io/emitter/snap_missed",
		Measure:     mDiscoverySnapshotMissed,
		Description: "The number of snapshots updates going missed. this can happen in heavy load. missed snapshot will be re-tried after a second.",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	view.Register(
		discoverysnapshotInView,
		discoverysnapshotOutView,
		discoverysnapshotMissedView,
		discoveryResourcesInView,
	)
}

type DiscoverySnapshotEmitter interface {
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *DiscoverySnapshot, <-chan error, error)
}

type DiscoveryEmitter interface {
	DiscoverySnapshotEmitter
	Register() error
	Pod() github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.PodClient
	Upstream() gloo_solo_io.UpstreamClient
	Deployment() github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.DeploymentClient
}

func NewDiscoveryEmitter(podClient github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.PodClient, upstreamClient gloo_solo_io.UpstreamClient, deploymentClient github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.DeploymentClient) DiscoveryEmitter {
	return NewDiscoveryEmitterWithEmit(podClient, upstreamClient, deploymentClient, make(chan struct{}))
}

func NewDiscoveryEmitterWithEmit(podClient github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.PodClient, upstreamClient gloo_solo_io.UpstreamClient, deploymentClient github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.DeploymentClient, emit <-chan struct{}) DiscoveryEmitter {
	return &discoveryEmitter{
		pod:        podClient,
		upstream:   upstreamClient,
		deployment: deploymentClient,
		forceEmit:  emit,
	}
}

type discoveryEmitter struct {
	forceEmit  <-chan struct{}
	pod        github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.PodClient
	upstream   gloo_solo_io.UpstreamClient
	deployment github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.DeploymentClient
}

func (c *discoveryEmitter) Register() error {
	if err := c.pod.Register(); err != nil {
		return err
	}
	if err := c.upstream.Register(); err != nil {
		return err
	}
	if err := c.deployment.Register(); err != nil {
		return err
	}
	return nil
}

func (c *discoveryEmitter) Pod() github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.PodClient {
	return c.pod
}

func (c *discoveryEmitter) Upstream() gloo_solo_io.UpstreamClient {
	return c.upstream
}

func (c *discoveryEmitter) Deployment() github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.DeploymentClient {
	return c.deployment
}

func (c *discoveryEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *DiscoverySnapshot, <-chan error, error) {

	if len(watchNamespaces) == 0 {
		watchNamespaces = []string{""}
	}

	for _, ns := range watchNamespaces {
		if ns == "" && len(watchNamespaces) > 1 {
			return nil, nil, errors.Errorf("the \"\" namespace is used to watch all namespaces. Snapshots can either be tracked for " +
				"specific namespaces or \"\" AllNamespaces, but not both.")
		}
	}

	errs := make(chan error)
	var done sync.WaitGroup
	ctx := opts.Ctx
	/* Create channel for Pod */
	type podListWithNamespace struct {
		list      github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.PodList
		namespace string
	}
	podChan := make(chan podListWithNamespace)

	var initialPodList github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.PodList
	/* Create channel for Upstream */
	type upstreamListWithNamespace struct {
		list      gloo_solo_io.UpstreamList
		namespace string
	}
	upstreamChan := make(chan upstreamListWithNamespace)

	var initialUpstreamList gloo_solo_io.UpstreamList
	/* Create channel for Deployment */
	type deploymentListWithNamespace struct {
		list      github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.DeploymentList
		namespace string
	}
	deploymentChan := make(chan deploymentListWithNamespace)

	var initialDeploymentList github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.DeploymentList

	currentSnapshot := DiscoverySnapshot{}

	for _, namespace := range watchNamespaces {
		/* Setup namespaced watch for Pod */
		{
			pods, err := c.pod.List(namespace, clients.ListOpts{Ctx: opts.Ctx, Selector: opts.Selector})
			if err != nil {
				return nil, nil, errors.Wrapf(err, "initial Pod list")
			}
			initialPodList = append(initialPodList, pods...)
		}
		podNamespacesChan, podErrs, err := c.pod.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Pod watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, podErrs, namespace+"-pods")
		}(namespace)
		/* Setup namespaced watch for Upstream */
		{
			upstreams, err := c.upstream.List(namespace, clients.ListOpts{Ctx: opts.Ctx, Selector: opts.Selector})
			if err != nil {
				return nil, nil, errors.Wrapf(err, "initial Upstream list")
			}
			initialUpstreamList = append(initialUpstreamList, upstreams...)
		}
		upstreamNamespacesChan, upstreamErrs, err := c.upstream.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Upstream watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, upstreamErrs, namespace+"-upstreams")
		}(namespace)
		/* Setup namespaced watch for Deployment */
		{
			deployments, err := c.deployment.List(namespace, clients.ListOpts{Ctx: opts.Ctx, Selector: opts.Selector})
			if err != nil {
				return nil, nil, errors.Wrapf(err, "initial Deployment list")
			}
			initialDeploymentList = append(initialDeploymentList, deployments...)
		}
		deploymentNamespacesChan, deploymentErrs, err := c.deployment.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Deployment watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(ctx, errs, deploymentErrs, namespace+"-deployments")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-ctx.Done():
					return
				case podList := <-podNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case podChan <- podListWithNamespace{list: podList, namespace: namespace}:
					}
				case upstreamList := <-upstreamNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case upstreamChan <- upstreamListWithNamespace{list: upstreamList, namespace: namespace}:
					}
				case deploymentList := <-deploymentNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case deploymentChan <- deploymentListWithNamespace{list: deploymentList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}
	/* Initialize snapshot for Pods */
	currentSnapshot.Pods = initialPodList.Sort()
	/* Initialize snapshot for Upstreams */
	currentSnapshot.Upstreams = initialUpstreamList.Sort()
	/* Initialize snapshot for Deployments */
	currentSnapshot.Deployments = initialDeploymentList.Sort()

	snapshots := make(chan *DiscoverySnapshot)
	go func() {
		// sent initial snapshot to kick off the watch
		initialSnapshot := currentSnapshot.Clone()
		snapshots <- &initialSnapshot

		timer := time.NewTicker(time.Second * 1)
		previousHash := currentSnapshot.Hash()
		sync := func() {
			currentHash := currentSnapshot.Hash()
			if previousHash == currentHash {
				return
			}

			sentSnapshot := currentSnapshot.Clone()
			select {
			case snapshots <- &sentSnapshot:
				stats.Record(ctx, mDiscoverySnapshotOut.M(1))
				previousHash = currentHash
			default:
				stats.Record(ctx, mDiscoverySnapshotMissed.M(1))
			}
		}
		podsByNamespace := make(map[string]github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.PodList)
		upstreamsByNamespace := make(map[string]gloo_solo_io.UpstreamList)
		deploymentsByNamespace := make(map[string]github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.DeploymentList)

		for {
			record := func() { stats.Record(ctx, mDiscoverySnapshotIn.M(1)) }

			select {
			case <-timer.C:
				sync()
			case <-ctx.Done():
				close(snapshots)
				done.Wait()
				close(errs)
				return
			case <-c.forceEmit:
				sentSnapshot := currentSnapshot.Clone()
				snapshots <- &sentSnapshot
			case podNamespacedList := <-podChan:
				record()

				namespace := podNamespacedList.namespace

				skstats.IncrementResourceCount(
					ctx,
					namespace,
					"pod",
					mDiscoveryResourcesIn,
				)

				// merge lists by namespace
				podsByNamespace[namespace] = podNamespacedList.list
				var podList github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.PodList
				for _, pods := range podsByNamespace {
					podList = append(podList, pods...)
				}
				currentSnapshot.Pods = podList.Sort()
			case upstreamNamespacedList := <-upstreamChan:
				record()

				namespace := upstreamNamespacedList.namespace

				skstats.IncrementResourceCount(
					ctx,
					namespace,
					"upstream",
					mDiscoveryResourcesIn,
				)

				// merge lists by namespace
				upstreamsByNamespace[namespace] = upstreamNamespacedList.list
				var upstreamList gloo_solo_io.UpstreamList
				for _, upstreams := range upstreamsByNamespace {
					upstreamList = append(upstreamList, upstreams...)
				}
				currentSnapshot.Upstreams = upstreamList.Sort()
			case deploymentNamespacedList := <-deploymentChan:
				record()

				namespace := deploymentNamespacedList.namespace

				skstats.IncrementResourceCount(
					ctx,
					namespace,
					"deployment",
					mDiscoveryResourcesIn,
				)

				// merge lists by namespace
				deploymentsByNamespace[namespace] = deploymentNamespacedList.list
				var deploymentList github_com_solo_io_solo_kit_pkg_api_v1_resources_common_kubernetes.DeploymentList
				for _, deployments := range deploymentsByNamespace {
					deploymentList = append(deploymentList, deployments...)
				}
				currentSnapshot.Deployments = deploymentList.Sort()
			}
		}
	}()
	return snapshots, errs, nil
}
