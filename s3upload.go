package main

import (
	"fmt"
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

func upload(path string, f os.FileInfo, err error) error {
	if !f.IsDir() { //if path is a file
		waiter.Add(1)
		go uploadFile(path)
	}

	return nil
}

func uploadFile(path string) {
	dest := strings.Replace(path, cfg.Locations.Source[locationPtr], "", 1)
	dest = strings.Trim(dest, "\\/")            //remove *nix and Windows dir separators
	dest = strings.Replace(dest, "\\", "/", -1) //in case its a Windows path

	defer waiter.Done()

	if !cfg.Setup.Overwrite {
		//if file exists do not upload again
		if _, e := s3util.Open(cfg.Locations.Destination[locationPtr]+"/"+dest, nil); e == nil {
			fmt.Printf("%s...Already exists.\n", dest)
			return
		}
	}

	r, err := os.Open(path)
	if err != nil {
		logger.Printf("%s...%s", dest, err.Error())
		return
	}
	defer r.Close()

	w, err := s3util.Create(cfg.Locations.Destination[locationPtr]+"/"+dest, nil, nil)
	if err != nil {
		logger.Printf("%s...%s\n", dest, err.Error())
		return
	}
	defer w.Close()

	io.Copy(w, r)
	fmt.Printf("%s.....Done\n", dest)
}

func checkFile(file, dest string) bool {
	r, err := s3util.Open(dest+file, nil)
	if err != nil {
		return false
	}
	defer r.Close()

	return true
}
