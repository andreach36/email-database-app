package main

import (
	datasetindex "email-database-api/dataset-index"
	"email-database-api/emails"
	"log"
	"net/http"
	"os"
	"runtime/pprof"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	// CPU Profiling
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer cpuFile.Close()
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()


	// Cargar las variables de entorno
	loadEnvVars()

	// Obtener la cantidad de datos total en ZincSearch
	hits, err := emails.GetTotalHits()
	if err != nil {
		log.Printf("error obteniendo los hits: %v\n", err)
		return
	}

	if hits == 0 {
		log.Println("No se encontraron datos. Iniciando subida de datos ...")
		datasetindex.IndexAndCreateJson()
	} else {
		log.Println("Los datos ya se encuentran cargados en la base de datos.")
		log.Printf("Total datos cargados: %d", hits)
	}

	// Configurar servidor
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Mount("/emails", emails.EmailsRoutes())
	log.Println("Starting server on :3000...")
	http.ListenAndServe(os.Getenv("PORT"), r)


	// Memory Profiling
	memFile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer memFile.Close()
	pprof.WriteHeapProfile(memFile)
}





// Obtener el mapping
// mapping, err := emails.GetMapping()
// if err != nil {
// 	fmt.Printf("Error obteniendo el mapping: %v\n", err)
// 	return
// }

// Mostrar el mapping
// fmt.Println("Mapping del índice:")
// fmt.Println(string(mapping))
