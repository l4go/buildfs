package main

import (
	"log"
	"text/template"
	"time"
	"os"
	"path"
)

var OutputFile = "./buildtime_test.go"
var TmplFile = "./buildtime_test.go.tmpl"
type TmplVar struct {
	Name string
	BuildMicroSeconds int64
}

func main() {
	log.Println("Start: generation")
	defer log.Println("Done: generation")

	tmpl, err := template.ParseFiles(TmplFile)
	if err != nil {
		log.Fatal(err)
	}

	out, err := os.OpenFile(OutputFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	param := TmplVar{
		Name: path.Base(TmplFile),
		BuildMicroSeconds: time.Now().UnixMicro(),
	}

	tmpl.Execute(out, param)
}
