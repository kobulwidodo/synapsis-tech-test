package main

import (
	"go-clean/src/business/domain"
	"go-clean/src/business/usecase"
	"go-clean/src/handler/rest"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/configreader"
	"go-clean/src/lib/midtrans"
	"go-clean/src/lib/redis"
	"go-clean/src/lib/sql"
	"go-clean/src/utils/config"

	_ "go-clean/docs/swagger"
)

// @contact.name   Rakhmad Giffari Nurfadhilah
// @contact.url    https://fadhilmail.tech/
// @contact.email  rakhmadgiffari14@gmail.com

// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization

const (
	configFile string = "./etc/cfg/config.json"
)

func main() {
	cfg := config.Init()
	configReader := configreader.Init(configreader.Options{
		ConfigFile: configFile,
	})
	configReader.ReadConfig(&cfg)

	auth := auth.Init()

	midtrans := midtrans.Init(cfg.Midtrans)

	db := sql.Init(cfg.SQL)

	redis := redis.Init(cfg.Redis)

	d := domain.Init(db, midtrans, redis)

	uc := usecase.Init(auth, d)

	r := rest.Init(cfg.Gin, configReader, uc, auth)

	r.Run()
}
