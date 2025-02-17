package app

import (
	"context"
	"github.com/Gustcat/auth/internal/closer"
	"github.com/Gustcat/auth/internal/config"
	"github.com/Gustcat/auth/internal/interceptor"
	descAccess "github.com/Gustcat/auth/pkg/access_v1"
	descAuth "github.com/Gustcat/auth/pkg/auth_v1"
	descUser "github.com/Gustcat/auth/pkg/user_v1"
	_ "github.com/Gustcat/auth/statik"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		app.initConfig,
		app.initServiceProvider,
		app.initGRPCServer,
		app.initHTTPServer,
		app.initSwaggerServer,
	}
	for _, init := range inits {
		if err := init(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		err := app.runGRPCServer()
		if err != nil {
			log.Fatalf("GRPC server run failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := app.runHTTPServer()
		if err != nil {
			log.Fatalf("HTTP server run failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := app.runSwaggerServer()
		if err != nil {
			log.Fatalf("Swagger server run failed: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (app *App) initConfig(_ context.Context) error {
	err := config.Load("local.env")
	if err != nil {
		return err
	}
	return nil
}

func (app *App) initServiceProvider(_ context.Context) error {
	app.serviceProvider = newServiceProvider()
	return nil
}

func (app *App) initGRPCServer(ctx context.Context) error {
	app.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	reflection.Register(app.grpcServer)

	descUser.RegisterUserV1Server(app.grpcServer, app.serviceProvider.UserImpl(ctx))
	descAuth.RegisterAuthV1Server(app.grpcServer, app.serviceProvider.AuthImpl(ctx))
	descAccess.RegisterAccessV1Server(app.grpcServer, app.serviceProvider.AccessImpl(ctx))

	return nil
}

func (app *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := descUser.RegisterUserV1HandlerFromEndpoint(ctx, mux, app.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Content-length", "Accept"},
		AllowCredentials: true,
	})

	app.httpServer = &http.Server{
		Addr:    app.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (app *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	app.swaggerServer = &http.Server{
		Addr:    app.serviceProvider.SwaggerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (app *App) runGRPCServer() error {
	log.Printf("GRPC server listen:  %s", app.serviceProvider.GRPCConfig().Address())
	lis, err := net.Listen("tcp", app.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = app.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) runHTTPServer() error {
	log.Printf("HTTP server listen:  %s", app.serviceProvider.HTTPConfig().Address())

	err := app.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) runSwaggerServer() error {
	log.Printf("Swagger server listen: %s", app.serviceProvider.SwaggerConfig().Address())

	err := app.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
