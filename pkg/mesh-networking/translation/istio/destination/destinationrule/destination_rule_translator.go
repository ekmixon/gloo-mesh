package destinationrule

import (
	"context"
	"reflect"

	v1alpha3sets "github.com/solo-io/external-apis/pkg/api/istio/networking.istio.io/v1alpha3/sets"
	discoveryv1sets "github.com/solo-io/gloo-mesh/pkg/api/discovery.mesh.gloo.solo.io/v1/sets"
	settingsv1 "github.com/solo-io/gloo-mesh/pkg/api/settings.mesh.gloo.solo.io/v1"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/istio/decorators/tls"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/istio/decorators/trafficshift"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/istio/destination/utils"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/utils/selectorutils"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/skv2/contrib/pkg/sets"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/rotisserie/eris"
	discoveryv1 "github.com/solo-io/gloo-mesh/pkg/api/discovery.mesh.gloo.solo.io/v1"
	"github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/input"
	networkingv1 "github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/v1"
	v1 "github.com/solo-io/gloo-mesh/pkg/api/networking.mesh.gloo.solo.io/v1"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/reporting"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/istio/decorators"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/utils/fieldutils"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/utils/hostutils"
	"github.com/solo-io/gloo-mesh/pkg/mesh-networking/translation/utils/metautils"
	"github.com/solo-io/skv2/pkg/equalityutils"
	"github.com/solo-io/skv2/pkg/ezkube"
	networkingv1alpha3spec "istio.io/api/networking/v1alpha3"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

//go:generate mockgen -source ./destination_rule_translator.go -destination mocks/destination_rule_translator.go

// the DestinationRule translator translates a Destination into a DestinationRule.
type Translator interface {
	// Translate translates the appropriate DestinationRule for the given Destination.
	// returns nil if no DestinationRule is required for the Destination (i.e. if no DestinationRule features are required, such as subsets).
	//
	// Errors caused by invalid user config will be reported using the Reporter.
	//
	// Note that the input snapshot DestinationSet contains the given Destination.
	Translate(
		ctx context.Context,
		in input.LocalSnapshot,
		destination *discoveryv1.Destination,
		sourceMeshInstallation *discoveryv1.MeshSpec_MeshInstallation,
		reporter reporting.Reporter,
		destinationRuleTrafficPolicyParents []ezkube.ResourceId,
	) *networkingv1alpha3.DestinationRule

	// Return true if the Destination should be translated given the event objects
	// Also return all parent TrafficPolicies
	ShouldTranslate(
		destination *discoveryv1.Destination,
		eventObjs map[schema.GroupVersionKind][]ezkube.ResourceId,
	) (bool, []ezkube.ResourceId)
}

type translator struct {
	settings             *settingsv1.Settings
	userDestinationRules v1alpha3sets.DestinationRuleSet
	clusterDomains       hostutils.ClusterDomainRegistry
	decoratorFactory     decorators.Factory
	destinations         discoveryv1sets.DestinationSet
}

func NewTranslator(
	settings *settingsv1.Settings,
	userDestinationRules v1alpha3sets.DestinationRuleSet,
	clusterDomains hostutils.ClusterDomainRegistry,
	decoratorFactory decorators.Factory,
	destinations discoveryv1sets.DestinationSet,
) Translator {
	return &translator{
		settings:             settings,
		userDestinationRules: userDestinationRules,
		clusterDomains:       clusterDomains,
		decoratorFactory:     decoratorFactory,
		destinations:         destinations,
	}
}

// Translate the Destination into a VirtualService if any of the following has changed:
//  1. the Destination
//  2. applied TrafficPolicy
// Return boolean indicating whether to translate, and all parent TrafficPolicies if they exist
// Note: this could be further optimized if we looked at whether DestinationRule-related fields within the TrafficPolicy changed
func (t *translator) ShouldTranslate(
	destination *discoveryv1.Destination,
	eventObjs map[schema.GroupVersionKind][]ezkube.ResourceId,
) (bool, []ezkube.ResourceId) {
	shouldTranslate := false
	var trafficPolicyParents []ezkube.ResourceId

	for gvk, objs := range eventObjs {
		for _, obj := range objs {
			switch gvk {
			case discoveryv1.DestinationGVK:
				if ezkube.RefsMatch(obj, destination) {
					shouldTranslate = true
				}
			case networkingv1.TrafficPolicyGVK:
				for _, appliedTrafficPolicy := range destination.Status.GetAppliedTrafficPolicies() {
					if ezkube.RefsMatch(obj, appliedTrafficPolicy.Ref) {
						shouldTranslate = true
					}
				}

				// TODO(harveyxia): refactor
				if _, ok := obj.(*networkingv1.TrafficPolicy); ok {
					if utils.ReferencedByTrafficShiftSubset(destination, obj.(*networkingv1.TrafficPolicy)) {
						shouldTranslate = true
						trafficPolicyParents = append(trafficPolicyParents, obj)
					}
				}
			}
		}
	}
	return shouldTranslate, trafficPolicyParents
}

// translate the appropriate DestinationRule for the given Destination.
// returns nil if no DestinationRule is required for the Destination (i.e. if no DestinationRule features are required, such as subsets).
// The input snapshot DestinationSet contains n the
func (t *translator) Translate(
	ctx context.Context,
	in input.LocalSnapshot,
	destination *discoveryv1.Destination,
	sourceMeshInstallation *discoveryv1.MeshSpec_MeshInstallation,
	reporter reporting.Reporter,
	destinationRuleTrafficPolicyParents []ezkube.ResourceId,
) *networkingv1alpha3.DestinationRule {
	kubeService := destination.Spec.GetKubeService()

	if kubeService == nil {
		// TODO(ilackarms): non kube services currently unsupported
		return nil
	}

	sourceClusterName := kubeService.Ref.ClusterName
	if sourceMeshInstallation != nil {
		sourceClusterName = sourceMeshInstallation.Cluster
	}

	destinationRule, err := t.initializeDestinationRule(destination, t.settings.Spec.Mtls, sourceMeshInstallation)
	if err != nil {
		contextutils.LoggerFrom(ctx).Error(err)
		return nil
	}

	// register the owners of the destinationrule fields
	destinationRuleFields := fieldutils.NewOwnershipRegistry()
	drDecorators := t.decoratorFactory.MakeDecorators(decorators.Parameters{
		ClusterDomains: t.clusterDomains,
		Snapshot:       in,
	})

	// Apply decorators which map a single applicable TrafficPolicy to a field on the DestinationRule.
	for _, policy := range destination.Status.AppliedTrafficPolicies {

		// Don't translate the trafficPolicy if the sourceClusterName is not selected by the SourceSelectors
		if !selectorutils.WorkloadSelectorContainsCluster(policy.Spec.SourceSelector, sourceClusterName) {
			continue
		}

		registerField := registerFieldFunc(destinationRuleFields, destinationRule, policy.Ref)
		for _, decorator := range drDecorators {

			if destinationRuleDecorator, ok := decorator.(decorators.TrafficPolicyDestinationRuleDecorator); ok {
				if err := destinationRuleDecorator.ApplyTrafficPolicyToDestinationRule(
					policy,
					destination,
					&destinationRule.Spec,
					registerField,
				); err != nil {
					reporter.ReportTrafficPolicyToDestination(destination, policy.Ref, eris.Wrapf(err, "%v", decorator.DecoratorName()))
				}
			}
		}
	}

	// Append the Destination and all applied TrafficPolicies as parents
	parents := map[schema.GroupVersionKind][]ezkube.ResourceId{
		discoveryv1.DestinationGVK: {destination},
	}
	for _, appliedTp := range destination.Status.GetAppliedTrafficPolicies() {
		parents[networkingv1.TrafficPolicyGVK] = append(parents[networkingv1.TrafficPolicyGVK], appliedTp.Ref)
	}
	for _, tpParent := range destinationRuleTrafficPolicyParents {
		parents[networkingv1.TrafficPolicyGVK] = append(parents[networkingv1.TrafficPolicyGVK], tpParent)
	}
	metautils.AnnotateParents(ctx, destinationRule, parents)

	if t.userDestinationRules == nil {
		return destinationRule
	}

	// detect and report error on intersecting config if enabled in settings
	if errs := conflictsWithUserDestinationRule(
		t.userDestinationRules,
		destinationRule,
	); len(errs) > 0 {
		for _, err := range errs {
			for _, policy := range destination.Status.AppliedTrafficPolicies {
				reporter.ReportTrafficPolicyToDestination(destination, policy.Ref, err)
			}
		}
		return nil
	}

	return destinationRule
}

// construct the callback for registering fields in the virtual service
func registerFieldFunc(
	destinationRuleFields fieldutils.FieldOwnershipRegistry,
	destinationRule *networkingv1alpha3.DestinationRule,
	policy ezkube.ResourceId,
) decorators.RegisterField {
	return func(fieldPtr, val interface{}) error {
		fieldVal := reflect.ValueOf(fieldPtr).Elem().Interface()

		if equalityutils.DeepEqual(fieldVal, val) {
			return nil
		}
		if err := destinationRuleFields.RegisterFieldOwnership(
			destinationRule,
			fieldPtr,
			[]ezkube.ResourceId{policy},
			&v1.TrafficPolicy{},
			0, //TODO(ilackarms): priority
		); err != nil {
			return err
		}
		return nil
	}
}

func (t *translator) initializeDestinationRule(
	destination *discoveryv1.Destination,
	mtlsDefault *v1.TrafficPolicySpec_Policy_MTLS,
	sourceMeshInstallation *discoveryv1.MeshSpec_MeshInstallation,
) (*networkingv1alpha3.DestinationRule, error) {
	var meta metav1.ObjectMeta
	if sourceMeshInstallation != nil {
		meta = metautils.FederatedObjectMeta(
			destination.Spec.GetKubeService().Ref,
			sourceMeshInstallation,
			nil,
		)
	} else {
		meta = metautils.TranslatedObjectMeta(
			destination.Spec.GetKubeService().Ref,
			nil,
		)
	}
	hostname := t.clusterDomains.GetDestinationFQDN(meta.ClusterName, destination.Spec.GetKubeService().Ref)

	destinationRule := &networkingv1alpha3.DestinationRule{
		ObjectMeta: meta,
		Spec: networkingv1alpha3spec.DestinationRule{
			Host:          hostname,
			TrafficPolicy: &networkingv1alpha3spec.TrafficPolicy{},
			Subsets: trafficshift.MakeDestinationRuleSubsetsForDestination(
				destination,
				t.destinations,
				sourceMeshInstallation.GetCluster(),
			),
		},
	}

	// Initialize Istio TLS mode with default declared in Settings
	istioTlsMode, err := tls.MapIstioTlsMode(mtlsDefault.GetIstio().GetTlsMode())
	if err != nil {
		return nil, err
	}
	destinationRule.Spec.TrafficPolicy.Tls = &networkingv1alpha3spec.ClientTLSSettings{
		Mode: istioTlsMode,
	}

	return destinationRule, nil
}

// Return errors for each user-supplied VirtualService that applies to the same hostname as the translated VirtualService
func conflictsWithUserDestinationRule(
	userDestinationRules v1alpha3sets.DestinationRuleSet,
	translatedDestinationRule *networkingv1alpha3.DestinationRule,
) []error {
	// For each user DR, check whether any hosts match any hosts from translated DR
	var errs []error

	// destination rules from RemoteSnapshot only contain non-translated objects
	userDestinationRules.List(func(dr *networkingv1alpha3.DestinationRule) bool {
		// check if common hostnames exist
		commonHostname := utils.CommonHostnames([]string{dr.Spec.Host}, []string{translatedDestinationRule.Spec.Host})
		if len(commonHostname) > 0 {
			errs = append(
				errs,
				eris.Errorf("Unable to translate AppliedTrafficPolicies to DestinationRule, applies to host %s that is already configured by the existing DestinationRule %s", commonHostname[0], sets.Key(dr)),
			)
		}
		return false
	})

	return errs
}
