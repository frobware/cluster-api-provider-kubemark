/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// NOTE: Boilerplate only.  Ignore this file.

// Package v1beta1 contains API Schema definitions for the kubemarkproviderconfig v1beta1 API group
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=package,register
// +k8s:conversion-gen=github.com/openshift/cluster-api-provider-kubemark/pkg/apis/kubemarkproviderconfig
// +k8s:defaulter-gen=TypeMeta
// +groupName=kubemarkproviderconfig.k8s.io
package v1beta1

import (
	"bytes"
	"fmt"

	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "kubemarkproviderconfig.k8s.io", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

// KubemarkProviderConfigCodec is a runtime codec for the provider configuration
// +k8s:deepcopy-gen=false
type KubemarkProviderConfigCodec struct {
	encoder runtime.Encoder
	decoder runtime.Decoder
}

// NewScheme creates a new Scheme
func NewScheme() (*runtime.Scheme, error) {
	return SchemeBuilder.Build()
}

// NewCodec creates a serializer/deserializer for the provider configuration
func NewCodec() (*KubemarkProviderConfigCodec, error) {
	scheme, err := NewScheme()
	if err != nil {
		return nil, err
	}
	codecFactory := serializer.NewCodecFactory(scheme)
	encoder, err := newEncoder(&codecFactory)
	if err != nil {
		return nil, err
	}
	codec := KubemarkProviderConfigCodec{
		encoder: encoder,
		decoder: codecFactory.UniversalDecoder(SchemeGroupVersion),
	}
	return &codec, nil
}

// DecodeProviderSpec deserialises an object from the provider config
func (codec *KubemarkProviderConfigCodec) DecodeProviderSpec(providerSpec *clusterv1.ProviderSpec, out runtime.Object) error {
	if providerSpec.Value != nil {
		_, _, err := codec.decoder.Decode(providerSpec.Value.Raw, nil, out)
		if err != nil {
			return fmt.Errorf("decoding failure: %v", err)
		}
	}
	return nil
}

// EncodeProviderSpec serialises an object to the provider config
func (codec *KubemarkProviderConfigCodec) EncodeProviderSpec(in runtime.Object) (*clusterv1.ProviderSpec, error) {
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}
	return &clusterv1.ProviderSpec{
		Value: &runtime.RawExtension{Raw: buf.Bytes()},
	}, nil
}

// EncodeProviderStatus serialises the provider status
func (codec *KubemarkProviderConfigCodec) EncodeProviderStatus(in runtime.Object) (*runtime.RawExtension, error) {
	var buf bytes.Buffer
	if err := codec.encoder.Encode(in, &buf); err != nil {
		return nil, fmt.Errorf("encoding failed: %v", err)
	}
	return &runtime.RawExtension{Raw: buf.Bytes()}, nil
}

// DecodeProviderStatus deserialises the provider status
func (codec *KubemarkProviderConfigCodec) DecodeProviderStatus(providerStatus *runtime.RawExtension, out runtime.Object) error {
	if providerStatus != nil {
		_, _, err := codec.decoder.Decode(providerStatus.Raw, nil, out)
		if err != nil {
			return fmt.Errorf("decoding failure: %v", err)
		}
	}
	return nil
}

func newEncoder(codecFactory *serializer.CodecFactory) (runtime.Encoder, error) {
	serializerInfos := codecFactory.SupportedMediaTypes()
	if len(serializerInfos) == 0 {
		return nil, fmt.Errorf("unable to find any serlializers")
	}
	encoder := codecFactory.EncoderForVersion(serializerInfos[0].Serializer, SchemeGroupVersion)
	return encoder, nil
}
