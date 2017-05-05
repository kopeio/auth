package apiserver

import (
	"fmt"
	"io"
	"net"

	//"github.com/spf13/cobra"

	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"kope.io/auth/pkg/apis/auth/v1alpha1"
)

const defaultEtcdPathPrefix = "/registry/auth.kope.io"

type AuthServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions

	StdOut io.Writer
	StdErr io.Writer
}

func NewAuthServerOptions(out, errOut io.Writer) *AuthServerOptions {
	o := &AuthServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(defaultEtcdPathPrefix, Scheme, Codecs.LegacyCodec(v1alpha1.SchemeGroupVersion)),

		StdOut: out,
		StdErr: errOut,
	}

	return o
}

//// NewCommandStartAuthServer provides a CLI handler for 'start server' command
//func NewCommandStartAuthServer(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
//	o := NewAuthServerOptions(out, errOut)
//
//	cmd := &cobra.Command{
//		Short: "Launch a auth API server",
//		Long:  "Launch a auth API server",
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

func (o AuthServerOptions) Validate(args []string) error {
	return nil
}

func (o *AuthServerOptions) Complete() error {
	return nil
}

func (o AuthServerOptions) Config() (*Config, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, net.ParseIP("127.0.0.1")); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	// 1.7: serverConfig := genericapiserver.NewConfig(Codecs)
	serverConfig := genericapiserver.NewConfig().WithSerializer(Codecs)
	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
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
