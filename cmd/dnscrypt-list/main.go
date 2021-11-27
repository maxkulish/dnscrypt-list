package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/maxkulish/dnscrypt-list/lib/config"
	"github.com/maxkulish/dnscrypt-list/lib/db"
	"github.com/maxkulish/dnscrypt-list/lib/download"
	"github.com/maxkulish/dnscrypt-list/lib/files"
	"github.com/maxkulish/dnscrypt-list/lib/logger"
	"github.com/maxkulish/dnscrypt-list/lib/output"
	"github.com/maxkulish/dnscrypt-list/lib/target"
	"go.uber.org/zap"
	"os"
	"time"
)

// version contains the version of the application
var version = "0.0.0-src" // set via ldflags

func main() {

	start := time.Now()

	conf, err := config.Get()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	logger.Info("dnscrypt-list starting", zap.String("version", version))

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

	logger.Info("collecting remoteTargets ...")
	remoteTargets, err := target.CollectTargets(conf)
	if err != nil {
		logger.Error("remoteTargets error", zap.Error(err))
	}

	logger.Info("found remoteTargets", zap.Int("total", remoteTargets.Length()))

	localFiles, err := download.GetAndSaveTargets(conf.TempDir, remoteTargets)
	if err != nil {
		logger.Error("get and save remoteTargets error", zap.Error(err))
	}

	// Read files and save them to the whitelist db
	err = download.ReadFilesAndSaveToDB(localFiles, whitelist, target.WhiteList)
	if err != nil {
		logger.Error("whitelist read and save error", zap.Error(err))
	}

	// Read files and save them to the blacklist db
	err = download.ReadFilesAndSaveToDB(localFiles, blacklist, target.BlackList)
	if err != nil {
		logger.Error("blacklist read and save error", zap.Error(err))
	}

	var domainCounter int64 // count how many domain collected
	// Read whitelist from the DB and save to the output file
	keys := whitelist.GetAllKeys()
	domainCounter += int64(len(keys))

	err = output.SaveDomainToFile(conf.Output.Whitelist, keys)
	if err != nil {
		logger.Error("save domains to the file error", zap.Error(err))
	}

	// Read blacklist from the DB and save to the output file
	keys = blacklist.GetAllKeys()
	domainCounter += int64(len(keys))

	err = output.SaveDomainToFile(conf.Output.Blacklist, keys)
	if err != nil {
		logger.Error("save domains to the file error", zap.Error(err))
	}

	var toDelete []string
	for _, tmpFile := range localFiles {
		if tmpFile.Temp {
			toDelete = append(toDelete, tmpFile.Path)
		}
	}

	err = files.DeleteAllFiles(toDelete...)
	if err != nil {
		logger.Debug("temporary files deletion error", zap.Error(err))
	}

	logger.Info(
		fmt.Sprintf("elapsed: %.2fs", time.Since(start).Seconds()),
		zap.Int("files", len(localFiles)),
		zap.String("domains", humanize.Comma(domainCounter)))
}
