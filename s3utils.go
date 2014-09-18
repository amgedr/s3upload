package main

import (
	"fmt"
	"github.com/kr/s3/s3util" //Ref: http://godoc.org/github.com/kr/s3/s3util
	"io"
	"os"
	"strings"
)

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

	fullPath := cfg.Locations.Destination[locationPtr] + "/" + dest

	if !cfg.Setup.Overwrite {
		//if file exists do not upload again
		if _, e := s3util.Open(fullPath, nil); e == nil {
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

	w, err := s3util.Create(fullPath, nil, nil)
	if err != nil {
		logger.Printf("%s...%s\n", dest, err.Error())
		return
	}
	defer w.Close()

	io.Copy(w, r)
	fmt.Printf("%s.....Done\n", dest)
}
