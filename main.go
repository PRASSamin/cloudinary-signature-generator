package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"github.com/cloudinary/cloudinary-go/v2/api"
)

func signatureHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Parse the JSON body
	var requestData map[string]string
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Extract the required fields from the JSON body
	apiSecretKey := requestData["api_secret"]
	folder := requestData["folder"]
	publicId := requestData["public_id"]
	timestamp := requestData["timestamp"]

	// Validate the required fields
	if apiSecretKey == "" || folder == "" || publicId == "" || timestamp == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Generate the signature using Cloudinary's SDK
	paramsToSign := url.Values{
		"folder":    []string{folder},
		"public_id": []string{publicId},
		"timestamp": []string{timestamp},
	}
	signature, err := api.SignParameters(paramsToSign, apiSecretKey)
	if err != nil {
		http.Error(w, "Error generating signature", http.StatusInternalServerError)
		return
	}

	// Respond with the generated signature
	response := map[string]string{
		"signature": signature,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	initialJson := map[string]interface{}{
		"name":        "Cloudinary Signature Generator",
		"version":     "1.0.0",
		"author":      "PRAS",
		"author_url":  "https://pras.me",
		"repository":  "https://github.com/PRASSamin/cloudinary-signature-generator",
		"license":     "MIT",
		"description": "A high-performance API for generating Cloudinary signatures. Ideal for projects requiring signature generation without backend support or where serverless platforms (e.g., Cloudflare Workers) are not compatible with Cloudinary's native signature API.",
		"api": map[string]interface{}{
			"endpoint": "/api/gen/signature",
			"method":   "POST",
			"body": map[string]string{
				"api_secret": "Your Cloudinary API Secret Key",
				"folder":         "Target folder for the asset",
				"public_id":      "Unique identifier for the asset",
				"timestamp":      "Unix timestamp of the request",
			},
			"response": map[string]string{
				"signature": "Generated signature for authentication",
			},
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(initialJson)
	})
	
	http.HandleFunc("/api/gen/signature", signatureHandler)

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
