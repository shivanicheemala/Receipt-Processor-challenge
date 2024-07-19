package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ProcessResponse struct {
	ID string `json:"id"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

var receipts = make(map[string]Receipt)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/receipts/process", processReceiptHandler)
	mux.HandleFunc("/receipts/", getPointsHandler)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	receipts[id] = receipt

	log.Printf("Receipt stored with ID: %s", id)
	log.Printf("Current stored receipts: %+v", receipts) // Log the current state of stored receipts

	response := ProcessResponse{ID: id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")
	log.Printf("Fetching points for ID: %s", id)
	receipt, exists := receipts[id]
	if !exists {
		log.Printf("Receipt not found for ID: %s", id)
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	points := calculatePoints(receipt)

	response := PointsResponse{Points: points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func calculatePoints(receipt Receipt) int {
	points := 0

	for _, c := range receipt.Retailer {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			points++
		}
	}
	//log.Printf("points: ", points)

	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil {
		if total == float64(int(total)) {
			points += 50
		}
		if int(total*100)%25 == 0 {
			points += 25
		}
	}
	//log.Printf("points: ", points)

	points += (len(receipt.Items) / 2) * 5
	//log.Printf("points: ", points)

	for _, item := range receipt.Items {
		trimmedLength := len(item.ShortDescription)
		if trimmedLength%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				log.Printf("Original Price: %f", price)
				log.Printf("Price * 0.2: %f", price*0.2)
				log.Printf("Points per item: %f", math.Ceil(price*0.2))
				points += int(math.Ceil(price * 0.2))
			}
		}
	}
	//log.Printf("points: ", points)

	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil {
		if purchaseDate.Day()%2 != 0 {
			points += 6
		}
	}
	//log.Printf("points: ", points)

	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil {
		if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
			points += 10
		}
	}
	//log.Printf("points: ", points)

	return points
}
