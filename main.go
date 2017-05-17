package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func generateServicePart(fn func(s *Service) error, s *Service, err *error) {
	if *err == nil {
		*err = fn(s)
	}
}

// Usage is a replacement usage function for the flags package
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s [service_file]\n", os.Args[0])
}

func main() {
	binaryName := os.Args[0]
	fileName := os.Args[1]
	flag.Usage = Usage
	flag.Parse()

	log.SetPrefix(fmt.Sprintf("%s:", binaryName))

	fil, err := parser.ParseFile(token.NewFileSet(), fileName, nil, parser.AllErrors)
	if err != nil {
		log.Println(err)
		return
	}

	srv, err := ParseService(fil)
	if err != nil {
		log.Println(err)
		return
	}

	generateServicePart(GenerateEndpoints, srv, &err)
	generateServicePart(GenerateService, srv, &err)

	if err != nil {
		log.Println(err)
		return
	}
}
