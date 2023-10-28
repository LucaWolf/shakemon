# Background story

I already had some experience with the scope of this exercise, in particular of what APIs I wanted
to use and a basic structure of the application's flow.

The only question remaining was how would Codeium do when I threw at it all this information I wante to build on.

# Setup and usage
I have used the VS Code as text editor and the plugin provided for it. Upon installing the extension,
the firs use going to be prompting you to login (I have used my google account) and then present you with an 
`{...}` icon in the extensions panel. Expand it and start typing requests at the prompt.

> [!BRILLIANT]
> Not only Codeium generates a block of code for but also renders a decent description of what the code does
and supplements it by some recommendations. For keeping this post short, I will only copy the first
commentary section and skip listing for the rest of the output.


# The `main` frame
In order to follow an older version of this application, I wanted to create a simple web service that takes
a Pokemon name parameter and returns a Shakesperean translation of its description. So let's start issuing "demands":

> write me a gorrilamux server with an GET endpoint taking two path parameters
>>
>> Sure! Here's an example of a Gorilla Mux server with a GET endpoint that takes two path parameters:
```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/users/{name}/{age}", getUser).Methods("GET")

    log.Fatal(http.ListenAndServe(":8080", r))
}

func getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    age := vars["age"]

    fmt.Fprintf(w, "Name: %s\n", name)
    fmt.Fprintf(w, "Age: %s\n", age)
}
```
>> This server sets up a Gorilla Mux router and defines a GET endpoint at "/users/{name}/{age}". 
>> When a GET request is made to that endpoint, the getUser function is called.
>> The function extracts the name and age path parameters from the request using mux.Vars and writes them back as the response.
>>
>> Remember to import the Gorilla Mux package (github.com/gorilla/mux) before running this code.

Decent enough! Let me make some scope adjustments and stub the desired functionality. So now it looke like
```go
r.HandleFunc("/translate/{kind}/{name}", translatePokemon).Methods("GET")

kind := vars["kind"]
name := vars["name"]

if kind != "pokemon" { return some json error body and http.StatusBadRequest }

if description, err := getDescription(name); err != nil {
    // grab some official description
    return some error body and http.StatusFailedDependency }
} else if translation, err := getTranslation(description); err != nil {
    // use the Bard
    return some error body and http.StatusFailedDependency } 
else {
    // fianlly return the expectations
    return result(translation) and http.StatusOK 
}
```
I won't get into what errors and positive outcome is, basically a simpel json body with the inherite errors or data.


# Making of `getDescription`
> implement getDescription as a resty module request to pokemon API returning the description of the name or an error
>>
>> Certainly! Here's an example of how you can implement the getDescription function using the Resty module 
>> to make a request to the Pokemon API and retrieve the description of a given name:
```go
type PokemonAPIResponse struct {
	Species struct {
		Name        string `json:"name"`
		FlavorTexts []struct {
			FlavorText string `json:"flavor_text"`
			Language   struct {
				Name string `json:"name"`
			} `json:"language"`
		} `json:"flavor_text_entries"`
	} `json:"species"`
}

func getDescription(name string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetPathParams(map[string]string{
			"name": name,
		}).
		Get("https://pokeapi.co/api/v2/pokemon-species/{name}")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", errors.New("Failed to get description")
	}

	var apiResponse PokemonAPIResponse
	err = json.Unmarshal(resp.Body(), &apiResponse)
	if err != nil {
		return "", err
	}

	// Extract the description from the API response
	description := ""
	for _, entry := range apiResponse.Species.FlavorTexts {
		if entry.Language.Name == "en" {
			description = entry.FlavorText
			break
		}
	}

	if description == "" {
		return "", errors.New("No description found")
	}

	return description, nil
}
```
>> In this example, we define a struct PokemonAPIResponse to represent the response from the Pokemon API. 
>> We use the Resty module to make a GET request to the API using the provided name as a path parameter. 
>> We then parse the response and extract the English flavor text as the description.
>>
>> Please note that you need to import the Resty module (github.com/go-resty/resty/v2) before running this code.

Very nice. I'd like to make it slighlty mine though, so let's just call `SetResult(PokemonAPIResponse{}).` on the `client.R()` request to let the 
library do the json parsing. I don't care about of the `200` http code test for now, there is always room for improvement later on.

Impresive that:
 - it listed the right API endpoint by `pokemon` and `description` tokens only
 - it built path parameters correctly.
 - was able to build the reponse structure by API definition
   (wherever it took it from, online or "trained memory". I did not have any local referenc to that API).
 - it compiled and worked correctly (before my changes)
 - it filtered by English `en` langauge with passing this requirement. This context awareness is somewhat scary, 
   unless this functionality has been put forward in a lot of interviews :-) and is somewhat very well knwon and scores high on code references.


