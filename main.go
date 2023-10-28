package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/translate/{kind}/{name}", translatePokemon).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func translatePokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kind := vars["kind"]
	name := vars["name"]

	if kind != "pokemon" {
		w.WriteHeader(http.StatusBadRequest)
		resp := apiError{
			Error: apiErrorDetails{
				Message: "invalid kind: only 'pokemon' is supported",
				Code:    "bad request",
			},
		}
		e := json.NewEncoder(w)
		e.Encode(resp)
		return
	}

	if description, err := getDescription(name); err != nil {
		resp := apiError{
			Error: apiErrorDetails{
				Message: err.Error(),
				Code:    StringOrInt("failed to get description for " + name),
			},
		}
		w.WriteHeader(http.StatusFailedDependency)
		e := json.NewEncoder(w)
		e.Encode(resp)
	} else if translation, err := getTranslation(description); err != nil {
		resp := apiError{
			Error: apiErrorDetails{
				Message: err.Error(),
				Code:    StringOrInt("failed to get translation for " + name),
			},
		}
		w.WriteHeader(http.StatusFailedDependency)
		e := json.NewEncoder(w)
		e.Encode(resp)
	} else {
		resp := apiReply{
			Name: name,
			Desc: translation,
		}
		w.WriteHeader(http.StatusOK)
		e := json.NewEncoder(w)
		e.Encode(resp)
	}
}
