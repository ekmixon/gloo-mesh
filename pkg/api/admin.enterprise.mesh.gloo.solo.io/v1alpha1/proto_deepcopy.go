// Code generated by skv2. DO NOT EDIT.

// This file contains generated Deepcopy methods for proto-based Spec and Status fields

package v1alpha1

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/solo-io/protoc-gen-ext/pkg/clone"
)

// DeepCopyInto for the IstioInstallation.Spec
func (in *IstioInstallationSpec) DeepCopyInto(out *IstioInstallationSpec) {
	var p *IstioInstallationSpec
	if h, ok := interface{}(in).(clone.Cloner); ok {
		p = h.Clone().(*IstioInstallationSpec)
	} else {
		p = proto.Clone(in).(*IstioInstallationSpec)
	}
	*out = *p
}

// DeepCopyInto for the IstioInstallation.Status
func (in *IstioInstallationStatus) DeepCopyInto(out *IstioInstallationStatus) {
	var p *IstioInstallationStatus
	if h, ok := interface{}(in).(clone.Cloner); ok {
		p = h.Clone().(*IstioInstallationStatus)
	} else {
		p = proto.Clone(in).(*IstioInstallationStatus)
	}
	*out = *p
}
