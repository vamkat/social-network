package configs

type Config struct {
	Port        string
	FileService FileService
	Clients     Clients
}

type FileService struct {
	Endpoint  string
	AccessKey string
	Secret    string
	Buckets   Buckets
}

// !!! Only use string types here !!!
type Buckets struct {
	Originals string
	Variants  string
}

type Clients struct {
}
