package main

func InitAuth(config *Config) {
	jwtKey = []byte(config.JWTSecret)
}