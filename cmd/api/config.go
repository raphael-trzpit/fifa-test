package main

type Config struct {
	HttpPort string `envconfig:"HTTP_PORT"`
	DBDSN    string `envconfig:"DB_DSN"`
}
