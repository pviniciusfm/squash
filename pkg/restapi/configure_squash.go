// Code generated by go-swagger; DO NOT EDIT.

package restapi

import (
	"crypto/tls"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	graceful "github.com/tylerb/graceful"

	"github.com/solo-io/squash/pkg/platforms"
	"github.com/solo-io/squash/pkg/platforms/debug"
	"github.com/solo-io/squash/pkg/platforms/kubernetes"
	"github.com/solo-io/squash/pkg/restapi/operations"
	"github.com/solo-io/squash/pkg/restapi/operations/debugattachment"
	"github.com/solo-io/squash/pkg/restapi/operations/debugrequest"
	"github.com/solo-io/squash/pkg/server"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target ../pkg --name Squash --spec ../api.yaml --exclude-main

var (
	options = &struct {
		Cluster    string `long:"cluster" short:"c" description:"Cluster type. currently only kube"`
		KubeConfig string `long:"kubeconfug" description:"Path to kube config file to use to connect to the cluster. in cluster is used if not provided"`
		KubeUrl    string `long:"kubeurl" description:"Kube url (useful with kubectl proxy)"`
	}{"kube", "", ""}
)

func configureFlags(api *operations.SquashAPI) {

	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			LongDescription:  "cluster related operations",
			ShortDescription: "cluster",
			Options:          options,
		},
	}
}

func configureAPI(api *operations.SquashAPI) http.Handler {
	log.SetLevel(log.DebugLevel)
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	var cl platforms.ContainerLocator = nil

	data := server.NewServerData()
	dataStore := server.NewNOPDataStore()

	switch options.Cluster {
	case "kube":
		var config *rest.Config = nil
		if options.KubeConfig != "" {
			var err error
			config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				&clientcmd.ClientConfigLoadingRules{ExplicitPath: options.KubeConfig},
				&clientcmd.ConfigOverrides{
					ClusterInfo: clientcmdapi.Cluster{
						Server: options.KubeUrl,
					},
				}).ClientConfig()

			if err != nil {
				log.Fatalln(err)
			}
		} else if options.KubeUrl != "" {
			config = &rest.Config{Host: options.KubeUrl}
		}

		kube, err := kubernetes.NewKubeOperations(nil, config)
		if err != nil {
			log.Fatalln(err)
		}
		cl = kube
		dataStore, err = kubernetes.NewCRDDataStore(config, data)
		if err != nil {
			log.Fatalln(err)
		}
	case "debug":
		var d debug.DebugPlatform
		cl = &d
	default:
		panic("Invalid cluster option. perhaps you meant 'debug'?")
	}

	handler := server.NewRestHandler(data, cl, dataStore)

	api.DebugattachmentAddDebugAttachmentHandler = debugattachment.AddDebugAttachmentHandlerFunc(handler.DebugattachmentAddDebugAttachmentHandler)
	api.DebugrequestCreateDebugRequestHandler = debugrequest.CreateDebugRequestHandlerFunc(handler.DebugrequestCreateDebugRequestHandler)
	api.DebugattachmentDeleteDebugAttachmentHandler = debugattachment.DeleteDebugAttachmentHandlerFunc(handler.DebugattachmentDeleteDebugAttachmentHandler)
	api.DebugrequestDeleteDebugRequestHandler = debugrequest.DeleteDebugRequestHandlerFunc(handler.DebugrequestDeleteDebugRequestHandler)
	api.DebugattachmentGetDebugAttachmentHandler = debugattachment.GetDebugAttachmentHandlerFunc(handler.DebugattachmentGetDebugAttachmentHandler)
	api.DebugattachmentGetDebugAttachmentsHandler = debugattachment.GetDebugAttachmentsHandlerFunc(handler.DebugattachmentGetDebugAttachmentsHandler)
	api.DebugrequestGetDebugRequestsHandler = debugrequest.GetDebugRequestsHandlerFunc(handler.DebugrequestGetDebugRequestsHandler)
	api.DebugrequestGetDebugRequestHandler = debugrequest.GetDebugRequestHandlerFunc(handler.DebugrequestGetDebugRequestHandler)
	api.DebugattachmentPatchDebugAttachmentHandler = debugattachment.PatchDebugAttachmentHandlerFunc(handler.DebugattachmentPatchDebugAttachmentHandler)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
