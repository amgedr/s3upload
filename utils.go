package main

import (
	"code.google.com/p/gcfg" //Ref: http://code.google.com/p/gcfg/
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

type Config struct {
	Auth struct {
		AccessKey string
		SecretKey string
	}

	Setup struct {
		Log       string
		Overwrite bool
	}

	Locations struct {
		Source      []string
		Destination []string
	}
}

//Get the configurations from ~/.s3upload.conf
func getConfigs() Config {
	usr, err := user.Current()
	if err != nil {
		logger.Fatalln(err.Error())
	}

	file := usr.HomeDir + string(os.PathSeparator) + ".s3upload.conf"

	//if the file does not exist create it
	if _, err := os.Stat(file); os.IsNotExist(err) {
		err := ioutil.WriteFile(file, []byte(newConfig), 644)
		if err != nil {
			//logger is not initialized yet so use Go's default log
			log.Fatalln("Could not create .s3upload.conf in the current user's home folder.")
		} else {
			log.Fatalln(".s3upload.conf was successfully created in the current user's home folder.",
				" Please edit it and enter correct settings.")
		}
	}

	var cfg Config
	cfg_err := gcfg.ReadFileInto(&cfg, file)
	if cfg_err != nil {
		fmt.Println(cfg_err.Error())
		os.Exit(1)
	}
	return cfg
}

var newConfig = `[Auth]
AccessKey="your_amazon_accesskey"  # Get from Amazon
SecretKey="your_amazon_secretkey"  # Get from Amazon

[Setup]
#Setting this option to true may incur additional cost
overwrite=false

[Locations]
source=""       # local directory to upload
destination=""  # S3 full path

# You can make s3upload upload multiple location by uncommenting the following
# lines and repeating them as much as you need.
# source=""
# destination=""
`
