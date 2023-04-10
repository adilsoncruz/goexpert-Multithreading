package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"
)

type Endereco struct {
	Cep        string `json:"cep,code"`
	Logradouro string `json:"logradouro,address"`
	Cidade     string `json:"localidade,city"`
	Estado     string `json:"uf,state"`
	Bairro     string `json:"bairro,district"`
}

type EnderecoViaCep struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Cidade     string `json:"localidade"`
	Estado     string `json:"uf"`
	Bairro     string `json:"bairro"`
}

type EnderecoAPICep struct {
	Cep        string `json:"code"`
	Logradouro string `json:"address"`
	Cidade     string `json:"city"`
	Estado     string `json:"state"`
	Bairro     string `json:"district"`
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

func requestCep[T any](ch chan T, serverName string, cep string) {
	tmpl := getTemplate(serverName)
	var serverRoute bytes.Buffer
	err := tmpl.Execute(&serverRoute, Busca{Cep: cep})
	if err != nil {
		panic(err)
	}
	req, err := http.Get(serverRoute.String())
	if err != nil || req.StatusCode != 200 {
		panic(errors.New("Falha na requisição"))
	}
	fmt.Println("The status code we got is:", req.StatusCode)
	fmt.Println("The status code text we got is:", http.StatusText(req.StatusCode))
	res, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	req.Body.Close()
	var msg T
	err = json.Unmarshal(res, &msg)
	if err != nil {
		panic(err)
	}

	ch <- msg
}

func main() {
	cApiCep := make(chan EnderecoAPICep)
	cViaCep := make(chan EnderecoViaCep)

	go requestCep(cViaCep, "viaCep", "05734-080")
	go requestCep(cApiCep, "apiCep", "05734-080")

	select {
	case msg1 := <-cApiCep:
		fmt.Printf("ApiCep return: %+v\n", msg1)
	case msg2 := <-cViaCep:
		fmt.Printf("ViaCep return: %+v\n", msg2)
	case <-time.After(time.Second * 1):
		println("timeout")
	}
}
