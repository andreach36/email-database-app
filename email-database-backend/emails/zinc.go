package emails

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func makeRequestZinc(method, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(os.Getenv("ZINC_USER"), os.Getenv("ZINC_PASSWORD"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	// fmt.Println(resp)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo la respuesta: %w", err)
	}

	// Manejar diferentes códigos de estado
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"solicitud fallida con código %d, respuesta: %s",
			resp.StatusCode, string(responseBody),
		)
	}

	return responseBody, nil
}

type TotalHitsResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
	} `json:"hits"`
}

func GetTotalHits() (int, error) {
	// Configuración del host y el índice
	index := os.Getenv("ZINC_INDEX")
	zincHost := os.Getenv("ZINC_HOST")
	zincURL := zincHost + "/api/" + index + "/_search"

	query := `{"query":{"match_all":{}}}`

	req, err := http.NewRequest(http.MethodPost, zincURL, strings.NewReader(query))
	if err != nil {
		return 0, fmt.Errorf("error creando la solicitud: %v", err)
	}
	// Autenticación básica
	req.SetBasicAuth(os.Getenv("ZINC_USER"), os.Getenv("ZINC_PASSWORD"))

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error realizando la solicitud: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error leyendo la respuesta de ZincSearch: %v", err)
	}

	var totalHitsResponse TotalHitsResponse
	if err := json.Unmarshal(body, &totalHitsResponse); err != nil {
		fmt.Println(string(body))
		return 0, fmt.Errorf("error parseando la respuesta de ZincSearch: %v", err)
	}

	return totalHitsResponse.Hits.Total.Value, nil
}

func GetMapping() ([]byte, error) {
	// Configuración del host y el índice
	index := os.Getenv("ZINC_INDEX")
	zincHost := os.Getenv("ZINC_HOST")
	zincURL := zincHost + "/api/" + index + "/_mapping"

	// Crear una nueva solicitud GET
	req, err := http.NewRequest(http.MethodGet, zincURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando la solicitud: %v", err)
	}

	// Autenticación básica
	req.SetBasicAuth(os.Getenv("ZINC_USER"), os.Getenv("ZINC_PASSWORD"))

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error realizando la solicitud: %v", err)
	}
	defer resp.Body.Close()

	// Validar el código de respuesta
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en la respuesta: %s, body: %s", resp.Status, string(body))
	}

	// Leer el cuerpo de la respuesta
	return io.ReadAll(resp.Body)
}
