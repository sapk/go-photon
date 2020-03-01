package main

import (
	"archive/zip"
	"fmt"
	"strings"

	"gopkg.in/ini.v1"
)

type SL1File struct {
	Config   *ini.File
	FilePath string
	Layers   int
}

func ReadSL1File(filepath string) (*SL1File, error) {

	sl1File := &SL1File{
		FilePath: filepath,
	}
	r, err := zip.OpenReader(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SL1 file: %v", err)
	}
	defer r.Close()

	//Count layers
	for _, f := range r.File {
		//fmt.Printf("Contents of %s:\n", f.Name)
		if f.Name == "config.ini" {
			//Read config file
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open SL1 config file: %v", err)
			}
			cfg, err := ini.Load(rc)
			if err != nil {
				return nil, fmt.Errorf("failed to read SL1 config file: %v", err)
			}
			sl1File.Config = cfg

		} else if strings.HasSuffix(f.Name, ".png") {
			//Count layers
			layerNumber++
		}
	}
	sl1File.Layers = layerNumber

	return sl1File, nil
}
