package main

import (
	"github.com/kr/s3/s3util" //Ref: http://godoc.org/github.com/kr/s3/s3util
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	cfg         Config
	locationPtr int //contains the index of the current location being processed
	waiter      sync.WaitGroup
	logger      *log.Logger
)

func main() {
	cfg = getConfigs()

	//log errors to the path in the config file (if opened successfully) and stderr
	logfile, err := os.OpenFile(cfg.Setup.Log, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logger = log.New(io.MultiWriter(logfile, os.Stderr), "",
			log.Ldate|log.Ltime|log.Lshortfile)
	}
	defer logfile.Close()

	if len(cfg.Locations.Source) != len(cfg.Locations.Destination) {
		logger.Fatalln("Error: The number of sources and destinations should be the same.")
	}

	if cfg.Auth.AccessKey != "" && cfg.Auth.SecretKey != "" {
		s3util.DefaultConfig.AccessKey = cfg.Auth.AccessKey
		s3util.DefaultConfig.SecretKey = cfg.Auth.SecretKey
	} else {
		logger.Fatalln("Error: AccessKey and/or SecretKey not set.")
	}

	for i := 0; i < len(cfg.Locations.Destination); i++ {
		//remove trailing slashes to unify the way paths are combined
		//ie: wether they are entered in the config file or not it still works
		cfg.Locations.Destination[i] = strings.Trim(cfg.Locations.Destination[i], "/")
	}

	//loop through all the source paths to check if they exist before processing
	for _, src := range cfg.Locations.Source {
		if _, err := os.Stat(src); os.IsNotExist(err) {
			logger.Fatalln(src + " does not exist.")
		}
	}

	for i, src := range cfg.Locations.Source {
		locationPtr = i
		_ = filepath.Walk(src, upload)
	}

	waiter.Wait()
}
