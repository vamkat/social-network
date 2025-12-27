package entry

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"social-network/services/media/internal/application"
	"social-network/services/media/internal/client"
	"social-network/services/media/internal/configs"
	"social-network/services/media/internal/convertor"
	"social-network/services/media/internal/db/dbservice"
	"social-network/services/media/internal/handler"
	"social-network/services/media/internal/validator"
	"social-network/shared/gen-go/media"
	"social-network/shared/go/ct"
	"social-network/shared/go/gorpc"
	postgresql "social-network/shared/go/postgre"

	"syscall"

	"github.com/minio/minio-go/v7"
)

func Run() error {
	cfgs := getConfigs()

	ctx, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopSignal()

	// Todo: Check ctx functionality on shut down
	// close := tele.InitTelemetry(ctx, "media", ct.CommonKeys(), cfgs.EnableDebugLogs, cfgs.SimplePrint)
	// defer close()

	// tele.Info(ctx, "initialized telemetry")

	pool, err := postgresql.NewPool(ctx, cfgs.DB.URL)
	if err != nil {
		return fmt.Errorf("failed to connect db: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to media database")

	// Internal client for backend operations
	fileServiceClient, err := NewMinIOConn(cfgs.FileService, cfgs.FileService.Endpoint, false)
	if err != nil {
		return err
	}

	// Optional public client for URL generation (e.g. localhost in dev)
	var publicFileServiceClient *minio.Client
	if cfgs.FileService.PublicEndpoint != "" {
		publicFileServiceClient, err = NewMinIOConn(cfgs.FileService, cfgs.FileService.PublicEndpoint, true)
		if err != nil {
			log.Printf("Warning: failed to initialize public MinIO client: %v", err)
		} else {
			log.Println("Initialized public MinIO client for URL generation")
		}
	}

	querier := dbservice.NewQuerier(pool)
	app, err := application.NewMediaService(
		pool,
		&client.Clients{
			Configs:           cfgs.FileService,
			MinIOClient:       fileServiceClient,
			PublicMinIOClient: publicFileServiceClient, //TODO look into eliminating this one, and just using the normal minio client by changing the args based on dev/prod mode
			Validator: &validator.ImageValidator{
				Config: cfgs.FileService.FileConstraints,
			},
			ImageConvertor: convertor.NewImageconvertor(
				cfgs.FileService.FileConstraints),
		},
		querier,
		cfgs,
	)
	if err != nil {
		log.Fatalf("failed to initialize media application: %v", err)
	}

	w := dbservice.NewWorker(querier)

	app.StartVariantWorker(ctx, cfgs.FileService.VariantWorkerInterval)
	w.StartStaleFilesWorker(ctx, cfgs.DB.StaleFilesWorkerInterval)

	service := &handler.MediaHandler{
		Application: app,
		// Configs:     cfgs.Server,
	}

	log.Println("Running gRpc service...")
	startServerFunc, endServerFunc, err := gorpc.CreateGRpcServer[media.MediaServiceServer](
		media.RegisterMediaServiceServer,
		service,
		cfgs.Server.GrpcServerPort,
		ct.CommonKeys(),
	)
	if err != nil {
		return err
	}
	defer endServerFunc()

	go func() {
		err := startServerFunc()
		if err != nil {
			log.Fatal("server failed to start")
		}
		fmt.Println("server finished")
	}()

	// wait here for process termination signal to initiate graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")
	endServerFunc()
	log.Println("Server stopped")
	return nil
}

// type configurations struct {
// 	// Server
// 	ServicePort string `env:"SERVICE_PORT"`

// 	// Database
// 	DatabaseURL              string `env:"DATABASE_URL"`
// 	StaleFilesWorkerInterval int    `env:"STALE_FILES_WORKER_INTERVAL_SECONDS"`

// 	// File service – buckets
// 	BucketOriginals string `env:"BUCKET_ORIGINALS"`
// 	BucketVariants  string `env:"BUCKET_VARIANTS"`

// 	// File service – workers
// 	VariantWorkerInterval int `env:"VARIANT_WORKER_INTERVAL_SECONDS"`

// 	// File constraints
// 	MaxImageUpload int64  `env:"MAX_IMAGE_UPLOAD_BYTES"`
// 	MaxWidth       int    `env:"MAX_IMAGE_WIDTH"`
// 	MaxHeight      int    `env:"MAX_IMAGE_HEIGHT"`
// 	AllowedMIMEs   string `env:"ALLOWED_MIMES"` // comma-separated
// 	AllowedExt     string `env:"ALLOWED_EXT"`   // comma-separated

// 	// MinIO
// 	MinioEndpoint       string `env:"MINIO_ENDPOINT"`
// 	MinioPublicEndpoint string `env:"MINIO_PUBLIC_ENDPOINT"`
// 	MinioAccessKey      string `env:"MINIO_ACCESS_KEY"`
// 	MinioSecretKey      string `env:"MINIO_SECRET_KEY"`
// }

// func getConfigs() configurations {
// 	//vaggelis TODO: add sensible
// 	cfg := configurations{
// 		ServicePort:              "8080",
// 		DatabaseURL:              "postgres://media_user:media_password@localhost:5432/media_db?sslmode=disable",
// 		StaleFilesWorkerInterval: 3600, // 1 hour
// 		BucketOriginals:          "uploads-originals",
// 		BucketVariants:           "uploads-variants",
// 		VariantWorkerInterval:    30, // 30 seconds
// 		MaxImageUpload:           5 << 20,
// 		MaxWidth:                 4096,
// 		MaxHeight:                4096,
// 		AllowedMIMEs:             "image/jpeg,image/jpg,image/png,image/gif,image/webp",
// 		AllowedExt:               ".jpg,.jpeg,.png,.gif,.webp",
// 		MinioEndpoint:            "localhost:9000",
// 		MinioPublicEndpoint:      "",
// 		MinioAccessKey:           "",
// 		MinioSecretKey:           "",
// 	}

// 	if err := configutil.LoadConfigs(&cfg); err != nil {
// 		log.Fatalf("failed to load env variables into config struct: %v", err)
// 	}
// 	return cfg
// }

func getConfigs() configs.Config {
	return configs.Config{
		Server: configs.Server{
			GrpcServerPort: os.Getenv("GRPC_SERVER_PORT"),
		},
		DB: configs.Db{
			URL:                      os.Getenv("DATABASE_URL"),
			StaleFilesWorkerInterval: 1 * time.Hour,
		},
		FileService: configs.FileService{
			Buckets: configs.Buckets{
				Originals: "uploads-originals",
				Variants:  "uploads-variants",
			},
			VariantWorkerInterval: 30 * time.Second,
			FileConstraints: configs.FileConstraints{
				MaxImageUpload: 5 << 20, // 5MB
				MaxWidth:       4096,
				MaxHeight:      4096,
				AllowedMIMEs: map[string]bool{
					"image/jpeg": true,
					"image/jpg":  true,
					"image/png":  true,
					"image/gif":  true,
					"image/webp": true,
				},
				AllowedExt: map[string]bool{
					".jpg":  true,
					".jpeg": true,
					".png":  true,
					".gif":  true,
					".webp": true,
				},
			},
			Endpoint:       os.Getenv("MINIO_ENDPOINT"),
			PublicEndpoint: os.Getenv("MINIO_PUBLIC_ENDPOINT"),
			AccessKey:      os.Getenv("MINIO_ACCESS_KEY"),
			Secret:         os.Getenv("MINIO_SECRET_KEY"),
		},
		EnableDebugLogs: true,
		SimplePrint:     true,
	}
}
