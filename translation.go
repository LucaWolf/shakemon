package main

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

type TranslationAPIResponse struct {
	Contents struct {
		Translated string `json:"translated"`
	} `json:"contents"`
}

func getTranslation(text string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetFormData(map[string]string{
			"text": text,
		}).
		SetError(apiError{}).
		SetResult(TranslationAPIResponse{}).
		Post("https://api.funtranslations.com/translate/shakespeare.json")

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		err, _ := resp.Error().(*apiError)
		return "", &err.Error
	}

	// Extract the description from the API response
	data, ok := resp.Result().(*TranslationAPIResponse)
	if !ok {
		return "", errors.New("response was not a Translation API format")
	}

	if len(data.Contents.Translated) == 0 {
		return "", errors.New("no translation found")
	}

	return data.Contents.Translated, nil
}
