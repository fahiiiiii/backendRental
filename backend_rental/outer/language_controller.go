
// controllers/language_controller.go
package controllers

import (
	"backend_rental/models"
	"log"
	"fmt"
	"time"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/logs"
)

// LanguageController handles language-related web requests
type LanguageController struct {
	beego.Controller
}

// Get handles GET request for languages
// This will map to GET /v1/languages
func (c *LanguageController) Get() {
	languageModel := &models.LanguageModel{}

	// Fetch languages from API
	languages, err := languageModel.FetchLanguagesFromAPI()
	if err != nil {
		logs.Error("Failed to fetch languages: %v", err)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Respond with languages
	c.Data["json"] = languages
	c.ServeJSON()
}

// CountryController handles country-related web requests
type CountryController struct {
	beego.Controller
}

// Get handles GET request for countries
// This will map to GET /v1/countries
func (c *CountryController) Get() {
	// Add detailed logging
	log.Println("Countries API: Fetching countries started")
	
	start := time.Now()
	countries, err := FetchCountriesAPI()
	duration := time.Since(start)

	log.Printf("Countries API: Fetch completed in %v", duration)

	if err != nil {
		log.Printf("Countries API Error: %v", err)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	log.Printf("Countries API: Fetched %d countries", len(countries))

	c.Data["json"] = countries
	c.ServeJSON()
}


// GetUniqueCountries retrieves unique countries
func (c *LanguageController) GetUniqueCountries() {
	languageModel := &models.LanguageModel{}

	// Fetch languages from API
	languages, err := languageModel.FetchLanguagesFromAPI()
	if err != nil {
		logs.Error("Failed to fetch languages: %v", err)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Extract unique countries
	countries := languageModel.ExtractUniqueCountries(languages)

	// Respond with countries
	c.Data["json"] = countries
	c.ServeJSON()
}

// FetchCountriesAPI retrieves countries from an API or data source
func FetchCountriesAPI() ([]string, error) {
    languageModel := &models.LanguageModel{}

    // Fetch languages from API
    languages, err := languageModel.FetchLanguagesFromAPI()
    if err != nil {
        return nil, fmt.Errorf("failed to fetch languages: %v", err)
    }

    // Extract unique countries
    countries := languageModel.ExtractUniqueCountries(languages)

    // Optional: Export countries to JSON file
    if err := languageModel.ExportCountriesToJSON(countries); err != nil {
        log.Printf("Warning: Failed to export countries to JSON: %v", err)
    }

    return countries, nil
}








