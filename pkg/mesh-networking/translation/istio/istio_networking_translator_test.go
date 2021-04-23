package istio

import (
	"context"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-multierror"
	. "github.com/onsi/ginkgo"
	"github.com/rotisserie/eris"
	discoveryv1 "github.com/solo-io/gloo-mesh/pkg/api/discovery.mesh.gloo.solo.io/v1"
	"github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/input"
	mock_istio_output "github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/output/istio/mocks"
	mock_local_output "github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/output/local/mocks"
	mock_reporting "github.com/solo-io/gloo-mesh/pkg/mesh-networking/reporting/mocks"
	mock_destination "github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/istio/destination/mocks"
	mock_extensions "github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/istio/extensions/mocks"
	mock_istio "github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/istio/internal/mocks"
	mock_mesh "github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/istio/mesh/mocks"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/utils/metautils"
	"github.com/solo-io/go-utils/contextutils"
	multiclusterv1alpha1 "github.com/solo-io/skv2/pkg/api/multicluster.solo.io/v1alpha1"
	"github.com/solo-io/skv2/pkg/ezkube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("IstioNetworkingTranslator", func() {
	var (
		ctrl                      *gomock.Controller
		ctx                       context.Context
		ctxWithValue              context.Context
		mockIstioExtender         *mock_extensions.MockIstioExtender
		mockReporter              *mock_reporting.MockReporter
		mockIstioOutputs          *mock_istio_output.MockBuilder
		mockLocalOutputs          *mock_local_output.MockBuilder
		mockDestinationTranslator *mock_destination.MockTranslator
		mockMeshTranslator        *mock_mesh.MockTranslator
		mockDependencyFactory     *mock_istio.MockDependencyFactory
		translator                Translator
		localEventObjs            = map[schema.GroupVersionKind][]ezkube.ResourceId{}
		remoteEventObjs           = map[schema.GroupVersionKind][]ezkube.ClusterResourceId{}
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		ctx = context.TODO()
		ctxWithValue = contextutils.WithLogger(context.TODO(), "istio-translator-0")
		mockIstioExtender = mock_extensions.NewMockIstioExtender(ctrl)
		mockReporter = mock_reporting.NewMockReporter(ctrl)
		mockDestinationTranslator = mock_destination.NewMockTranslator(ctrl)
		mockMeshTranslator = mock_mesh.NewMockTranslator(ctrl)
		mockDependencyFactory = mock_istio.NewMockDependencyFactory(ctrl)
		mockIstioOutputs = mock_istio_output.NewMockBuilder(ctrl)
		mockLocalOutputs = mock_local_output.NewMockBuilder(ctrl)
		translator = &istioTranslator{dependencies: mockDependencyFactory, extender: mockIstioExtender}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should translate", func() {
		var err error
		err = multierror.Append(err, eris.New(""))
		in := input.NewInputLocalSnapshotManualBuilder("").
			AddKubernetesClusters([]*multiclusterv1alpha1.KubernetesCluster{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "kube-cluster",
						Namespace: "namespace",
					},
				},
			}).
			AddMeshes([]*discoveryv1.Mesh{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mesh-1",
					},
					Spec: discoveryv1.MeshSpec{
						Type: &discoveryv1.MeshSpec_Istio_{
							Istio: &discoveryv1.MeshSpec_Istio{
								Installation: &discoveryv1.MeshSpec_MeshInstallation{
									Cluster: "cluster1",
								},
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mesh-2",
					},
					Spec: discoveryv1.MeshSpec{
						Type: &discoveryv1.MeshSpec_Istio_{
							Istio: &discoveryv1.MeshSpec_Istio{
								Installation: &discoveryv1.MeshSpec_MeshInstallation{
									Cluster: "cluster2",
								},
							},
						},
					},
				},
			}).
			AddWorkloads([]*discoveryv1.Workload{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "mesh-workload-1",
						Labels: metautils.TranslatedObjectLabels(),
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "mesh-workload-2",
						Labels: metautils.TranslatedObjectLabels(),
					},
				},
			}).
			AddDestinations([]*discoveryv1.Destination{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "traffic-target-1",
						Labels: metautils.TranslatedObjectLabels(),
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "traffic-target-2",
						Labels: metautils.TranslatedObjectLabels(),
					},
				},
			}).Build()

		contextMatcher := gomock.Any()
		mockDependencyFactory.
			EXPECT().
			MakeDestinationTranslator(contextMatcher, nil, in.KubernetesClusters(), in.Destinations()).
			Return(mockDestinationTranslator)

		for _, destination := range in.Destinations().List() {
			mockDestinationTranslator.
				EXPECT().
				Translate(localEventObjs, remoteEventObjs, in, destination, mockIstioOutputs, mockReporter)
		}

		mockDependencyFactory.
			EXPECT().
			MakeMeshTranslator(ctxWithValue, in.Secrets(), in.Workloads()).
			Return(mockMeshTranslator)
		for _, mesh := range in.Meshes().List() {
			mockMeshTranslator.
				EXPECT().
				ShouldTranslate(mesh, localEventObjs).
				Return(true)
			mockMeshTranslator.
				EXPECT().
				Translate(in, mesh, mockIstioOutputs, mockLocalOutputs, mockReporter)
		}

		mockIstioOutputs.
			EXPECT().
			AddCluster("cluster1")
		mockIstioOutputs.
			EXPECT().
			AddCluster("cluster2")

		mockIstioExtender.EXPECT().PatchOutputs(contextMatcher, in, mockIstioOutputs)

		translator.Translate(ctx, localEventObjs, remoteEventObjs, in, nil, mockIstioOutputs, mockLocalOutputs, mockReporter)
	})
})
