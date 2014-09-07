package main

import (
	"fmt"
	"github.com/kr/s3/s3util" //Ref: http://godoc.org/github.com/kr/s3/s3util
	"io"
	"os"
	"path/filepath"
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
	for i := 0; i < count; i++ {
		src := cfg.Locations.Source[i]
		locationPtr = i

		_ = filepath.Walk(src, upload)
	}
}

func upload(path string, f os.FileInfo, err error) error {
	if !f.IsDir() { //if path is a file
		uploadFile(path, f.Name())
	} else { //else its a directory
		//TODO: create the directory
	}

	return nil
}

func uploadFile(file, name string) {
	s3util.DefaultConfig.AccessKey = cfg.Auth.AccessKey
	s3util.DefaultConfig.SecretKey = cfg.Auth.SecretKey

	r, rerr := os.Open(file)
	if rerr != nil {
		fmt.Println(rerr)
		os.Exit(2)
	}
	defer r.Close()

	w, werr := s3util.Create(cfg.Locations.Destination[locationPtr]+name, nil, nil)
	if werr != nil {
		fmt.Println(werr)
		os.Exit(3)
	}
	defer w.Close()

	io.Copy(w, r)
}

func checkFile(file, dest string) bool {
	s3util.DefaultConfig.AccessKey = cfg.Auth.AccessKey
	s3util.DefaultConfig.SecretKey = cfg.Auth.SecretKey

	r, err := s3util.Open(dest+file, nil)
	if err != nil {
		return false
	}
	defer r.Close()

	return true
}
