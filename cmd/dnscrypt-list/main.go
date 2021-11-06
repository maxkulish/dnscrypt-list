package main

import (
	"fmt"
	"github.com/maxkulish/dnscrypt-list/lib/config"
	"github.com/maxkulish/dnscrypt-list/lib/db"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"github.com/maxkulish/dnscrypt-list/lib/target"
	"go.uber.org/zap"
	"os"
)

// BuildVersion contains the version of the application
var BuildVersion string

func main() {
	conf, err := config.Get()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Info("dnscypt-list starting", zap.String("version", BuildVersion))

	conn, err := db.NewConn(conf.TempDB)
	if err != nil {
		logger.Error("database connection error", zap.Error(err))
	}
	defer conn.Close()

	logger.Info("collecting targets ...")
	targets, err := target.CollectTargets(conf)
	if err != nil {
		logger.Error("targets error", zap.Error(err))
	}

	fmt.Println("Found targets:", targets.Length())

}
