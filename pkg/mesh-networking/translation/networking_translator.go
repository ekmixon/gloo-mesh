package translation

import (
	"bytes"
	"context"
	"fmt"

	"github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/input"
	appmeshoutput "github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/output/appmesh"
	istiooutput "github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/output/istio"
	localoutput "github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/output/local"
	smioutput "github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/output/smi"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/reporting"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/appmesh"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/istio"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/osm"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/utils/metautils"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/skv2/contrib/pkg/output"
	"github.com/solo-io/skv2/pkg/multicluster"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate mockgen -source ./networking_translator.go -destination mocks/networking_translator.go

// the networking translator translates an istio input networking snapshot to an istiooutput snapshot of mesh config resources
type Translator interface {
	// errors reflect an internal translation error and should never happen
	Translate(
		ctx context.Context,
		in input.LocalSnapshot,
		userSupplied input.RemoteSnapshot,
		reporter reporting.Reporter,
	) (*Outputs, error)
}

type translator struct {
	totalTranslates   int // TODO(ilackarms): metric
	istioTranslator   istio.Translator
	appmeshTranslator appmesh.Translator
	osmTranslator     osm.Translator
}

func NewTranslator(
	istioTranslator istio.Translator,
	appmeshTranslator appmesh.Translator,
	osmTranslator osm.Translator,
) Translator {
	return &translator{
		istioTranslator:   istioTranslator,
		appmeshTranslator: appmeshTranslator,
		osmTranslator:     osmTranslator,
	}
}

func (t *translator) Translate(
	ctx context.Context,
	in input.LocalSnapshot,
	userSupplied input.RemoteSnapshot,
	reporter reporting.Reporter,
) (*Outputs, error) {
	t.totalTranslates++

	// TODO(ilackarms): remove this when we remove the applier
	// this is done to distinguish the number of actual translations from the applier-invoked ones
	currentTranslation := t.totalTranslates/2 + 1

	ctx = contextutils.WithLogger(ctx, fmt.Sprintf("translation-%v", currentTranslation))

	istioOutputs := istiooutput.NewBuilder(ctx, fmt.Sprintf("networking-istio-%v", currentTranslation))
	appmeshOutputs := appmeshoutput.NewBuilder(ctx, fmt.Sprintf("networking-appmesh-%v", currentTranslation))
	smiOutputs := smioutput.NewBuilder(ctx, fmt.Sprintf("networking-smi-%v", currentTranslation))
	localOutputs := localoutput.NewBuilder(ctx, fmt.Sprintf("networking-local-%v", currentTranslation))

	t.istioTranslator.Translate(ctx, in, userSupplied, istioOutputs, localOutputs, reporter)

	t.appmeshTranslator.Translate(ctx, in, appmeshOutputs, reporter)

	t.osmTranslator.Translate(ctx, in, smiOutputs, reporter)

	return &Outputs{
		Istio:   istioOutputs,
		Appmesh: appmeshOutputs,
		Smi:     smiOutputs,
		Local:   localOutputs,
	}, nil
}

type Outputs struct {
	Istio   istiooutput.Builder
	Appmesh appmeshoutput.Builder
	Smi     smioutput.Builder
	Local   localoutput.Builder
}

func (t *Outputs) snapshots() (outputSnapshots, error) {
	istioSnapshot, err := t.Istio.BuildSinglePartitionedSnapshot(metautils.TranslatedObjectLabels())
	if err != nil {
		return outputSnapshots{}, err
	}

	appmeshSnapshot, err := t.Appmesh.BuildSinglePartitionedSnapshot(metautils.TranslatedObjectLabels())
	if err != nil {
		return outputSnapshots{}, err
	}

	smiSnapshot, err := t.Smi.BuildSinglePartitionedSnapshot(metautils.TranslatedObjectLabels())
	if err != nil {
		return outputSnapshots{}, err
	}

	localSnapshot, err := t.Local.BuildSinglePartitionedSnapshot(metautils.TranslatedObjectLabels())
	if err != nil {
		return outputSnapshots{}, err
	}

	return outputSnapshots{
		istio:   istioSnapshot,
		appmesh: appmeshSnapshot,
		smi:     smiSnapshot,
		local:   localSnapshot,
	}, nil
}

func (t *Outputs) MarshalJSON() ([]byte, error) {
	snaps, err := t.snapshots()
	if err != nil {
		return nil, err
	}
	return snaps.MarshalJSON()
}

func (t *Outputs) ApplyMultiCluster(
	ctx context.Context,
	clusterClient client.Client,
	multiClusterClient multicluster.Client,
	opts output.OutputOpts,
) error {
	snaps, err := t.snapshots()
	if err != nil {
		return err
	}
	// Apply mesh resources to registered clusters
	snaps.istio.ApplyMultiCluster(ctx, multiClusterClient, opts)
	snaps.appmesh.ApplyMultiCluster(ctx, multiClusterClient, opts)
	snaps.smi.ApplyMultiCluster(ctx, multiClusterClient, opts)
	// Apply local resources only to management cluster
	snaps.local.ApplyLocalCluster(ctx, clusterClient, opts)

	return nil
}

type outputSnapshots struct {
	istio   istiooutput.Snapshot
	appmesh appmeshoutput.Snapshot
	smi     smioutput.Snapshot
	local   localoutput.Snapshot
}

func (t outputSnapshots) MarshalJSON() ([]byte, error) {

	istioByt, err := t.istio.MarshalJSON()
	if err != nil {
		return nil, err
	}
	appmeshByt, err := t.appmesh.MarshalJSON()
	if err != nil {
		return nil, err
	}
	smiByt, err := t.smi.MarshalJSON()
	if err != nil {
		return nil, err
	}
	localByt, err := t.local.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return bytes.Join([][]byte{istioByt, appmeshByt, smiByt, localByt}, []byte("\n")), nil
}
