package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	cep := "01153000"

	// Channels to receive the results
	ch1 := make(chan string)
	ch2 := make(chan string)

	// Anonymous function for BrasilAPI
	go func() {
		url := "https://brasilapi.com.br/api/cep/v1/" + cep
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			ch1 <- fmt.Sprintf("BrasilAPI: Erro criado request: %v", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ch1 <- fmt.Sprintf("BrasilAPI: Erro no request: %v", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ch1 <- fmt.Sprintf("BrasilAPI: Erro lendo resposta: %v", err)
			return
		}
		ch1 <- fmt.Sprintf("BrasilAPI: %s", body)
	}()

	// Anonymous function for ViaCEP
	go func() {
		url := "http://viacep.com.br/ws/" + cep + "/json/"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			ch2 <- fmt.Sprintf("ViaCEP:  Erro criado request: %v", err)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ch2 <- fmt.Sprintf("ViaCEP: Erro no request: %v", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ch2 <- fmt.Sprintf("ViaCEP: Erro lendo resposta: %v", err)
			return
		}
		ch2 <- fmt.Sprintf("ViaCEP: %s", body)
	}()

	start := time.Now()
	// Select the first response or timeout                                                                                                                          │
	select {
	case res := <-ch1:
		fmt.Println("Respondeu antes: BrasilAPI")
		fmt.Printf("Em: %v\n", time.Since(start))
		fmt.Println(res)
	case res := <-ch2:
		fmt.Println("Respondeu antes: ViaCEP")
		fmt.Printf("Em: %v\n", time.Since(start))
		fmt.Println(res)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: Nenhuma API respondeu dentro de 1 segundo.")
	}
}
