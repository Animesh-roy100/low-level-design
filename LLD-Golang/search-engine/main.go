package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

// Entity: Document
type Document struct {
	ID        string
	Title     string
	Content   string
	Timestamp time.Time
}

// Entity: Posting
type Posting struct {
	DocID     string
	Frequency int
	Positions []int
}

// InvertedIndex (core data structure)
type InvertedIndex map[string][]Posting

// Function to build index (Indexing Strategy)
func BuildIndex(docs []Document) InvertedIndex {
	index := make(InvertedIndex)
	for _, doc := range docs {
		tokens := tokenize(doc.Content) // Simple split, add stemming in prod
		termFreq := make(map[string]int)
		for _, token := range tokens {
			termFreq[token]++
			// Add position if needed
		}
		for term, freq := range termFreq {
			index[term] = append(index[term], Posting{DocID: doc.ID, Frequency: freq})
		}
	}
	return index
}

func tokenize(text string) []string {
	// Basic: lower, split, remove stops. Use libraries like "golang.org/x/text" for advanced.
	stops := map[string]bool{"the": true, "a": true} // Example stops
	words := strings.Fields(strings.ToLower(text))
	var tokens []string
	for _, w := range words {
		if !stops[w] {
			tokens = append(tokens, w)
		}
	}
	return tokens
}

// Ranking: Simple TF-IDF
func RankResults(queryTerms []string, index InvertedIndex, totalDocs int) []string {
	scores := make(map[string]float64)
	for _, term := range queryTerms {
		postings, ok := index[term]
		if !ok {
			continue
		}
		idf := math.Log(float64(totalDocs) / float64(len(postings)))
		for _, p := range postings {
			tf := float64(p.Frequency)
			scores[p.DocID] += tf * idf
		}
	}
	// Sort by score descending
	var ranked []string
	for docID := range scores {
		ranked = append(ranked, docID)
	}
	sort.Slice(ranked, func(i, j int) bool {
		return scores[ranked[i]] > scores[ranked[j]]
	})
	return ranked
}

// Autocomplete: Simple Trie
type TrieNode struct {
	Children map[rune]*TrieNode
	IsEnd    bool
	Freq     int // For popularity-based suggestions
}

type Trie struct {
	Root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{Root: &TrieNode{Children: make(map[rune]*TrieNode)}}
}

func (t *Trie) Insert(word string, freq int) {
	node := t.Root
	for _, ch := range word {
		if _, ok := node.Children[ch]; !ok {
			node.Children[ch] = &TrieNode{Children: make(map[rune]*TrieNode)}
		}
		node = node.Children[ch]
	}
	node.IsEnd = true
	node.Freq = freq
}

func (t *Trie) Suggest(prefix string, limit int) []string {
	// Traverse to prefix node, then DFS for suggestions, sort by freq
	// Implementation omitted for brevity; return top suggestions.
	return []string{} // Placeholder
}

// Example Usage
func main() {
	docs := []Document{{ID: "1", Content: "Hello world search engine"}}
	index := BuildIndex(docs)
	results := RankResults([]string{"search"}, index, len(docs))
	// Output results

	fmt.Println(results)
}
