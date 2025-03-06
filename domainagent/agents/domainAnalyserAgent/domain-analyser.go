package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type DomainAnalysis struct {
	Domain     string   `json:"domain"`
	Name       string   `json:"name"`
	Extension  string   `json:"extension"`
	Words      []string `json:"words"`
	Categories []string `json:"categories"`
}

var extensionMap = map[string][]string{
	"com":      {"Commercial", "General"},
	"co.uk":    {"UK Business", "United Kingdom"},
	"uk":       {"United Kingdom", "General"},
	"org":      {"Organization", "Non-profit"},
	"net":      {"Network", "Technology"},
	"io":       {"Technology", "Startup"},
	"ai":       {"Artificial Intelligence", "Technology"},
	"app":      {"Application", "Software"},
	"dev":      {"Development", "Technology"},
	"tech":     {"Technology", "Innovation"},
	"store":    {"E-commerce", "Retail"},
	"shop":     {"E-commerce", "Retail"},
	"blog":     {"Content", "Personal"},
	"edu":      {"Education", "Learning"},
	"gov":      {"Government", "Official"},
	"info":     {"Information", "General"},
	"biz":      {"Business", "Commercial"},
	"travel":   {"Travel", "Tourism"},
	"agency":   {"Service", "Business"},
	"design":   {"Creative", "Art"},
	"health":   {"Healthcare", "Wellness"},
	"media":    {"Media", "Entertainment"},
	"finance":  {"Finance", "Banking"},
	"property": {"Real Estate", "Property"},
}

// todo create categorisation model
var wordToCategories = map[string][]string{
	"property":   {"Real Estate", "Property Investment", "Housing"},
	"mortgage":   {"Finance", "Banking", "Home Loans", "Property"},
	"calculator": {"Tools", "Utility", "Finance"},
	"meat":       {"Food", "Cuisine", "Butcher", "Restaurant"},
	"rice":       {"Food", "Cuisine", "Grocery"},
	"flat":       {"Real Estate", "Property", "Apartment"},
	"prices":     {"E-commerce", "Price Comparison", "Shopping"},
	"redesign":   {"Design", "Renovation", "Improvement"},
	"mandrill":   {"Animal", "Wildlife", "Email Service"}, // Mandrill is also an email service
}

func AnalyseDomain(domainName string) (*DomainAnalysis, error) {
	// clean
	domainName = strings.ToLower(strings.TrimSpace(domainName))

	// valid format
	if !isValidDomain(domainName) {
		return nil, fmt.Errorf("invalid domain name format: %s", domainName)
	}

	// split domain and extension/s
	parts := strings.Split(domainName, ".")
	var name, extension string

	if len(parts) == 2 {
		// simple like example.com
		name = parts[0]
		extension = parts[1]
	} else if len(parts) == 3 {
		// e.g. subdomain or country code co.uk
		name = parts[0]
		extension = strings.Join(parts[1:], ".")
	} else {
		return nil, fmt.Errorf("invalid domain - needs extension and not to be a subdomain: %s", domainName)
	}

	// valid name
	words := extractWords(name)

	categories := determineCategories(words, extension)

	analysis := &DomainAnalysis{
		Domain:     domainName,
		Name:       name,
		Extension:  extension,
		Words:      words,
		Categories: categories,
	}

	return analysis, nil
}

func isValidDomain(domain string) bool {
	pattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, domain)
	return match
}

func extractWords(name string) []string {
	// change hyphens to space
	name = strings.ReplaceAll(name, "-", " ")

	// Common prefixes to remove
	prefixes := []string{"www", "my", "the", "get", "best"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(name, prefix) {
			name = strings.TrimPrefix(name, prefix)
			break
		}
	}

	var words []string

	// First check if there are spaces already
	if strings.Contains(name, " ") {
		words = strings.Fields(name)
	} else {
		// Use camelCase splitting
		re := regexp.MustCompile(`[A-Z][a-z]+|[a-z]+`)
		camelCaseWords := re.FindAllString(name, -1)

		if len(camelCaseWords) > 1 {
			words = camelCaseWords
		} else {
			// todo:
			// Try dictionary-based word extraction (simplified here)
			// In a real implementation, this would use a proper dictionary
			// or NLP tokenization method
			words = []string{name}
		}
	}

	return words

}

// todo: replace with big llm call to get good info
func determineCategories(words []string, extension string) []string {
	categoryMap := make(map[string]bool)

	// add categories based on extension
	if cats, ok := extensionMap[extension]; ok {
		for _, cat := range cats {
			categoryMap[cat] = true
		}
	}

	// based on domain words
	for _, word := range words {
		if cats, ok := wordToCategories[word]; ok {
			for _, cat := range cats {
				categoryMap[cat] = true
			}
		}
	}

	// back to slice
	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}

	return categories
}

func analyseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
		return
	}

	// parse domain from request
	var requestData struct {
		Domain string `json:"domain"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.Domain == "" {
		http.Error(w, "Domain name is required", http.StatusBadRequest)
		return
	}

	// analyse domain
	analysis, err := AnalyseDomain(requestData.Domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// return results as json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)

}

func main() {
	// get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// register handler
	http.HandleFunc("/", analyseHandler)

	// start server
	fmt.Printf("Domain Analyser Agent starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
