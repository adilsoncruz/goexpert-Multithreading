package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"
)

type Endereco struct {
	Cep        string `json:"cep" json:"code"`
	Logradouro string `json:"logradouro" json:"address"`
	Cidade     string `json:"localidade" json:"city"`
	Estado     string `json:"uf" json:"state"`
	Bairro     string `json:"bairro" json:"district"`
}

var m = map[string]string{
	"apiCep": "https://cdn.apicep.com/file/apicep/{{.Cep}}.json",
	"viaCep": "http://viacep.com.br/ws/{{.Cep}}/json/",
}

type Busca struct {
	Cep string
}

func getTemplate(serverName string) *template.Template {
	serverRoute, exists := m[serverName]
	if exists != true {
		panic("not found server name")
	}
	tmpl, err := template.New(serverName).Parse(serverRoute)
	if err != nil {
		panic(err)
	}

	return tmpl
}

func requestCep(ch chan Endereco, serverName string, cep string) {
	tmpl := getTemplate(serverName)
	var serverRoute bytes.Buffer
	err := tmpl.Execute(&serverRoute, Busca{Cep: cep})
	if err != nil {
		panic(err)
	}
	req, err := http.Get(serverRoute.String())
	if err != nil {
		panic(err)
	}
	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	req.Body.Close()
	var msg Endereco
	err = json.Unmarshal(res, &msg)
	if err != nil {
		panic(err)
	}

	ch <- msg
}

func main() {
	cApiCep := make(chan Endereco)
	cViaCep := make(chan Endereco)

	go requestCep(cApiCep, "apiCep", "04094-000")
	go requestCep(cViaCep, "viaCep", "04094-000")

	select {
	case msg1 := <-cApiCep:
		fmt.Printf("ApiCep return: %+v\n", msg1)
	case msg2 := <-cViaCep:
		fmt.Printf("ViaCep return: %+v\n", msg2)
	case <-time.After(time.Second * 1):
		println("timeout")
	}
}
