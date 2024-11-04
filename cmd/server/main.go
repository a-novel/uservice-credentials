package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"buf.build/gen/go/a-novel/proto/grpc/go/credentials/v1/credentialsv1grpc"

	"github.com/a-novel/golib/database"
	anovelgrpc "github.com/a-novel/golib/grpc"
	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/adapters"
	"github.com/a-novel/golib/loggers/formatters"

	"github.com/a-novel/uservice-credentials/config"
	"github.com/a-novel/uservice-credentials/migrations"
	"github.com/a-novel/uservice-credentials/pkg/dao"
	"github.com/a-novel/uservice-credentials/pkg/handlers"
	"github.com/a-novel/uservice-credentials/pkg/services"
)

var rpcServices = []grpc.ServiceDesc{
	healthpb.Health_ServiceDesc,
	credentialsv1grpc.CreateService_ServiceDesc,
	credentialsv1grpc.GetService_ServiceDesc,
	credentialsv1grpc.ListService_ServiceDesc,
	credentialsv1grpc.SearchService_ServiceDesc,
	credentialsv1grpc.UpdateService_ServiceDesc,
}

func getDepsCheck(database *bun.DB) *anovelgrpc.DepsCheck {
	return &anovelgrpc.DepsCheck{
		Dependencies: anovelgrpc.DepCheckCallbacks{
			"postgres": database.Ping,
		},
		Services: anovelgrpc.DepCheckServices{
			"create": {"postgres"},
			"get":    {"postgres"},
			"list":   {"postgres"},
			"search": {"postgres"},
			"update": {"postgres"},
		},
	}
}

func main() {
	logger := config.Logger.Formatter

	loader := formatters.NewLoader(
		fmt.Sprintf("Acquiring database connection at %s...", config.App.Postgres.DSN),
		spinner.Meter,
	)
	logger.Log(loader, loggers.LogLevelInfo)

	postgresDB, closePostgresDB, err := database.OpenDB(config.App.Postgres.DSN)
	if err != nil {
		logger.Log(formatters.NewError(err, "open database conn"), loggers.LogLevelFatal)
	}
	defer closePostgresDB()

	logger.Log(
		loader.SetDescription("Database connection successfully acquired.").SetCompleted(),
		loggers.LogLevelInfo,
	)

	if err := database.Migrate(postgresDB, migrations.SQLMigrations, logger); err != nil {
		logger.Log(formatters.NewError(err, "migrate database"), loggers.LogLevelFatal)
	}

	loader = formatters.NewLoader("Setup services...", spinner.Meter)
	logger.Log(loader, loggers.LogLevelInfo)

	grpcReporter := adapters.NewGRPC(logger)

	createCredentialsDAO := dao.NewCreateCredentials(postgresDB)
	getCredentialsDAO := dao.NewGetCredentials(postgresDB)
	listCredentialsDAO := dao.NewListCredentials(postgresDB)
	searchCredentialsDAO := dao.NewSearchCredentials(postgresDB)
	updateCredentialsDAO := dao.NewUpdateCredentials(postgresDB)

	createCredentialsService := services.NewCreateCredentials(createCredentialsDAO)
	getCredentialsService := services.NewGetCredentials(getCredentialsDAO)
	listCredentialsService := services.NewListCredentials(listCredentialsDAO)
	searchCredentialsService := services.NewSearchCredentials(searchCredentialsDAO)
	updateCredentialsService := services.NewUpdateCredentials(updateCredentialsDAO)

	createCredentialsHandler := handlers.NewCreateCredentials(createCredentialsService, grpcReporter)
	getCredentialsHandler := handlers.NewGetCredentials(getCredentialsService, grpcReporter)
	listCredentialsHandler := handlers.NewListCredentials(listCredentialsService, grpcReporter)
	searchCredentialsHandler := handlers.NewSearchCredentials(searchCredentialsService, grpcReporter)
	updateCredentialsHandler := handlers.NewUpdateCredentials(updateCredentialsService, grpcReporter)

	logger.Log(loader.SetDescription("Services successfully setup.").SetCompleted(), loggers.LogLevelInfo)

	listener, server, err := anovelgrpc.StartServer(config.App.Server.Port)
	if err != nil {
		logger.Log(formatters.NewError(err, "start server"), loggers.LogLevelFatal)
	}
	defer anovelgrpc.CloseServer(listener, server)

	reflection.Register(server)
	healthpb.RegisterHealthServer(server, anovelgrpc.NewHealthServer(getDepsCheck(postgresDB), time.Minute))
	credentialsv1grpc.RegisterCreateServiceServer(server, createCredentialsHandler)
	credentialsv1grpc.RegisterGetServiceServer(server, getCredentialsHandler)
	credentialsv1grpc.RegisterListServiceServer(server, listCredentialsHandler)
	credentialsv1grpc.RegisterSearchServiceServer(server, searchCredentialsHandler)
	credentialsv1grpc.RegisterUpdateServiceServer(server, updateCredentialsHandler)

	report := formatters.NewDiscoverGRPC(rpcServices, config.App.Server.Port)
	logger.Log(report, loggers.LogLevelInfo)

	if err := server.Serve(listener); err != nil {
		logger.Log(formatters.NewError(err, "serve"), loggers.LogLevelFatal)
	}
}
