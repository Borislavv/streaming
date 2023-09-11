package resource

type config struct {
	// api
	ApiVersionPrefix    string `env:"API_VERSION_PREFIX" envDefault:"/api/v1"`
	RenderVersionPrefix string `env:"RENDER_VERSION_PREFIX" envDefault:""`
	StaticVersionPrefix string `env:"STATIC_VERSION_PREFIX" envDefault:""`
	// server
	Host      string `env:"RESOURCES_SERVER_HOST" envDefault:"0.0.0.0"`
	Port      string `env:"RESOURCES_SERVER_PORT" envDefault:"8000"`
	Transport string `env:"RESOURCES_SERVER_TRANSPORT_PROTOCOL" envDefault:"tcp"`
	// database
	MongoUri string `env:"MONGO_URI" envDefault:"mongodb://mongodb:27017/streaming"`
	MongoDb  string `env:"MONGO_DATABASE" envDefault:"streaming"`
	// application
	ResourceFormFilename      string `env:"RESOURCE_FORM_FILENAME" envDefault:"resource"`
	InMemoryFileSizeThreshold int64  `env:"IN_MEMORY_FILE_SIZE_THRESHOLD" envDefault:"104857600"`
}
