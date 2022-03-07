package structs

import "os"

type NatsConfig struct {
	Host string
	User string
	Pass string
}

func (c *NatsConfig) Defaults() *NatsConfig {

	*c = NatsConfig{
		User: "ruser",
		Pass: "T0pS3cr3t",
	}

	if os.Getenv("NATS_SERVER") != "" {
		c.Host = os.Getenv("NATS_SERVER")
	} else {
		c.Host = "localhost:4222"
	}

	return c
}
