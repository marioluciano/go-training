package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	contacatoUrl     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	requestTimeout   = 200
	databaseTimetout = 10

	dbDriver  = "sqlite3"
	dbName    = "cotacao.db"
	tableName = "cotacao_data"
)

type Cotacao struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	log.Println("Iniciando http server")
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", buscaCotacaoHandler)
	http.ListenAndServe(":8080", mux)
}

func buscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cotacao, err := buscaCotacao()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Iniciando conexão com banco de dados.")
	db, err := initDB()
	if err != nil {
		log.Fatalf("Error ao inicializar banco de dados: %v", err)
	}
	defer db.Close()

	err = salvaCotacao(db, cotacao)
	if err != nil {
		log.Fatalf("Error saving currency data to DB: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cotacao)
}

func buscaCotacao() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", contacatoUrl, nil)
	if err != nil {
		log.Println("Erro ao criar request.\n", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao fazer requisição para a API de cotação.\n", err)
		return nil, err
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cotacao Cotacao
	err = json.Unmarshal(result, &cotacao)
	if err != nil {
		return nil, err
	}

	log.Println("Requisição para API externa de cotação feito com sucesso.")
	return &cotacao, nil
}

func salvaCotacao(db *sql.DB, cotacato *Cotacao) error {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimetout*time.Millisecond)
	defer cancel()

	insertSQL := "INSERT INTO " + tableName + " (moeda, cotacao, timestamp) VALUES (?, ?, ?)"

	_, err := db.ExecContext(ctx, insertSQL, "USD", cotacato.USDBRL.Bid, time.Now())
	if err != nil {
		return fmt.Errorf("error inserting data into table: %w", err)
	}

	return nil
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbName)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS ` + tableName + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		moeda TEXT NOT NULL,
		cotacao DECIMAL NOT NULL,
		timestamp DATETIME NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	return db, nil
}
