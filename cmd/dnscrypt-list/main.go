package main

import (
	"fmt"
	"github.com/maxkulish/dnscrypt-list/lib/config"
	"github.com/maxkulish/dnscrypt-list/lib/db"
	"github.com/maxkulish/dnscrypt-list/lib/download"
	"github.com/maxkulish/dnscrypt-list/lib/files"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"github.com/maxkulish/dnscrypt-list/lib/output"
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

	// create whitelist db
	whitelist, err := db.NewConn(conf.WhiteListDB)
	if err != nil {
		logger.Error("whitelist database connection error", zap.Error(err))
	}
	defer whitelist.Close()

	// create blacklist db
	blacklist, err := db.NewConn(conf.BlackListDB)
	if err != nil {
		logger.Error("blacklist database connection error", zap.Error(err))
	}
	defer blacklist.Close()

	logger.Info("collecting targets ...")
	targets, err := target.CollectTargets(conf)
	if err != nil {
		logger.Error("targets error", zap.Error(err))
	}

	logger.Info("found targets", zap.Int("total", targets.Length()))

	tempFiles, err := download.GetAndSaveTargets(conf.TempDir, targets)
	if err != nil {
		logger.Error("get and save targets error", zap.Error(err))
	}

	// Read files and save them to the whitelist db
	err = download.ReadFilesAndSaveToDB(tempFiles, whitelist, target.WhiteList)
	if err != nil {
		logger.Error("whitelist read and save error", zap.Error(err))
	}

	// Read files and save them to the blacklist db
	err = download.ReadFilesAndSaveToDB(tempFiles, blacklist, target.BlackList)
	if err != nil {
		logger.Error("blacklist read and save error", zap.Error(err))
	}

	// Read whitelist from the DB and save to the output file
	keys := whitelist.GetAllKeys()

	err = output.SaveDomainToFile(conf.Output.Whitelist, keys)
	if err != nil {
		logger.Error("save domains to the file error", zap.Error(err))
	}

	// Read blacklist from the DB and save to the output file
	keys = blacklist.GetAllKeys()

	err = output.SaveDomainToFile(conf.Output.Blacklist, keys)
	if err != nil {
		logger.Error("save domains to the file error", zap.Error(err))
	}

	err = files.DeleteAllFiles(tempFiles...)
	if err != nil {
		logger.Debug("temporary files deletion error", zap.Error(err))
	}
}
