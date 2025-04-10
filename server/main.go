package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/luisteixeira74/go-expert-desafio-multithreading/model"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/consulta-cep/{cep}", consultaCepHandler).Methods("GET")

	fmt.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func consultaCepHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cep := vars["cep"]

	cepResult, err := consultaCep(cep)
	if err != nil {
		http.Error(w, "CEP não encontrado ou timeout", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cepResult)
}

func consultaCep(cep string) (model.Cep, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cepChan := make(chan model.Cep)

	go func() {
		log.Println("Consultando ViaCEP")
		if result, err := consultaViaCep(cep); err == nil {
			select {
			case cepChan <- result:
			case <-ctx.Done():
			}
		} else {
			log.Println("Erro ViaCEP:", err)
		}
	}()

	go func() {
		log.Println("Consultando BrasilAPI")
		if result, err := consultaBrasilApi(cep); err == nil {
			select {
			case cepChan <- result:
			case <-ctx.Done():
			}
		} else {
			log.Println("Erro BrasilAPI:", err)
		}
	}()

	select {
	case res := <-cepChan:
		log.Println("Resposta recebida")
		return res, nil
	case <-ctx.Done():
		log.Println("Timeout atingido")
		return model.Cep{}, fmt.Errorf("timeout ao consultar CEP")
	}
}

func consultaViaCep(cep string) (model.Cep, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return model.Cep{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Cep{}, fmt.Errorf("erro ao consultar ViaCEP: %s", resp.Status)
	}

	var viaCepResponse model.ViaCepResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCepResponse); err != nil {
		return model.Cep{}, err
	}

	cepReturn := model.Cep{
		Cep:        viaCepResponse.Cep,
		Logradouro: viaCepResponse.Logradouro,
		Bairro:     viaCepResponse.Bairro,
		Cidade:     viaCepResponse.Localidade,
		Estado:     viaCepResponse.Uf,
		Sender:     "ViaCEP",
	}

	return cepReturn, nil
}

func consultaBrasilApi(cep string) (model.Cep, error) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := http.Get(url)
	if err != nil {
		return model.Cep{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Cep{}, fmt.Errorf("erro ao consultar BrasilAPI: %s", resp.Status)
	}

	var brasilApiResponse model.BrasilApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&brasilApiResponse); err != nil {
		return model.Cep{}, err
	}

	cepReturn := model.Cep{
		Cep:        brasilApiResponse.Cep,
		Logradouro: brasilApiResponse.Street,
		Bairro:     brasilApiResponse.Neighborhood,
		Cidade:     brasilApiResponse.City,
		Estado:     brasilApiResponse.State,
		Sender:     "BrasilAPI",
	}

	return cepReturn, nil
}
