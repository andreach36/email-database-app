package main

import (
	datasetindex "email-database-api/dataset-index"
	"log"
	"os"
	"runtime/pprof"

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

	loadEnvVars()
	datasetindex.IndexAndCreateJson()

	// Memory Profiling
	memFile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer memFile.Close()
	pprof.WriteHeapProfile(memFile)
}
