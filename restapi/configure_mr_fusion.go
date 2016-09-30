package restapi

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"golang.org/x/net/context"

	"github.com/influxdata/mrfusion"
	"github.com/influxdata/mrfusion/bolt"
	"github.com/influxdata/mrfusion/dist"
	"github.com/influxdata/mrfusion/handlers"
	"github.com/influxdata/mrfusion/influx"
	"github.com/influxdata/mrfusion/mock"
	"github.com/influxdata/mrfusion/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name  --spec ../swagger.yaml --with-context

var devFlags = struct {
	Develop bool `short:"d" long:"develop" description:"Run server in develop mode."`
}{}

var influxFlags = struct {
	Server string `short:"s" long:"server" description:"Full URL of InfluxDB server (http://localhost:8086)" env:"INFLUX_HOST"`
}{}

var storeFlags = struct {
	BoltPath string `short:"b" long:"bolt-path" description:"Full path to boltDB file (/Users/somebody/mrfusion.db)" env:"BOLT_PATH"`
}{}

func configureFlags(api *operations.MrFusionAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			ShortDescription: "Develop Mode server",
			LongDescription:  "Server will use the ui/build directory directly.",
			Options:          &devFlags,
		},
		swag.CommandLineOptionsGroup{
			ShortDescription: "Default Time Series Backend",
			LongDescription:  "Specify the url of an InfluxDB server",
			Options:          &influxFlags,
		},
		swag.CommandLineOptionsGroup{
			ShortDescription: "Default Store Backend",
			LongDescription:  "Specify the path to a BoltDB file",
			Options:          &storeFlags,
		},
	}
}

func assets() mrfusion.Assets {
	if devFlags.Develop {
		return &dist.DebugAssets{
			Dir:     "ui/build",
			Default: "ui/build/index.html",
		}
	}
	return &dist.BindataAssets{
		Prefix:  "ui/build",
		Default: "index.html",
	}
}

