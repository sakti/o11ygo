package config

type Configuration struct {
	Debug           bool
	FEPort          int    `required:"true"`
	AuthsvcPort     int    `required:"true"`
	AddersvcPort    int    `required:"true"`
	MultiplysvcPort int    `required:"true"`
	OTLPEndpoint    string `required:"true"`
	DBPath          string `required:"true"`
}
