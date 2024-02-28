package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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

// receive-only: diz que esse canal somente recebe dados (seta do lado direito)
func BuscarEnderecoViaBrasilAPI(url string, ch chan<- *EnderecoBrasilAPI) error {
	client := http.Client{Timeout: time.Second}
	resp, err := client.Get(url)
	if err != nil {
		ch <- nil
		return err
	}
	defer resp.Body.Close()

	var endereco EnderecoBrasilAPI
	if err := json.NewDecoder(resp.Body).Decode(&endereco); err != nil {
		ch <- nil
		return err
	}

	ch <- &endereco
	return nil
}

// receive-only: diz que esse canal somente recebe dados (seta do lado direito)
func BuscarEnderecoViaCepAPI(url string, ch chan<- *EnderecoViaCEP) error {
	client := http.Client{Timeout: time.Second}
	resp, err := client.Get(url)
	if err != nil {
		ch <- nil
		return err
	}
	defer resp.Body.Close()

	var endereco EnderecoViaCEP
	if err := json.NewDecoder(resp.Body).Decode(&endereco); err != nil {
		ch <- nil
		return err
	}

	ch <- &endereco
	return nil
}

func main() {

	// CEP para consulta
	cep := "01153000"

	canalBrasil := make(chan *EnderecoBrasilAPI)
	canalViaCEP := make(chan *EnderecoViaCEP)

	// As duas requisições serão feitas simultaneamente para as seguintes APIs:
	urlBrasilAPI := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	urlViaCepAPI := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	go BuscarEnderecoViaCepAPI(urlViaCepAPI, canalViaCEP)
	go BuscarEnderecoViaBrasilAPI(urlBrasilAPI, canalBrasil)

	var brasilAPI *EnderecoBrasilAPI
	var viaCepAPI *EnderecoViaCEP

	var vencedora string
	var resultado interface{}

	// selecionando a api vencedora
	select {
	case brasilAPI = <-canalBrasil:
		if brasilAPI != nil {
			vencedora = "BrasilAPI"
			resultado = brasilAPI
		} else {
			fmt.Println("Falha ao tentar buscar o endereço via BrasilAPI")
		}
	case viaCepAPI = <-canalViaCEP:
		if viaCepAPI != nil {
			vencedora = "ViaCEP"
			resultado = viaCepAPI
		} else {
			fmt.Println("Falha ao tentar buscar o endereço via ViaCepAPI")
		}
	case <-time.After(time.Second):
		fmt.Println("Timeout: nenhuma resposta foi recebida em 1 segundo.")
	}

	// mostrar o resultado da vencedora
	if vencedora != "" {
		fmt.Printf("API vencedora: %s\n", vencedora)
		fmt.Println("Resultado:")
		jsonData, err := json.Marshal(resultado)
		if err != nil {
			fmt.Println("Falha ao tentar serializar o resultado da busca:", err)
		} else {
			fmt.Println(string(jsonData))
		}
	}
}