func configureAPI(api *operations.MrFusionAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	mockHandler := mock.NewHandler()

	api.GetHandler = operations.GetHandlerFunc(mockHandler.AllRoutes)

	if len(storeFlags.BoltPath) > 0 {
		c := bolt.NewClient()
		c.Path = storeFlags.BoltPath
		if err := c.Open(); err != nil {
			panic(err)
		}
		h := handlers.ExplorationStore{
			ExplorationStore: c.ExplorationStore,
		}
		api.DeleteSourcesIDUsersUserIDExplorationsExplorationIDHandler = operations.DeleteSourcesIDUsersUserIDExplorationsExplorationIDHandlerFunc(h.DeleteExploration)
		api.GetSourcesIDUsersUserIDExplorationsExplorationIDHandler = operations.GetSourcesIDUsersUserIDExplorationsExplorationIDHandlerFunc(h.Exploration)
		api.GetSourcesIDUsersUserIDExplorationsHandler = operations.GetSourcesIDUsersUserIDExplorationsHandlerFunc(h.Explorations)
		api.PatchSourcesIDUsersUserIDExplorationsExplorationIDHandler = operations.PatchSourcesIDUsersUserIDExplorationsExplorationIDHandlerFunc(h.UpdateExploration)
		api.PostSourcesIDUsersUserIDExplorationsHandler = operations.PostSourcesIDUsersUserIDExplorationsHandlerFunc(h.NewExploration)
	} else {
		api.DeleteSourcesIDUsersUserIDExplorationsExplorationIDHandler = operations.DeleteSourcesIDUsersUserIDExplorationsExplorationIDHandlerFunc(mockHandler.DeleteExploration)
		api.GetSourcesIDUsersUserIDExplorationsExplorationIDHandler = operations.GetSourcesIDUsersUserIDExplorationsExplorationIDHandlerFunc(mockHandler.Exploration)
		api.GetSourcesIDUsersUserIDExplorationsHandler = operations.GetSourcesIDUsersUserIDExplorationsHandlerFunc(mockHandler.Explorations)
		api.PatchSourcesIDUsersUserIDExplorationsExplorationIDHandler = operations.PatchSourcesIDUsersUserIDExplorationsExplorationIDHandlerFunc(mockHandler.UpdateExploration)
		api.PostSourcesIDUsersUserIDExplorationsHandler = operations.PostSourcesIDUsersUserIDExplorationsHandlerFunc(mockHandler.NewExploration)
	}

	api.DeleteSourcesIDHandler = operations.DeleteSourcesIDHandlerFunc(mockHandler.RemoveSource)
	api.PatchSourcesIDHandler = operations.PatchSourcesIDHandlerFunc(mockHandler.UpdateSource)

	api.GetSourcesHandler = operations.GetSourcesHandlerFunc(mockHandler.Sources)
	api.GetSourcesIDHandler = operations.GetSourcesIDHandlerFunc(mockHandler.SourcesID)
	api.PostSourcesHandler = operations.PostSourcesHandlerFunc(mockHandler.NewSource)

	if len(influxFlags.Server) > 0 {
		c, err := influx.NewClient(influxFlags.Server)
		if err != nil {
			panic(err)
		}
		// TODO: Change to bolt when finished
		h := handlers.InfluxProxy{
			Srcs:       mock.DefaultSourcesStore,
			TimeSeries: c,
		}
		api.PostSourcesIDProxyHandler = operations.PostSourcesIDProxyHandlerFunc(h.Proxy)
	} else {
		api.PostSourcesIDProxyHandler = operations.PostSourcesIDProxyHandlerFunc(mockHandler.Proxy)
	}

	api.DeleteSourcesIDRolesRoleIDHandler = operations.DeleteSourcesIDRolesRoleIDHandlerFunc(func(ctx context.Context, params operations.DeleteSourcesIDRolesRoleIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .DeleteSourcesIDRolesRoleID has not yet been implemented")
	})

	api.DeleteSourcesIDUsersUserIDHandler = operations.DeleteSourcesIDUsersUserIDHandlerFunc(func(ctx context.Context, params operations.DeleteSourcesIDUsersUserIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .DeleteSourcesIDUsersUserID has not yet been implemented")
	})

	api.DeleteDashboardsIDHandler = operations.DeleteDashboardsIDHandlerFunc(func(ctx context.Context, params operations.DeleteDashboardsIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .DeleteDashboardsID has not yet been implemented")
	})
	api.GetDashboardsHandler = operations.GetDashboardsHandlerFunc(func(ctx context.Context, params operations.GetDashboardsParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetDashboards has not yet been implemented")
	})
	api.GetDashboardsIDHandler = operations.GetDashboardsIDHandlerFunc(func(ctx context.Context, params operations.GetDashboardsIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetDashboardsID has not yet been implemented")
	})

	api.GetSourcesIDPermissionsHandler = operations.GetSourcesIDPermissionsHandlerFunc(func(ctx context.Context, params operations.GetSourcesIDPermissionsParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetSourcesIDPermissions has not yet been implemented")
	})
	api.GetSourcesIDRolesHandler = operations.GetSourcesIDRolesHandlerFunc(func(ctx context.Context, params operations.GetSourcesIDRolesParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetSourcesIDRoles has not yet been implemented")
	})
	api.GetSourcesIDRolesRoleIDHandler = operations.GetSourcesIDRolesRoleIDHandlerFunc(func(ctx context.Context, params operations.GetSourcesIDRolesRoleIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetSourcesIDRolesRoleID has not yet been implemented")
	})

	api.GetSourcesIDUsersHandler = operations.GetSourcesIDUsersHandlerFunc(func(ctx context.Context, params operations.GetSourcesIDUsersParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetSourcesIDUsers has not yet been implemented")
	})
	api.GetSourcesIDUsersUserIDHandler = operations.GetSourcesIDUsersUserIDHandlerFunc(func(ctx context.Context, params operations.GetSourcesIDUsersUserIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .GetSourcesIDUsersUserID has not yet been implemented")
	})

	api.PatchSourcesIDRolesRoleIDHandler = operations.PatchSourcesIDRolesRoleIDHandlerFunc(func(ctx context.Context, params operations.PatchSourcesIDRolesRoleIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .PatchSourcesIDRolesRoleID has not yet been implemented")
	})

	api.PatchSourcesIDUsersUserIDHandler = operations.PatchSourcesIDUsersUserIDHandlerFunc(func(ctx context.Context, params operations.PatchSourcesIDUsersUserIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .PatchSourcesIDUsersUserID has not yet been implemented")
	})
	api.PostDashboardsHandler = operations.PostDashboardsHandlerFunc(func(ctx context.Context, params operations.PostDashboardsParams) middleware.Responder {
		return middleware.NotImplemented("operation .PostDashboards has not yet been implemented")
	})

	api.PostSourcesIDRolesHandler = operations.PostSourcesIDRolesHandlerFunc(func(ctx context.Context, params operations.PostSourcesIDRolesParams) middleware.Responder {
		return middleware.NotImplemented("operation .PostSourcesIDRoles has not yet been implemented")
	})
	api.PostSourcesIDUsersHandler = operations.PostSourcesIDUsersHandlerFunc(func(ctx context.Context, params operations.PostSourcesIDUsersParams) middleware.Responder {
		return middleware.NotImplemented("operation .PostSourcesIDUsers has not yet been implemented")
	})

	api.PutDashboardsIDHandler = operations.PutDashboardsIDHandlerFunc(func(ctx context.Context, params operations.PutDashboardsIDParams) middleware.Responder {
		return middleware.NotImplemented("operation .PutDashboardsID has not yet been implemented")
	})

	api.GetSourcesIDMonitoredHandler = operations.GetSourcesIDMonitoredHandlerFunc(mockHandler.MonitoredServices)

	api.ServerShutdown = func() {}

	handler := setupGlobalMiddleware(api.Serve(setupMiddlewares))
	return handler
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		if strings.Contains(r.URL.Path, "/chronograf/v1") {
			handler.ServeHTTP(w, r)
			return
		} else if r.URL.Path == "//" {
			http.Redirect(w, r, "/index.html", http.StatusFound)
		} else {
			assets().Handler().ServeHTTP(w, r)
			return
		}
	})
}
