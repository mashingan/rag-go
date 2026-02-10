package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/kirill-scherba/word2vec"
)

// AdvancedSemanticContext represents a rich semantic understanding
type AdvancedSemanticContext struct {
	// Raw vector representation in high-dimensional space
	EmbeddingVector []float32

	// Similarity scores between this vector and other concepts
	SimilarityScores map[string]float32

	// Semantic relationships and their strengths
	RelationshipGraph map[string]RelationshipDetails
}

// RelationshipDetails provides nuanced semantic connections
type RelationshipDetails struct {
	Words       []string // Related words
	Strength    float32  // Semantic proximity (0-1)
	Connections []string // Types of relationships
}

// SemanticAnalyzer demonstrates advanced semantic context usage
func SemanticAnalyzer(model *word2vec.Model, query string) AdvancedSemanticContext {
	// 1. Generate embedding vector
	embedding := make([]float32, model.Size())
	if err := model.Embedding(query, embedding); err != nil {
		log.Println(err)
		return AdvancedSemanticContext{}
	}

	// 2. Find nearest neighbors
	neighbors := make([]word2vec.Nearest, 20)
	if err := model.NearestToVector(embedding, neighbors, nil); err != nil {
		log.Println(err)
		return AdvancedSemanticContext{}
	}

	// 3. Build semantic similarity and relationship graph
	similarityScores := make(map[string]float32)
	relationshipGraph := make(map[string]RelationshipDetails)

	for _, neighbor := range neighbors {
		// Compute cosine similarity
		similarity := neighbor.Distance

		// Categorize relationship strength
		strengthCategory := categorizeStrength(similarity)

		similarityScores[neighbor.Word] = similarity

		relationshipGraph[neighbor.Word] = RelationshipDetails{
			Words:       []string{neighbor.Word},
			Strength:    similarity,
			Connections: strengthCategory,
		}
	}

	return AdvancedSemanticContext{
		EmbeddingVector:   embedding,
		SimilarityScores:  similarityScores,
		RelationshipGraph: relationshipGraph,
	}
}

// Categorize relationship strength
func categorizeStrength(similarity float32) []string {
	switch {
	case similarity > 0.8:
		return []string{"very_strong", "core_concept", "direct_relation"}
	case similarity > 0.6:
		return []string{"strong", "close_semantic_link"}
	case similarity > 0.4:
		return []string{"moderate", "contextual_relation"}
	case similarity > 0.2:
		return []string{"weak", "distant_connection"}
	default:
		return []string{"minimal", "peripheral_concept"}
	}
}

// Visualize Visualization and Analysis
func (ctx AdvancedSemanticContext) Visualize() {
	fmt.Println("🔍 Semantic Context Analysis")
	fmt.Println("-----------------------------")

	// Sort relationships by strength
	sortedRelationships := make([]struct {
		Word        string
		Strength    float32
		Connections []string
	}, 0, len(ctx.RelationshipGraph))

	for word, details := range ctx.RelationshipGraph {
		sortedRelationships = append(sortedRelationships, struct {
			Word        string
			Strength    float32
			Connections []string
		}{
			Word:        word,
			Strength:    details.Strength,
			Connections: details.Connections,
		})
	}

	// Sort by strength in descending order
	sort.Slice(sortedRelationships, func(i, j int) bool {
		return sortedRelationships[i].Strength > sortedRelationships[j].Strength
	})

	// Print top relationships
	fmt.Println("Top Semantic Relationships:")
	for i, rel := range sortedRelationships {
		if i >= 5 { // Limit to top 5
			break
		}
		fmt.Printf("  %s:\n", rel.Word)
		fmt.Printf("    Strength:     %.2f\n", rel.Strength)
		fmt.Printf("    Connections:  %v\n", rel.Connections)
	}
}

// GenerateInsight Semantic Reasoning Example
func (ctx AdvancedSemanticContext) GenerateInsight(query string) string {
	// Find top relationships
	var topRelationships []string
	for word, details := range ctx.RelationshipGraph {
		if details.Strength > 0.6 {
			topRelationships = append(topRelationships, word)
		}
	}

	// Generate narrative based on top relationships
	if len(topRelationships) > 0 {
		return fmt.Sprintf(
			"The concept of '%s' is deeply connected to %s, revealing a rich semantic landscape of interconnected ideas.",
			query,
			topRelationships[0],
		)
	}

	return "No significant semantic relationships found."
}

func main() {
	// Load pre-trained word2vec model
	model, err := word2vec.Load("./vectors.bin")
	if err != nil {
		panic(err)
	}

	// Example queries
	queries := []string{
		"dog",
		"what is dog doing?",
		"programming",
		//"machine learning",
		"love and peace",
	}

	for _, query := range queries {
		fmt.Printf("\n🔬 Analyzing Query: %s\n", query)

		// Generate semantic context
		semanticContext := SemanticAnalyzer(model, query)

		// Visualize relationships
		semanticContext.Visualize()

		// Generate semantic insight
		insight := semanticContext.GenerateInsight(query)
		fmt.Println("\nSemantic Insight:")
		fmt.Println(insight)
	}
}
