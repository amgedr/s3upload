package main

import (
	"fmt"
	"github.com/kr/s3/s3util" //Ref: http://godoc.org/github.com/kr/s3/s3util
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	cfg         Config
	locationPtr int //contains the index of the current location being processed
)

func main() {
	cfg = getConfigs()

	if len(cfg.Locations.Source) != len(cfg.Locations.Destination) {
		fmt.Println("Error: The number of sources and destinations should be the same.")
		os.Exit(1)
	}

	count := len(cfg.Locations.Source)
	s3util.DefaultConfig.AccessKey = cfg.Auth.AccessKey
	s3util.DefaultConfig.SecretKey = cfg.Auth.SecretKey

	//loop through all the source paths to check if they exist before processing
	for i := 0; i < count; i++ {
		if _, err := os.Stat(cfg.Locations.Source[i]); os.IsNotExist(err) {
			fmt.Println(cfg.Locations.Source[i] + " does not exist.")
			os.Exit(1)
		}
	}

	for i := 0; i < count; i++ {
		src := cfg.Locations.Source[i]
		locationPtr = i

		_ = filepath.Walk(src, upload)
	}
}

func upload(path string, f os.FileInfo, err error) error {
	if !f.IsDir() { //if path is a file
		uploadFile(path)
	}

	return nil
}

func uploadFile(path string) {
	dest := strings.Replace(path, cfg.Locations.Source[locationPtr], "", 1)
	dest = strings.Trim(dest, "\\")             //in case its a Windows path
	dest = strings.Replace(dest, "\\", "/", -1) //in case its a Windows path

	fmt.Print(path)

	r, rerr := os.Open(path)
	if rerr != nil {
		fmt.Println("....." + rerr.Error())
		return
	}
	defer r.Close()

	w, werr := s3util.Create(cfg.Locations.Destination[locationPtr]+dest, nil, nil)
	if werr != nil {
		fmt.Println("....." + werr.Error())
		return
	}
	defer w.Close()

	io.Copy(w, r)

	fmt.Println(".....Done")
}

func checkFile(file, dest string) bool {
	r, err := s3util.Open(dest+file, nil)
	if err != nil {
		return false
	}
	defer r.Close()

	return true
}
