package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type Response struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

var statuses = []string{"PROCESSING", "PROCESSED"} // "INVALID"

func handler(w http.ResponseWriter, r *http.Request) {
	orderNumber := r.URL.Path[len("/api/orders/"):]
	if orderNumber == "" {
		http.Error(w, "Order number is required", http.StatusBadRequest)
		return
	}

	// Рандомный статус
	status := statuses[rand.Intn(len(statuses))]

	var accrual *float64
	if status == "PROCESSED" {
		val := rand.Float64() * 100 // случайное значение начисления
		accrual = &val
	}

	resp := Response{
		Order:   orderNumber,
		Status:  status,
		Accrual: accrual,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/api/orders/", handler)

	addr := "localhost:9999"
	fmt.Printf("Mock server running at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
