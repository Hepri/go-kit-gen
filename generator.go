package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

const (
	outDir = "generated"
)

func processTemplateWithService(tplName string, s *Service) (string, error) {
	gopath := os.Getenv("GOPATH")
	t, err := template.ParseFiles(filepath.Join(gopath, "src", "github.com", "Hepri", "go-kit-gen", "templates", tplName))
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	err = t.Execute(buf, s)
	if err != nil {
		return "", err
	}

	//return buf.String(), nil

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

func generateFileFromTemplate(fileName string, tplName string, s *Service) error {
	txt, err := processTemplateWithService(tplName, s)
	if err != nil {
		return err
	}

	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	outPath := path.Join(workDir, outDir)
	if _, err = os.Stat(outPath); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(outPath, os.ModePerm)
		}
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(path.Join(outPath, fileName), []byte(txt), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func GenerateEndpoints(s *Service) error {
	return generateFileFromTemplate("endpoints.go", "endpoints.html", s)
}

func GenerateService(s *Service) error {
	return generateFileFromTemplate("service.go", "service.html", s)
}
