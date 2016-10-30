package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kr/s3/s3util" //Ref: http://godoc.org/github.com/kr/s3/s3util
)

var (
	cfg         Config
	locationPtr int //contains the index of the current location being processed
	waiter      sync.WaitGroup
	infoLog     *log.Logger
	errLog      *log.Logger

	summary struct {
		Success int
		Fails   int
		Exists  int
	}
)

func main() {
	cfg = getConfigs()

	//log errors to the path in the config file (if opened successfully) and stderr
	logfile, err := os.OpenFile(cfg.Setup.Log, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		infoLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
		errLog = log.New(os.Stderr, "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		infoLog = log.New(io.MultiWriter(logfile, os.Stdout), "", log.Ldate|log.Ltime)
		errLog = log.New(io.MultiWriter(logfile, os.Stderr), "Error: ",
			log.Ldate|log.Ltime|log.Lshortfile)
	}
	defer logfile.Close()

	if len(cfg.Locations.Source) != len(cfg.Locations.Destination) {
		errLog.Fatalln("The number of sources and destinations should be the same.")
	}

	if cfg.Auth.AccessKey != "" && cfg.Auth.SecretKey != "" {
		s3util.DefaultConfig.AccessKey = cfg.Auth.AccessKey
		s3util.DefaultConfig.SecretKey = cfg.Auth.SecretKey
	} else {
		errLog.Fatalln("AccessKey and/or SecretKey not set.")
	}

	for i := 0; i < len(cfg.Locations.Destination); i++ {
		//remove trailing slashes to unify the way paths are combined
		//ie: wether they are entered in the config file or not it still works
		cfg.Locations.Destination[i] = strings.Trim(cfg.Locations.Destination[i], "/")
	}

	//loop through all the source paths to check if they exist before processing
	for _, src := range cfg.Locations.Source {
		if _, err := os.Stat(src); os.IsNotExist(err) {
			errLog.Fatalln(src + " does not exist.")
		}
	}

	for i, src := range cfg.Locations.Source {
		locationPtr = i
		_ = filepath.Walk(src, upload)
	}

	waiter.Wait()

	infoLog.Printf("Successful: %d, Existed: %d, Failed: %d\n\n",
		summary.Success, summary.Exists, summary.Fails)
}
