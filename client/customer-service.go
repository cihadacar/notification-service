package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email"`
	Address string `json:"address,omitempty"`
}

func CustomerServiceRequest(customerId string) (Customer, error) {
	url := fmt.Sprintf("http://customer-service:8080/customers/%s", customerId)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("request gonderilmeli: %v", err)
		return Customer{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Response okunmadi: %v", err)
		return Customer{}, err
	}

	var customer Customer
	err = json.Unmarshal(body, &customer)

	if err != nil {
		log.Printf("customer parse edilemedi: %v", err)
		return Customer{}, err
	}

	return customer, nil
}
