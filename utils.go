package main

import (
	"code.google.com/p/gcfg" //Ref: http://code.google.com/p/gcfg/
	"fmt"
	"os"
	"os/user"
)

type Config struct {
	Auth struct {
		AccessKey string
		SecretKey string
	}
}

//Get the configurations from ~/.s3upload.conf
func getConfigs() Config {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(10)
	}

	var cfg Config
	cfg_err := gcfg.ReadFileInto(&cfg,
		usr.HomeDir+string(os.PathSeparator)+".s3upload.conf")
	if cfg_err != nil {
		fmt.Println(cfg_err)
		os.Exit(11)
	}
	return cfg
}
