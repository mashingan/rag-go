package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"

	"github.com/kirill-scherba/word2vec"
)

const topN = 10 // neighbor to consider

func generateAnswer(query string, context []word2vec.Nearest) string {
	if len(context) == 0 {
		return "couldn't find any related information"
	}
	sort.Slice(context, func(i, j int) bool {
		return context[i].Distance > context[j].Distance
	})
	words := make([]string, 0, topN)
	seen := map[string]struct{}{}
	for _, n := range context {
		w := strings.ToLower(n.Word)
		if w == strings.ToLower(query) {
			continue
		}
		if _, ok := seen[w]; ok {
			continue
		}
		seen[w] = struct{}{}
		words = append(words, w)
		if len(words) == topN {
			break
		}
	}
	log.Println("words:", words)
	if len(words) < 2 {
		return fmt.Sprintf("limited context of %d for query %s", len(words), query)
	}
	templates := []string{
		"%s is fundamentally connected to concepts like %s, suggesting a complex interplay of %s and deeper characteristics involving %s.",
		"Exploring %s reveals intricate relationships with %s, which illuminate its nature through %s and the underlying dynamics of %s.",
		"The essence of %s can be understood through its connections to %s, manifesting in qualities like %s and the broader context of %s.",
		"Diving into %s uncovers a rich tapestry of associations including %s, characterized by %s and the nuanced dimensions of %s.",
	}
	templateIdx := rand.Intn(len(templates))
	for len(words) <= 4 {
		words = append(words, words[0])
	}
	return fmt.Sprintf(templates[templateIdx], words[0], words[1], words[2], words[3])
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	model, err := word2vec.Load("./scratch/w2v-pure-go/vectors.bin")
	if err != nil {
		log.Fatalf("failed to load model: %v", err)
	}
	query := "what is a dog?"
	context, err := retrieveContext(query, model)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query: %s\n", query)
	/*
		fmt.Printf("Context: \n")
		for _, word := range context {
			fmt.Printf("%20s\t%f\n", word.Word, word.Distance)
		}
	*/
	answer := generateAnswer(query, context)
	fmt.Printf("Answer: %s\n", answer)

}

func retrieveContext(query string, model *word2vec.Model) ([]word2vec.Nearest, error) {
	log.Println("model.size:", model.Size())
	size := model.Size()
	// size := 10
	emb := make([]float32, size)
	if err := model.Embedding(query, emb); err != nil {
		return nil, err
	}
	context := make([]word2vec.Nearest, len(emb))
	if err := model.NearestToVector(emb, context, nil); err != nil {
		return nil, err
	}
	return context, nil
}
