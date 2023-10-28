package main

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

type PokemonAPIResponse struct {
	Name        string `json:"name"`
	FlavorTexts []struct {
		FlavorText string `json:"flavor_text"`
		Language   struct {
			Name string `json:"name"`
		} `json:"language"`
	} `json:"flavor_text_entries"`
}

func getDescription(name string) (string, error) {
	client := resty.New().SetLogger(NewLogger())

	resp, err := client.R().
		SetPathParams(map[string]string{
			"name": name,
		}).
		SetResult(PokemonAPIResponse{}).
		Get("https://pokeapi.co/api/v2/pokemon-species/{name}")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", errors.New("failed to get description")
	}

	// Extract the description from the API response
	data, ok := resp.Result().(*PokemonAPIResponse)
	if !ok {
		return "", errors.New("response was not a Pokemon API format")
	}

	for _, entry := range data.FlavorTexts {
		if entry.Language.Name == "en" {
			return entry.FlavorText, nil
		}
	}

	return "", errors.New("no description found")
}
