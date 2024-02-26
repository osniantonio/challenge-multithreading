package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Endereco interface{}

type EnderecoBrasilAPI struct {
	Cep          string `json:"cep"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Service      string `json:"service"`
	State        string `json:"state"`
	Street       string `json:"street"`
}

type EnderecoViaCEP struct {
	Bairro      string `json:"bairro"`
	Cep         string `json:"cep"`
	Complemento string `json:"complemento"`
	DDD         string `json:"ddd"`
	GIA         string `json:"gia"`
	IBGE        string `json:"ibge"`
	Localidade  string `json:"localidade"`
	Logradouro  string `json:"logradouro"`
	Siafi       string `json:"siafi"`
	UF          string `json:"uf"`
}

func BuscarCEP(url string, ch chan<- *Endereco) error {
	startTime := time.Now()

	client := http.Client{Timeout: time.Second}
	resp, err := client.Get(url)
	if err != nil {
		ch <- nil
		return err
	}
	defer resp.Body.Close()

	var endereco Endereco
	if err := json.NewDecoder(resp.Body).Decode(&endereco); err != nil {
		ch <- nil
		return err
	}

	elapsedTime := time.Since(startTime)
	fmt.Printf("Tempo de execução (%s): %v\n", url, elapsedTime)

	ch <- &endereco

	return nil
}

func main() {

	cep := "01153000"

	ch := make(chan *Endereco)

	urlBrasilAPI := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	go BuscarCEP(urlBrasilAPI, ch)

	urlViaCepAPI := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	go BuscarCEP(urlViaCepAPI, ch)

	var api1, api2 *Endereco

	select {
	case api1 = <-ch:
		fmt.Println("Resposta da API 1 (brasilapi.com.br):")
		if api1 != nil {
			jsonData, _ := json.Marshal(api1)
			fmt.Println(string(jsonData))
		} else {
			fmt.Println("Erro ao buscar endereço na API 1")
		}
	case api2 = <-ch:
		fmt.Println("Resposta da API 2 (viacep.com.br):")
		if api2 != nil {
			jsonData, _ := json.Marshal(api2)
			fmt.Println(string(jsonData))
		} else {
			fmt.Println("Erro ao buscar endereço na API 2")
		}
	case <-time.After(time.Second):
		fmt.Println("Timeout: Nenhuma resposta recebida em 1 segundo.")
	}

}
