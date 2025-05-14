package app

import (
	"context"
	"log"

	"currency/internal/adapters/client/coindesk"
	"currency/internal/adapters/storage/postgres"
	"currency/internal/ports/http/public"
	"currency/internal/usecases"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

type Config struct {
	port          string
	connStr       string
	url           string
	baseUrlParams []string
}

func NewConfig() *Config {
	port := viper.GetString("port")
	connStr := viper.GetString("database.connStr")
	url := viper.GetString("externalAPI.url")
	baseUrlParams := viper.GetStringSlice("externalAPI.baseUrlParams.fsyms")

	return &Config{port: port, connStr: connStr, url: url, baseUrlParams: baseUrlParams}
}

func Run() error {
	viper.AddConfigPath("deployment/config")
	viper.SetConfigName("config")

	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "read config failed")
	}

	config := NewConfig()

	ctx := context.Background()

	storage, err := postgres.NewStorage(ctx, config.connStr)
	if err != nil {
		return errors.Wrap(err, "create storage failed")
	}

	client, err := coindesk.NewClient(config.url)
	if err != nil {
		return errors.Wrap(err, "create client failed")
	}

	service, err := usecases.NewService(storage, client)
	if err != nil {
		return errors.Wrap(err, "create service failed")
	}

	server, err := public.NewServer(service, config.port)
	if err != nil {
		return errors.Wrap(err, "create server failed")
	}

	go runCrone(service, config.baseUrlParams)

	err = server.Run()
	if err != nil {
		return errors.Wrap(err, "server run failed")
	}

	return nil
}

func runCrone(service *usecases.Service, titles []string) {
	ctx := context.Background()

	_, err := service.GetCoinsFromAPI(ctx, titles...)
	if err != nil {
		log.Println(err)
	}

	c := cron.New()
	updateFunc := func() {
		_, err := service.GetCoinsFromAPI(ctx)
		if err != nil {
			log.Println(err)
		}
	}
	c.AddFunc("@every 1m", updateFunc)
	c.Run()
}
