package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Language represents the language structure
type Language struct {
	Typename         string `json:"__typename"`
	Code             string `json:"code"`
	CodeAirportTaxis string `json:"codeAirportTaxis"`
	CountryFlag      string `json:"countryFlag"`
	Name             string `json:"name"`
}

// LanguagesResponse represents the API response structure
type LanguagesResponse struct {
	Data []Language `json:"data"`
}

// UniqueCountriesResponse represents unique countries
type UniqueCountriesResponse struct {
	Countries []string `json:"countries"`
}

// LanguageModel handles language-related operations
type LanguageModel struct{}

// FetchLanguagesFromAPI retrieves languages from Booking.com API
func (m *LanguageModel) FetchLanguagesFromAPI() ([]Language, error) {
    url := "https://booking-com18.p.rapidapi.com/languages"
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", os.Getenv("RAPID_API_KEY"))

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code: %d, body: %s", 
            resp.StatusCode, string(body))
    }

    var languagesResp LanguagesResponse
    if err = json.Unmarshal(body, &languagesResp); err != nil {
        return nil, fmt.Errorf("error parsing JSON: %w", err)
    }

    return languagesResp.Data, nil
}
// ExtractUniqueCountries extracts unique countries from languages
func (m *LanguageModel) ExtractUniqueCountries(languages []Language) []string {
	uniqueCountries := make(map[string]bool)
	
	for _, lang := range languages {
		country := m.extractCountryFromLanguage(lang)
		if country != "" {
			uniqueCountries[country] = true
		}
	}

	countries := make([]string, 0, len(uniqueCountries))
	for country := range uniqueCountries {
		countries = append(countries, country)
	}

	return countries
}

// extractCountryFromLanguage is a private method to extract country
func (m *LanguageModel) extractCountryFromLanguage(lang Language) string {
	// Extract country using multiple methods
	re := regexp.MustCompile(`\(([^)]+)\)$`)
	matches := re.FindStringSubmatch(lang.Name)
	if len(matches) > 1 {
		return strings.ToUpper(matches[1])
	}

	if lang.CountryFlag != "" {
		return strings.ToUpper(lang.CountryFlag)
	}

	parts := strings.Split(lang.Code, "-")
	if len(parts) > 1 {
		return strings.ToUpper(parts[1])
	}

	return ""
}

// ExportCountriesToJSON saves unique countries to a JSON file
func (m *LanguageModel) ExportCountriesToJSON(countries []string) error {
	countriesResp := UniqueCountriesResponse{Countries: countries}
	
	file, err := os.Create("unique_countries.json")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(countriesResp); err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	return nil
}