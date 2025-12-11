package configs

import "os"

type Configs struct {
	PassSecret    string
	JwtSecret     []byte
	EncrytpionKey string

	DbURL      string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	SslMode    string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Ports
	UsersPort         string
	PostsPort         string
	ChatPort          string
	NotificationsPort string
	MediaPort         string
}

var Cfgs Configs

func init() {
	Cfgs.PassSecret = os.Getenv("PASSWORD_SECRET")
	Cfgs.JwtSecret = []byte(os.Getenv("JWT_KEY"))

	// Primary database URL (usually for production/Kube)
	Cfgs.DbURL = os.Getenv("DATABASE_URL")

	// Explicit DB variables (Docker compose local dev)
	Cfgs.DbHost = os.Getenv("DB_HOST")
	Cfgs.DbPort = os.Getenv("DB_PORT")
	Cfgs.DbUser = os.Getenv("DB_USER")
	Cfgs.DbPassword = os.Getenv("DB_PASSWORD")
	Cfgs.DbName = os.Getenv("DB_NAME")
	Cfgs.SslMode = os.Getenv("SSL_MODE")

	// MinIO (used in media + gateway)
	Cfgs.MinioEndpoint = os.Getenv("MINIO_ENDPOINT")
	Cfgs.MinioAccessKey = os.Getenv("MINIO_ACCESS_KEY")
	Cfgs.MinioSecretKey = os.Getenv("MINIO_SECRET_KEY")

	// Redis
	Cfgs.RedisHost = os.Getenv("REDIS_HOST")
	Cfgs.RedisPort = os.Getenv("REDIS_PORT")
	Cfgs.RedisPassword = os.Getenv("REDIS_PASSWORD")

	// PORTS
	Cfgs.UsersPort = os.Getenv("USERS_PORT")
	Cfgs.PostsPort = os.Getenv("POSTS_PORT")
	Cfgs.ChatPort = os.Getenv("CHAT_PORT")
	Cfgs.NotificationsPort = os.Getenv("NOTIFICATIONS_PORT")
	Cfgs.MediaPort = os.Getenv("MEDIA_PORT")
}
