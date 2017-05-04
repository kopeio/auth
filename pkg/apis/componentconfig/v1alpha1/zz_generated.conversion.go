// +build !ignore_autogenerated

/*
Copyright 2017 The Kubernetes Authors.

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

// This file was autogenerated by conversion-gen. Do not edit it manually!

package v1alpha1

import (
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	componentconfig "kope.io/auth/pkg/apis/componentconfig"
)

func init() {
	SchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedConversionFuncs(
		Convert_v1alpha1_AuthConfiguration_To_componentconfig_AuthConfiguration,
		Convert_componentconfig_AuthConfiguration_To_v1alpha1_AuthConfiguration,
		Convert_v1alpha1_AuthConfigurationSpec_To_componentconfig_AuthConfigurationSpec,
		Convert_componentconfig_AuthConfigurationSpec_To_v1alpha1_AuthConfigurationSpec,
		Convert_v1alpha1_AuthProviderSpec_To_componentconfig_AuthProviderSpec,
		Convert_componentconfig_AuthProviderSpec_To_v1alpha1_AuthProviderSpec,
		Convert_v1alpha1_GenerateKubeconfig_To_componentconfig_GenerateKubeconfig,
		Convert_componentconfig_GenerateKubeconfig_To_v1alpha1_GenerateKubeconfig,
		Convert_v1alpha1_OAuthConfig_To_componentconfig_OAuthConfig,
		Convert_componentconfig_OAuthConfig_To_v1alpha1_OAuthConfig,
	)
}

func autoConvert_v1alpha1_AuthConfiguration_To_componentconfig_AuthConfiguration(in *AuthConfiguration, out *componentconfig.AuthConfiguration, s conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_AuthConfigurationSpec_To_componentconfig_AuthConfigurationSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1alpha1_AuthConfiguration_To_componentconfig_AuthConfiguration(in *AuthConfiguration, out *componentconfig.AuthConfiguration, s conversion.Scope) error {
	return autoConvert_v1alpha1_AuthConfiguration_To_componentconfig_AuthConfiguration(in, out, s)
}

func autoConvert_componentconfig_AuthConfiguration_To_v1alpha1_AuthConfiguration(in *componentconfig.AuthConfiguration, out *AuthConfiguration, s conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_componentconfig_AuthConfigurationSpec_To_v1alpha1_AuthConfigurationSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	return nil
}

func Convert_componentconfig_AuthConfiguration_To_v1alpha1_AuthConfiguration(in *componentconfig.AuthConfiguration, out *AuthConfiguration, s conversion.Scope) error {
	return autoConvert_componentconfig_AuthConfiguration_To_v1alpha1_AuthConfiguration(in, out, s)
}

func autoConvert_v1alpha1_AuthConfigurationSpec_To_componentconfig_AuthConfigurationSpec(in *AuthConfigurationSpec, out *componentconfig.AuthConfigurationSpec, s conversion.Scope) error {
	if in.AuthProviders != nil {
		in, out := &in.AuthProviders, &out.AuthProviders
		*out = make([]componentconfig.AuthProviderSpec, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_AuthProviderSpec_To_componentconfig_AuthProviderSpec(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.AuthProviders = nil
	}
	if in.GenerateKubeconfig != nil {
		in, out := &in.GenerateKubeconfig, &out.GenerateKubeconfig
		*out = new(componentconfig.GenerateKubeconfig)
		if err := Convert_v1alpha1_GenerateKubeconfig_To_componentconfig_GenerateKubeconfig(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.GenerateKubeconfig = nil
	}
	return nil
}

func Convert_v1alpha1_AuthConfigurationSpec_To_componentconfig_AuthConfigurationSpec(in *AuthConfigurationSpec, out *componentconfig.AuthConfigurationSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_AuthConfigurationSpec_To_componentconfig_AuthConfigurationSpec(in, out, s)
}

func autoConvert_componentconfig_AuthConfigurationSpec_To_v1alpha1_AuthConfigurationSpec(in *componentconfig.AuthConfigurationSpec, out *AuthConfigurationSpec, s conversion.Scope) error {
	if in.AuthProviders != nil {
		in, out := &in.AuthProviders, &out.AuthProviders
		*out = make([]AuthProviderSpec, len(*in))
		for i := range *in {
			if err := Convert_componentconfig_AuthProviderSpec_To_v1alpha1_AuthProviderSpec(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.AuthProviders = nil
	}
	if in.GenerateKubeconfig != nil {
		in, out := &in.GenerateKubeconfig, &out.GenerateKubeconfig
		*out = new(GenerateKubeconfig)
		if err := Convert_componentconfig_GenerateKubeconfig_To_v1alpha1_GenerateKubeconfig(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.GenerateKubeconfig = nil
	}
	return nil
}

func Convert_componentconfig_AuthConfigurationSpec_To_v1alpha1_AuthConfigurationSpec(in *componentconfig.AuthConfigurationSpec, out *AuthConfigurationSpec, s conversion.Scope) error {
	return autoConvert_componentconfig_AuthConfigurationSpec_To_v1alpha1_AuthConfigurationSpec(in, out, s)
}

func autoConvert_v1alpha1_AuthProviderSpec_To_componentconfig_AuthProviderSpec(in *AuthProviderSpec, out *componentconfig.AuthProviderSpec, s conversion.Scope) error {
	out.ID = in.ID
	out.Name = in.Name
	if in.OAuthConfig != nil {
		in, out := &in.OAuthConfig, &out.OAuthConfig
		*out = new(componentconfig.OAuthConfig)
		if err := Convert_v1alpha1_OAuthConfig_To_componentconfig_OAuthConfig(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.OAuthConfig = nil
	}
	out.PermitEmails = in.PermitEmails
	return nil
}

func Convert_v1alpha1_AuthProviderSpec_To_componentconfig_AuthProviderSpec(in *AuthProviderSpec, out *componentconfig.AuthProviderSpec, s conversion.Scope) error {
	return autoConvert_v1alpha1_AuthProviderSpec_To_componentconfig_AuthProviderSpec(in, out, s)
}

func autoConvert_componentconfig_AuthProviderSpec_To_v1alpha1_AuthProviderSpec(in *componentconfig.AuthProviderSpec, out *AuthProviderSpec, s conversion.Scope) error {
	out.ID = in.ID
	out.Name = in.Name
	if in.OAuthConfig != nil {
		in, out := &in.OAuthConfig, &out.OAuthConfig
		*out = new(OAuthConfig)
		if err := Convert_componentconfig_OAuthConfig_To_v1alpha1_OAuthConfig(*in, *out, s); err != nil {
			return err
		}
	} else {
		out.OAuthConfig = nil
	}
	out.PermitEmails = in.PermitEmails
	return nil
}

func Convert_componentconfig_AuthProviderSpec_To_v1alpha1_AuthProviderSpec(in *componentconfig.AuthProviderSpec, out *AuthProviderSpec, s conversion.Scope) error {
	return autoConvert_componentconfig_AuthProviderSpec_To_v1alpha1_AuthProviderSpec(in, out, s)
}

func autoConvert_v1alpha1_GenerateKubeconfig_To_componentconfig_GenerateKubeconfig(in *GenerateKubeconfig, out *componentconfig.GenerateKubeconfig, s conversion.Scope) error {
	out.Server = in.Server
	out.Name = in.Name
	return nil
}

func Convert_v1alpha1_GenerateKubeconfig_To_componentconfig_GenerateKubeconfig(in *GenerateKubeconfig, out *componentconfig.GenerateKubeconfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_GenerateKubeconfig_To_componentconfig_GenerateKubeconfig(in, out, s)
}

func autoConvert_componentconfig_GenerateKubeconfig_To_v1alpha1_GenerateKubeconfig(in *componentconfig.GenerateKubeconfig, out *GenerateKubeconfig, s conversion.Scope) error {
	out.Server = in.Server
	out.Name = in.Name
	return nil
}

func Convert_componentconfig_GenerateKubeconfig_To_v1alpha1_GenerateKubeconfig(in *componentconfig.GenerateKubeconfig, out *GenerateKubeconfig, s conversion.Scope) error {
	return autoConvert_componentconfig_GenerateKubeconfig_To_v1alpha1_GenerateKubeconfig(in, out, s)
}

func autoConvert_v1alpha1_OAuthConfig_To_componentconfig_OAuthConfig(in *OAuthConfig, out *componentconfig.OAuthConfig, s conversion.Scope) error {
	out.ClientID = in.ClientID
	out.ClientSecret = in.ClientSecret
	return nil
}

func Convert_v1alpha1_OAuthConfig_To_componentconfig_OAuthConfig(in *OAuthConfig, out *componentconfig.OAuthConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_OAuthConfig_To_componentconfig_OAuthConfig(in, out, s)
}

func autoConvert_componentconfig_OAuthConfig_To_v1alpha1_OAuthConfig(in *componentconfig.OAuthConfig, out *OAuthConfig, s conversion.Scope) error {
	out.ClientID = in.ClientID
	out.ClientSecret = in.ClientSecret
	return nil
}

func Convert_componentconfig_OAuthConfig_To_v1alpha1_OAuthConfig(in *componentconfig.OAuthConfig, out *OAuthConfig, s conversion.Scope) error {
	return autoConvert_componentconfig_OAuthConfig_To_v1alpha1_OAuthConfig(in, out, s)
}