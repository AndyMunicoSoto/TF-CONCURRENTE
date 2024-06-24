package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

type House struct {
	Size     float64
	Bedrooms float64
	Age      float64
	Location string
}

type PredictionRequest struct {
	House House
}

type PredictionResponse struct {
	Price float64
}

func calculatePrice(h House) float64 {
	// Calcular el precio base sin aplicar la variabilidad
	basePrice := 1200*h.Size + 500*h.Bedrooms
	ageMultiplier := 1 + 0.08*h.Age
	locationMultiplier := 1.0

	switch h.Location {
	case "A":
		locationMultiplier = 1.06
	case "B":
		locationMultiplier = 1.02
	case "D":
		locationMultiplier = 0.98
	}

	// Introducir variabilidad en el número de habitaciones y antigüedad
	rand.Seed(time.Now().UnixNano())

	// Variabilidad para el número de habitaciones (±1 o ±2 habitaciones)
	bedroomsVariability := float64(rand.Intn(3) - 1) // Entre -1 y 1 habitaciones
	bedroomsAdjusted := h.Bedrooms + bedroomsVariability

	// Variabilidad para la antigüedad (±1 o ±2 años)
	ageVariability := float64(rand.Intn(3) - 1) // Entre -1 y 1 años
	ageAdjusted := h.Age + ageVariability

	// Calcular el precio final aplicando la variabilidad
	finalPrice := basePrice * ageMultiplier * locationMultiplier

	// Utilizar las variables bedroomsAdjusted y ageAdjusted en el cálculo final
	finalPriceAdjusted := finalPrice + bedroomsAdjusted + ageAdjusted

	return finalPriceAdjusted
}

func calculateMAE(data [][]string) float64 {
	var totalError float64
	var count int

	for _, row := range data {
		size, _ := strconv.ParseFloat(row[0], 64)
		bedrooms, _ := strconv.ParseFloat(row[1], 64)
		age, _ := strconv.ParseFloat(row[2], 64)
		location := row[3]
		actualPrice, _ := strconv.ParseFloat(row[4], 64)

		predictedPrice := calculatePrice(House{Size: size, Bedrooms: bedrooms, Age: age, Location: location})
		error := math.Abs(predictedPrice - actualPrice)
		totalError += error
		count++
	}

	return totalError / float64(count)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := json.NewDecoder(conn)
	var request PredictionRequest
	if err := decoder.Decode(&request); err != nil {
		log.Println("Error decodificando la solicitud:", err)
		return
	}

	price := calculatePrice(request.House)
	fmt.Printf("Predicción para casa con características %v: %.2f dólares\n", request.House, price)

	response := PredictionResponse{Price: price}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&response); err != nil {
		log.Println("Error codificando la respuesta:", err)
		return
	}
}

func predictPriceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var request PredictionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	price := calculatePrice(request.House)
	fmt.Printf("Predicción para casa con características %v: %.2f dólares\n", request.House, price)

	response := PredictionResponse{Price: price}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func predictPriceHandler2(w http.ResponseWriter, r *http.Request) {
	var request PredictionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	price := calculatePrice(request.House)
	fmt.Printf("Predicción para casa con características %v: %.2f dólares\n", request.House, price)

	response := PredictionResponse{Price: price}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Descarga el archivo CSV durante la ejecución del contenedor
	resp, err := http.Get("https://raw.githubusercontent.com/AndyMunicoSoto/TA-4-CONCURRENTE/main/house_price_8.csv")
	if err != nil {
		log.Fatal("Error al descargar el archivo CSV:", err)
	}
	defer resp.Body.Close()

	// Lee el CSV descargado
	reader := csv.NewReader(resp.Body)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error leyendo el archivo CSV:", err)
	}

	// Calcular la precisión
	mae := calculateMAE(data[1:]) // Excluir encabezado
	fmt.Printf("Error absoluto medio (MAE) del modelo: %.2f dólares\n", mae)

	// Configurar el servidor HTTP
	http.HandleFunc("/predict", predictPriceHandler)

	// Servir el archivo index.html desde la raíz
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Println("El servidor está escuchando en el puerto 8000")
	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		log.Fatal(err)
	}
}
