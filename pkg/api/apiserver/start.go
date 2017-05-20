package apiserver

import (
	"fmt"
	"io"
	"net"

	//"github.com/spf13/cobra"

	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	authv1alpha1 "kope.io/auth/pkg/apis/auth/v1alpha1"
	componentconfigv1alpha1 "kope.io/auth/pkg/apis/componentconfig/v1alpha1"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"github.com/spf13/pflag"
)

// const defaultEtcdPathPrefix = "/registry/auth.kope.io"
const defaultEtcdPathPrefix = "/"

type AuthServerOptions struct {
	//RecommendedOptions *genericoptions.RecommendedOptions

	Etcd           *genericoptions.EtcdOptions
	SecureServing  *genericoptions.SecureServingOptions
	//Authentication *DelegatingAuthenticationOptions
	//Authorization  *DelegatingAuthorizationOptions
	//Audit          *genericoptions.AuditLogOptions
	//Features       *genericoptions.FeatureOptions

	StdOut io.Writer
	StdErr io.Writer
}

func NewAuthServerOptions(out, errOut io.Writer) *AuthServerOptions {
	prefix := defaultEtcdPathPrefix
	copier := Scheme
	codec := Codecs.LegacyCodec(authv1alpha1.SchemeGroupVersion, componentconfigv1alpha1.SchemeGroupVersion)

	o := &AuthServerOptions{
		//RecommendedOptions: genericoptions.NewRecommendedOptions(prefix, copier, codec),

		StdOut: out,
		StdErr: errOut,
	}

	o.Etcd = genericoptions.NewEtcdOptions(storagebackend.NewDefaultConfig(prefix, copier, codec))
	o.SecureServing = genericoptions.NewSecureServingOptions()
		//Authentication: NewDelegatingAuthenticationOptions(),
		//Authorization:  NewDelegatingAuthorizationOptions(),
		//Audit:          NewAuditLogOptions(),
		//Features:       NewFeatureOptions(),

	return o
}

//// NewCommandStartAuthServer provides a CLI handler for 'start server' command
//func NewCommandStartAuthServer(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
//	o := NewAuthServerOptions(out, errOut)
//
//	cmd := &cobra.Command{
//		Short: "Launch a user API server",
//		Long:  "Launch a user API server",
//		RunE: func(c *cobra.Command, args []string) error {
//			if err := o.Complete(); err != nil {
//				return err
//			}
//			if err := o.Validate(args); err != nil {
//				return err
//			}
//			if err := o.RunAuthServer(stopCh); err != nil {
//				return err
//			}
//			return nil
//		},
//	}
//
//	flags := cmd.Flags()
//	o.RecommendedOptions.AddFlags(flags)
//
//	return cmd
//}


func (o *AuthServerOptions) AddFlags(fs *pflag.FlagSet) {
	//o.RecommendedOptions.AddFlags(fs)
	o.Etcd.AddFlags(fs)
	o.SecureServing.AddFlags(fs)
	//o.Authentication.AddFlags(fs)
	//o.Authorization.AddFlags(fs)
	//o.Audit.AddFlags(fs)
	//o.Features.AddFlags(fs)
}


func (o AuthServerOptions) Validate(args []string) error {
	return nil
}

func (o *AuthServerOptions) Complete() error {
	return nil
}

func (o AuthServerOptions) Config() (*Config, error) {
	// TODO have a "real" external address
	if err := o.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewConfig(Codecs)
	// 1.6: serverConfig := genericapiserver.NewConfig().WithSerializer(Codecs)
	//if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
	//	return nil, err
	//}

	serverConfig.CorsAllowedOriginList = []string{ ".*" }

	if err := o.Etcd.ApplyTo(serverConfig); err != nil {
		return nil, err
	}
	if err := o.SecureServing.ApplyTo(serverConfig); err != nil {
		return nil, err
	}


	config := &Config{
		GenericConfig: serverConfig,
	}
	return config, nil
}

func (o AuthServerOptions) RunAuthServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.Complete().New()
	if err != nil {
		return err
	}
	return server.GenericAPIServer.PrepareRun().Run(stopCh)
}
