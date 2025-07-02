package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClientePotencialTest representa la estructura para pruebas
type ClientePotencialTest struct {
	ID           uint      `json:"id,omitempty"`
	Nombre       string    `json:"nombre"`
	Email        string    `json:"email"`
	Telefono     string    `json:"telefono"`
	Fuente       string    `json:"fuente"`
	FechaCaptura time.Time `json:"fecha_captura,omitempty"`
	Interes      string    `json:"interes"`
	Notas        string    `json:"notas"`
	Estado       string    `json:"estado,omitempty"`
	SucursalID   uint      `json:"sucursal_id"`
}

func main() {
	baseURL := "http://localhost:8080"

	// Datos de prueba
	cliente := ClientePotencialTest{
		Nombre:     "Juan Pérez",
		Email:      "juan.perez@email.com",
		Telefono:   "+34 123 456 789",
		Fuente:     "formulario_web",
		Interes:    "productos_financieros",
		Notas:      "Cliente interesado en hipotecas",
		SucursalID: 1,
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(cliente)
	if err != nil {
		fmt.Printf("Error al convertir a JSON: %v\n", err)
		return
	}

	fmt.Printf("Enviando datos: %s\n", string(jsonData))

	// Crear cliente HTTP
	client := &http.Client{Timeout: 10 * time.Second}

	// Realizar petición POST
	resp, err := client.Post(
		baseURL+"/api/clientes",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Printf("Error en la petición: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Leer respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error al leer respuesta: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Respuesta: %s\n", string(body))

	// Probar obtener todos los clientes
	fmt.Println("\n--- Probando GET /api/clientes ---")
	resp2, err := client.Get(baseURL + "/api/clientes")
	if err != nil {
		fmt.Printf("Error en GET: %v\n", err)
		return
	}
	defer resp2.Body.Close()

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		fmt.Printf("Error al leer respuesta GET: %v\n", err)
		return
	}

	fmt.Printf("Status GET: %d\n", resp2.StatusCode)
	fmt.Printf("Respuesta GET: %s\n", string(body2))
}
