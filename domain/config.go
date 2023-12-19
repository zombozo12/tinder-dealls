package domain

type Config struct {
	Server   Server   `json:"server" validate:"required"`
	Key      string   `json:"key"`
	Database Database `json:"database" validate:"required"`
	JWT      JWT      `json:"jwt" validate:"required"`
	Redis    Redis    `json:"redis" validate:"required"`
}

type Server struct {
	Port int `json:"port" validate:"required"`
}

type Database struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required"`
	Username string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type JWT struct {
	Key        string `json:"key" validate:"required"`
	ExpireTime int    `json:"expire_time" validate:"required"` // in hours
}

type Redis struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required"`
	Password string `json:"password"`
}
