package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	internalCotacaoUrl = "http://localhost:8080/cotacao"
	fileName           = "cotacao.txt"
	requestTimeout     = 300
)

type Cotacao struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	result, err := buscaCotacao()
	if err != nil {
		log.Fatal(err)
	}

	err = saveToFile(result)
	if err != nil {
		log.Fatal(err)
	}
}

func buscaCotacao() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", internalCotacaoUrl, nil)
	if err != nil {
		return nil, errors.New("Erro ao criar request.\n" + err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao chamar a cotacao server.\n", err)
		return nil, errors.New("Erro ao chamar a cotacao server.\n" + err.Error())
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Erro ao ler resultado da cotacao.\n" + err.Error())
	}

	var cotacao Cotacao
	err = json.Unmarshal(result, &cotacao)
	if err != nil {
		return nil, err
	}

	log.Printf("Cotacao retornada com sucesso. \n%v", cotacao)
	return &cotacao, nil
}

func saveToFile(cotacao *Cotacao) error {
	dataToSave := fmt.Sprintf("Dólar: %v\n", cotacao.USDBRL.Bid)

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.New("error opening file: " + err.Error())
	}
	defer file.Close()

	if _, err := file.WriteString(dataToSave); err != nil {
		return errors.New("error writing to file: " + err.Error())
	}

	return nil
}
