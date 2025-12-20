package configs

import "time"

type Config struct {
	Server      Server
	DB          Db
	FileService FileService
	Clients     Clients
}

type FileService struct {
	Endpoint              string
	PublicEndpoint        string // used only for dev. Should be empty for production
	AccessKey             string
	Secret                string
	Buckets               Buckets
	FileConstraints       FileConstraints
	VariantWorkerInterval time.Duration
}

// !!! Only use string types here !!!
type Buckets struct {
	Originals string
	Variants  string
}

type FileConstraints struct {
	MaxImageUpload int64
	AllowedMIMEs   map[string]bool
	AllowedExt     map[string]bool
	MaxWidth       int
	MaxHeight      int
}

type Server struct {
	Port string
}

type Clients struct {
}

type Db struct {
	URL                      string
	StaleFilesWorkerInterval time.Duration
}
