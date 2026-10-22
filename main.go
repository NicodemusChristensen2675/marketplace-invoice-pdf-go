package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/invoice", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "use POST", http.StatusMethodNotAllowed)
			return
		}
		var order Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, err := GenerateInvoice(context.Background(), client, order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
