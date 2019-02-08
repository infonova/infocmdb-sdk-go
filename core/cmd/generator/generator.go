package main

import (
	"fmt"
	"html/template"
	"os"
)

//go:generate go run generator.go

type webservice struct {
	name   string
	params []string
}

type webservices struct {
	webservice []webservice
}

var webservice_defaults = webservices{
	[]webservice{
		{name: "test", params: []string{"aa", "bbb"}},
		{name: "test1", params: []string{"aa", "bbb"}},
		{name: "test2", params: []string{"aa", "bbb"}},
		{name: "test3", params: []string{"aa", "bbb"}},
	},
}

type templatedata struct {
	Packagename string
	Webservices webservices
}

func main() {
	fmt.Println("Generating default webservices")
	webservice_template_file := "webservice_default.gotmpl"
	t := template.Must(template.New(webservice_template_file).ParseFiles(webservice_template_file))
	templatedata := &templatedata{Packagename: "WhatHWHATHATH", Webservices: webservice_defaults}
	err := t.Execute(os.Stdout, templatedata)
	if err != nil {
		panic(err)
	}
}
