package emails

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type EmailHandler struct {
	Message_ID string `json:"Message-ID"`
	Date       string `json:"Date"`
	From       string `json:"from"`
	To         string `json:"to"`
	Subject    string `json:"subject"`
	Body       string `json:"Body"`
}

// Función para listar todos los emails
func (e EmailHandler) GetAllEmails(w http.ResponseWriter, r *http.Request) {

	// Parsear los parámetros 'from' y 'size' con valores predeterminados
	from := r.URL.Query().Get("from")
	size := r.URL.Query().Get("size")

	if from == "" {
		from = "0"
	}
	if size == "" {
		size = "10"
	}

	// Construir la consulta
	index := os.Getenv("ZINC_INDEX")
	zincHost := os.Getenv("ZINC_HOST")
	zincUrl := zincHost + "/api/" + index + "/_search"
	query := `{
		"from": ` + from + `,
		"size": ` + size + `,
		"query": {
			"match_all": {}
		}
	}`

	// Hacer la solicitud a ZincSearch
	response, err := makeRequestZinc(http.MethodPost, zincUrl, []byte(query))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enviar la respuesta al cliente
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// Función para manejar la búsqueda por palabra clave
func (e EmailHandler) SearchEmails(w http.ResponseWriter, r *http.Request) {
	// Obtener los parámetros de búsqueda
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		http.Error(w, "Missing 'q' query parameter", http.StatusBadRequest)
		return
	}

	// Valores predeterminados para paginación
	from, size := 0, 10
	fromParam := r.URL.Query().Get("from")
	sizeParam := r.URL.Query().Get("size")

	// Convertir parámetros de paginación
	if fromParam != "" {
		var err error
		from, err = strconv.Atoi(fromParam)
		if err != nil {
			http.Error(w, "Invalid 'from' query parameter", http.StatusBadRequest)
			return
		}
	}
	if sizeParam != "" {
		var err error
		size, err = strconv.Atoi(sizeParam)
		if err != nil {
			http.Error(w, "Invalid 'size' query parameter", http.StatusBadRequest)
			return
		}
	}

	// Construir el JSON de la consulta para ZincSearch
	query := map[string]interface{}{
		"search_type": "match",
		"from":        from,
		"size":        size,
		"query": map[string]interface{}{
			"term":  keyword,
			"field": "_all",
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		http.Error(w, "Failed to serialize query to JSON", http.StatusInternalServerError)
		return
	}

	// Depurar el JSON de la consulta
	fmt.Printf("Query JSON generado: %s\n", queryJSON)

	// URL y credenciales para ZincSearch
	index := os.Getenv("ZINC_INDEX")
	zincHost := os.Getenv("ZINC_HOST")
	zincUrl := fmt.Sprintf("%s/api/%s/_search", zincHost, index)

	// Realizar la solicitud a ZincSearch
	response, err := makeRequestZinc(http.MethodPost, zincUrl, queryJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Enviar la respuesta al cliente
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
