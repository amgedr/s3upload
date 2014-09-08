package main

import (
	"code.google.com/p/gcfg" //Ref: http://code.google.com/p/gcfg/
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

type Config struct {
	Auth struct {
		AccessKey string
		SecretKey string
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
		fmt.Println(err)
		os.Exit(1)
	}

	file := usr.HomeDir + string(os.PathSeparator) + ".s3upload.conf"
	var cfg Config
	cfg_err := gcfg.ReadFileInto(&cfg, file)
	if cfg_err != nil {
		err := ioutil.WriteFile(file, []byte(newConfig), 644)
		if err != nil {
			fmt.Println("Error: Could not create .s3upload.conf in the current user's home folder.",
				" Please check the permission and run again.")
		} else {
			fmt.Println(".s3upload.conf was successfully created in the current user's home folder.",
				" Please edit it and enter correct settings.")
		}
		os.Exit(1)
	}
	return cfg
}

var newConfig = `[Auth]
AccessKey="your_amazon_accesskey"  # Get from Amazon
SecretKey="your_amazon_secretkey"  # Get from Amazon

[Locations]
source=""       # local directory to upload
destination=""  # S3 full path

# You can make s3upload upload multiple location by uncommenting the following
# lines and repeating them as much as you need.
# source=""
# destination=""
`