--- listing the rest of conversation

   implement getTranslation as an resty client to shakespeare style API translation endpoint of api.funtranslations.com

   Sure! Here's an example of how you can implement the getTranslation function using the Resty module to make a request to the Shakespeare style API translation endpoint of api.funtranslations.com:

   type TranslationAPIResponse struct {
	Contents struct {
		Translated string `json:"translated"`
	} `json:"contents"`
}

func getTranslation(text string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"text": text,
		}).
		Get("https://api.funtranslations.com/translate/shakespeare.json")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", errors.New("Failed to get translation")
	}

	var apiResponse TranslationAPIResponse
	err = json.Unmarshal(resp.Body(), &apiResponse)
	if err != nil {
		return "", err
	}

	translation := apiResponse.Contents.Translated

	if translation == "" {
		return "", errors.New("No translation found")
	}

	return translation, nil
}
In this example, we define a struct TranslationAPIResponse to represent the response from the translation API. We use the Resty module to make a GET request to the API, passing the text as a query parameter. We then parse the response and extract the translated text.

Please note that you need to import the Resty module (github.com/go-resty/resty/v2) before running this code.


implement json deserialiation for type StringOrInt. When string copy as is, when integer convert into string.

Sat, Oct 28, 2023, 2:06 pm
To implement JSON deserialization for a type StringOrInt where a string is copied as is and an integer is converted into a string, you can define a custom unmarshal method for the type. Here's an example:


Copied! ✅
📝 Copy Code
 ➕ Insert Code
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type StringOrInt struct {
	Value string
}

func (s *StringOrInt) UnmarshalJSON(data []byte) error {
	var rawValue interface{}
	err := json.Unmarshal(data, &rawValue)
	if err != nil {
		return err
	}

	switch v := rawValue.(type) {
	case string:
		s.Value = v
	case float64:
		s.Value = strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		s.Value = strconv.Itoa(v)
	default:
		return fmt.Errorf("unexpected value type: %T", v)
	}

	return nil
}

func main() {
	jsonData := `{"value": "hello"}`
	var obj StringOrInt
	err := json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", obj.Value)
	}

	jsonData = `{"value": 42}`
	err = json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", obj.Value)
	}
}
In this example, we define the StringOrInt struct with a Value field of type string. We then implement the UnmarshalJSON method for the StringOrInt type. Inside the method, we use type assertions to determine the type of the raw value and handle it accordingly. If the raw value is a string, we copy it as is. If it's a float64 or int, we convert it to a string using strconv. For any other type, we return an error.

In the main function, we demonstrate the usage of the StringOrInt type by unmarshaling JSON data containing a string and an integer. We print the Value field to verify the deserialization.

Please note that this implementation assumes that the JSON data has a field named value which corresponds to the Value field of the StringOrInt struct.



Sat, Oct 28, 2023, 2:06 pm
LL

refactor StringOrInt so is a simple type StringOrInt string redefinition. Unmarshal should accomodate this type overload

Sat, Oct 28, 2023, 2:11 pm
Certainly! Here's an updated version of the StringOrInt type with a custom UnmarshalJSON method that handles the overloaded type:


Copied! ✅
📝 Copy Code
 ➕ Insert Code
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type StringOrInt string

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

func main() {
	jsonData := `{"value": "hello"}`
	var obj struct {
		Value StringOrInt `json:"value"`
	}
	err := json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", obj.Value)
	}

	jsonData = `{"value": 42}`
	err = json.Unmarshal([]byte(jsonData), &obj)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value:", obj.Value)
	}
}
In this updated example, we redefine StringOrInt as a simple type of string. We then implement the UnmarshalJSON method for the StringOrInt type. Inside the method, we handle the overloaded type by switching on the type of the raw value and assigning it to the StringOrInt variable s. We convert float64 and int values to strings using strconv.

In the main function, we demonstrate the usage of the updated StringOrInt type by unmarshaling JSON data containing a string and an integer. We print the Value field to verify the deserialization.

Please note that this updated implementation assumes that the JSON data has a field named value which corresponds to the Value field of the StringOrInt struct.



Sat, Oct 28, 2023, 2:11 pm
LL

make apiErrorDetails implement the Error i/face

Sat, Oct 28, 2023, 2:13 pm
To make the apiErrorDetails struct implement the error interface in Go, you need to define an Error method on the struct. Here's an example:


Copied! ✅
📝 Copy Code
 ➕ Insert Code
package main

import (
	"fmt"
)

type apiErrorDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *apiErrorDetails) Error() string {
	return fmt.Sprintf("API Error: Code=%d, Message=%s", e.Code, e.Message)
}

func main() {
	err := &apiErrorDetails{
		Code:    404,
		Message: "Page not found",
	}

	fmt.Println(err.Error())
}
In this example, we define the apiErrorDetails struct with the Code and Message fields. We then implement the Error method on the struct, which returns a string representation of the error.

In the main function, we create an instance of apiErrorDetails with sample values, and then call the Error method on the error instance to print the error message.

By implementing the Error method on the apiErrorDetails struct, it now satisfies the error interface in Go.










