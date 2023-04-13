package api

import (
	"github.com/gin-gonic/gin"
	"github.com/qwddz/bitmexws/internal"
	"github.com/qwddz/bitmexws/internal/config"
	"github.com/qwddz/bitmexws/internal/statistics"
	"github.com/qwddz/bitmexws/pkg/logger"
	"github.com/qwddz/bitmexws/pkg/store"
	"github.com/sirupsen/logrus"
	"log"
)

type API struct {
	config *config.Config
	router *gin.Engine
	log    logger.Logger
	store  *store.Store

	statistic internal.Router
}

func New() (*API, error) {
	a := API{
		config: config.NewConfig(),
		router: gin.Default(),
		log:    logrus.New(),
	}

	if err := a.setStore(); err != nil {
		return nil, err
	}

	a.setRoutes()

	a.log.Infoln("setup api log: application has been successfully configured")

	return &a, nil
}

func (a *API) ServeTCP() error {
	a.log.Infoln("setup api log: application has been successfully started")

	rg := a.router.Group("statistics")
	{
		a.statistic.RegisterRouter(rg)
	}

	return a.router.Run(a.config.ApiConfig.BindAddr)
}

func (a *API) Shutdown() error {
	err := a.store.Close()
	if err != nil {
		return err
	}

	log.Println("setup api log: store connection closed...")

	return nil
}

func (a *API) setRoutes() {
	a.log.Infoln("setup api log: setup router")

	a.statistic = statistics.NewRouter(a.store, a.log)
}

func (a *API) setStore() error {
	a.log.Infoln("setup api log: setup client db connection")

	conf := store.Config{
		Host: store.Host{
			Master: a.config.DB.Host.Master,
			Slave:  a.config.DB.Host.Slave,
		},
		Name:     a.config.DB.Name,
		User:     a.config.DB.User,
		Password: a.config.DB.Password,
	}

	st, err := store.New(conf)
	if err != nil {
		return err
	}

	a.store = st

	return nil
}
