// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ELBListener) DeepCopyInto(out *ELBListener) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ELBListener.
func (in *ELBListener) DeepCopy() *ELBListener {
	if in == nil {
		return nil
	}
	out := new(ELBListener)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ELBPodInfo) DeepCopyInto(out *ELBPodInfo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ELBPodInfo.
func (in *ELBPodInfo) DeepCopy() *ELBPodInfo {
	if in == nil {
		return nil
	}
	out := new(ELBPodInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ELBService) DeepCopyInto(out *ELBService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ELBService.
func (in *ELBService) DeepCopy() *ELBService {
	if in == nil {
		return nil
	}
	out := new(ELBService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ELBService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ELBServiceList) DeepCopyInto(out *ELBServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ELBService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ELBServiceList.
func (in *ELBServiceList) DeepCopy() *ELBServiceList {
	if in == nil {
		return nil
	}
	out := new(ELBServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ELBServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ELBServiceSpec) DeepCopyInto(out *ELBServiceSpec) {
	*out = *in
	out.Listener = in.Listener
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ELBServiceSpec.
func (in *ELBServiceSpec) DeepCopy() *ELBServiceSpec {
	if in == nil {
		return nil
	}
	out := new(ELBServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ELBServiceStatus) DeepCopyInto(out *ELBServiceStatus) {
	*out = *in
	if in.PodInfos != nil {
		in, out := &in.PodInfos, &out.PodInfos
		*out = make([]ELBPodInfo, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ELBServiceStatus.
func (in *ELBServiceStatus) DeepCopy() *ELBServiceStatus {
	if in == nil {
		return nil
	}
	out := new(ELBServiceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodLabel) DeepCopyInto(out *PodLabel) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodLabel.
func (in *PodLabel) DeepCopy() *PodLabel {
	if in == nil {
		return nil
	}
	out := new(PodLabel)
	in.DeepCopyInto(out)
	return out
}
