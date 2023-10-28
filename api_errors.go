package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type StringOrInt string

type apiError struct {
	Error apiErrorDetails `json:"error"`
}

type apiErrorDetails struct {
	Message string      `json:"message"`
	Code    StringOrInt `json:"code"`
}

type apiReply struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func (s *StringOrInt) UnmarshalJSON(data []byte) error {
	var rawValue interface{}
	err := json.Unmarshal(data, &rawValue)
	if err != nil {
		return err
	}

	switch v := rawValue.(type) {
	case string:
		*s = StringOrInt(v)
	case float64:
		*s = StringOrInt(strconv.FormatFloat(v, 'f', -1, 64))
	case int:
		*s = StringOrInt(strconv.Itoa(v))
	default:
		return fmt.Errorf("unexpected value type: %T", v)
	}

	return nil
}

func (e *apiErrorDetails) Error() string {
	return fmt.Sprintf("API Error: Code=%s, Message=%s", e.Code, e.Message)
}
