package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type address struct {
	API          string `json:"api"`
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type brazilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type viaCEPResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"uf"`
	City         string `json:"localidade"`
	Neighborhood string `json:"bairro"`
	Street       string `json:"logradouro"`
}

func findAddressFromBrazilAPI(cep string, ch chan<- address) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var data brazilAPIResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return
	}

	ch <- address{
		API:          "brasilapi",
		Cep:          data.Cep,
		State:        data.State,
		City:         data.City,
		Neighborhood: data.Neighborhood,
		Street:       data.Street,
	}
}

func findAddressFromviaCEP(cep string, ch chan<- address) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)

	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var data viaCEPResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return
	}

	ch <- address{
		API:          "viacep",
		Cep:          data.Cep,
		State:        data.State,
		City:         data.City,
		Neighborhood: data.Neighborhood,
		Street:       data.Street,
	}
}

func main() {
	ch := make(chan address)
	cep := "01153000"

	go findAddressFromBrazilAPI(cep, ch)
	go findAddressFromviaCEP(cep, ch)

	select {
	case address := <-ch:
		fmt.Println("---------------------------------------")
		fmt.Printf("API: ** %s **\n", address.API)
		fmt.Printf("CEP: %s\n", address.Cep)
		fmt.Printf("State: %s\n", address.State)
		fmt.Printf("City: %s\n", address.City)
		fmt.Printf("Neighborhood: %s\n", address.Neighborhood)
		fmt.Printf("Street: %s\n", address.Street)
		fmt.Println("---------------------------------------")
	case <-time.After(1 * time.Second):
		fmt.Println("one second timeout exceeded!")
	}
}
