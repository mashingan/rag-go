package main

import (
	"fmt"
	"log"

	"github.com/kirill-scherba/word2vec"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	model, err := word2vec.Load("vectors.bin")
	if err != nil {
		log.Fatalf("failed to load model: %v", err)
	}
	queryWord := "example"
	nearestWords := make([]word2vec.Nearest, 10)
	if err := model.Lookup(queryWord, nearestWords); err != nil {
		log.Fatalf("lookup failed: %v", err)
	}
	log.Printf("Words closest to '%s': \n", queryWord)
	for _, word := range nearestWords {
		fmt.Printf("%20s\t%f\n", word.Word, word.Distance)
	}

	demonstrateEmbeddingNuances(model)
}

func demonstrateEmbeddingNuances(model *word2vec.Model) {
	queries := []string{
		"dog",
		"what is a dog",
		"playful canine companion",
	}

	for _, query := range queries {
		vec := make([]float32, model.Size())
		model.Embedding(query, vec)

		fmt.Printf("Query: %s\n", query)
		fmt.Printf("Embedding Vector Length: %d\n", len(vec))

		// Find nearest neighbors
		nearest := make([]word2vec.Nearest, 5)
		model.NearestToVector(vec, nearest, nil)

		fmt.Println("Semantic Neighbors:")
		for _, n := range nearest {
			fmt.Printf("  - %s (Distance: %f)\n", n.Word, n.Distance)
		}
		fmt.Println()
	}
}
