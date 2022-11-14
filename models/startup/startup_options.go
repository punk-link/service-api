package startup

type StartupOptions struct {
	EnvironmentName string
	GinMode         string
	JaegerEndpoint  string
	ServiceName     string
}
