# Background story

I wanted to investigate how AI could help me in my coding endeavors. Using an older version of a dummy application, I started
replicating a simple web service that takes a Pokemon character name as parameter and returns a Shakespearean translation of its 
description. For reference, the  full project code is available at my [shakemon](https://github.com/LucaWolf/shakemon) repository.

I've had already some experience with the scope of this exercise, in particular of what APIs and modules 
([GorillaMux](https://github.com/gorilla/mux) and [Resty](https://github.com/go-resty/resty)) I wanted to use and the basic flow of the application. 

The only question remaining was how would [Codeium](https://codeium.com/) AI system help me when I threw at it all this information.

# Setup and usage
I have used the VS Code as text editor and the [plugin](https://codeium.com/vscode_tutorial) provided for it. Upon installing the extension,
the first use is going to be prompting you to login (I have used my google account) and then present you with an 
`{...}` icon in the extensions panel. Expand it and start interacting the prompt, it is as simple as that.

> **NOTE**
> not only Codeium generates blocks of code but also renders a decent description of what those source code do
and supplements it with some recommendations.


# The `main` entry point
 So let's start by issuing "demands":

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
>> When a GET request is made to that endpoint, the `getUser` function is called.
>> The function extracts the `name` and `age` path parameters from the request using `mux.Vars` and writes them back as the response.
>>
>> Remember to import the Gorilla Mux package (`github.com/gorilla/mux`) before running this code.

Decent enough! Let me make some scope adjustments and stub out the desired functionality. 
Basically I want to get a description by name and then a "classy" translation. So it then looked like this:

```go
// route declaration in main
r.HandleFunc("/translate/{kind}/{name}", translatePokemon).Methods("GET")

// route handler
func translatePokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kind := vars["kind"]
	name := vars["name"]

	if kind != "pokemon" { return some json error body and http.StatusBadRequest }

	// grab some official description
	if description, err := getDescription(name); err != nil {		
		return some error body and http.StatusFailedDependency }
	
	// use the Bard
	} else if translation, err := getTranslation(description); err != nil {		
		return some error body and http.StatusFailedDependency }

	// finally return the expectations
	else {		
		return result(translation) and http.StatusOK 
	}
}

```
I won't get into what errors and positive outcome are, basically simple json bodies with the inherited errors or data
(search `apiErrorDetails` and `apiReply` in the source code for reference).
What I needed now was the implementation of `getDescription` and `getTranslation`.


# Making of `getDescription`
I had issued this command and got back:

> implement getDescription as a resty module request to pokemon API returning the description of the name or an error
>>
>> Certainly! Here's an example of how you can implement the `getDescription` function using the Resty module 
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
>> In this example, we define a struct `PokemonAPIResponse` to represent the response from the Pokemon API. 
>> We use the Resty module to make a GET request to the API using the provided `name` as a path parameter. 
>> We then parse the response and extract the English flavor text as the description.
>>
>> Please note that you need to import the Resty module (`github.com/go-resty/resty/v2`) before running this code.

Very nice. I'd like to make it slightly mine though, so let's just call `SetResult(PokemonAPIResponse{})` on the `client.R()` request to let the 
library do the heavy lifting on the json parsing. I don't care about of the `200` http code test for now, there is always room for improvement later on.

 🤯 **Mind blown that**:
 - it listed the right API endpoint from input of `pokemon` and `description` tokens (keywords) only.
 - it built path parameters for the API request correctly.
 - was able to build the `PokemonAPIResponse` response structure compatible with the API's definition
   (wherever it took it from, online or "trained memory"). I did not have any local reference to that API though I have an older similar project in Github.
 - the embedded `struct` in `struct` definition of `PokemonAPIResponse` took me by surprise; I wasn't fully aware of this syntax being allowed.
 - the produced code compiled and worked correctly, even before adding any of my changes.
 - it automatically filtered the results by English `en` language <ins>without asking</ins> for this requirement 😲. This context awareness is somewhat scary, 
   unless this functionality has been put forward in a lot of interviews :-) and is likely very well known and scores high on token references.


# Making of `getTranslation`

> implement getTranslation as an resty client to shakespeare style API translation endpoint of api.funtranslations.com
>>
>> Sure! Here's an example of how you can implement the `getTranslation` function using the Resty module to make a request to 
>> the Shakespeare style API translation endpoint of api.funtranslations.com:

```go
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

```
>> In this example, we define a struct `TranslationAPIResponse` to represent the response from the translation API. 
>> We use the Resty module to make a GET request to the API, passing the `text` as a query parameter.
>>  We then parse the response and extract the translated text.
>>
>> Please note that you need to import the Resty module (`github.com/go-resty/resty/v2`) before running this code.

Perhaps was due to the my simple way of expressing the request but this API endpoint was supposed to take a POST request
(thus using  `client.R().SetFormData` method instead of `SetQueryParams`). I also adjusted the usual json response parsing and reworked a bit the return block.

Nonetheless, the provided code was a solid starting base. My adjusted version reads:

```go
{
	// POST instead of GET
	resp, err := client.R().
		SetFormData(map[string]string{
			"text": text,
		}).
		SetError(apiError{}).
		SetResult(TranslationAPIResponse{}).
		Post("https://api.funtranslations.com/translate/shakespeare.json")

	...
	// replace the http.Code 200 tests
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

```

# Playing with json unmarshalling
At this point, after a few debug/testing sessions, I realized that there is a divergence between the applications API error type `apiErrorDetails`
returned by the endpoint handler (as inspired by MS Graph and using `Code` formatted as `string` => some meaningful message) and 
the error returned by the translation service (similar but using `Code` as `integer` => the http.code value). 
In order to consolidate them, I've applied the "polymorphic" unmarshalling trick,
 i.e. type wrapping the `Code` into own type `StringOrInt` and implementing a bespoke parsing routine over it. 
 
So my request to the AI engine was:
> implement json deserialisation for type StringOrInt. When string copy as is, when integer convert into string.
>>
>> To implement JSON deserialization for a type `StringOrInt` where a string is copied as is and an integer is converted into a string, you can define a custom 
>> unmarshal method for the type. Here's an example:

```go 
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

```
>> In this example, we define the `StringOrInt` struct with a `Value` field of type `string`. We then implement the `UnmarshalJSON` method for the `StringOrInt`
>> type. Inside the method, we use type assertions to determine the type of the raw value and handle it accordingly. 
>> If the raw value is a string, we copy it as is. If it's a `float64` or `int`, we convert it to a string using `strconv`. 
>> For any other type, we return an error.
>>
>> In the main function, we demonstrate the usage of the StringOrInt type by unmarshaling JSON data containing a string and an integer. 
>> We print the Value field to >> verify the deserialization.
>>
>> Please note that this implementation assumes that the JSON data has a field named value which corresponds to the Value field of the StringOrInt struct.

Excellent implementation, with usage example on top! But me being picky, I wanted the `StringOrInt` type to exclude inner fields as I don't like multiple
 dereferencing (e.g. `obj.Value`). Therefore I issued a more explicit request:

> refactor StringOrInt so is a simple type StringOrInt string redefinition. Unmarshal should accommodate this type overload
>>
>>Certainly! Here's an updated version of the `StringOrInt` type with a custom `UnmarshalJSON` method that handles the overloaded type:

```go
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

// same usage example 
func main() {...}
```

>> In this updated example, we redefine `StringOrInt` as a simple type of `string`. We then implement the `UnmarshalJSON` method for the `StringOrInt` type. 
>> Inside the method, we handle the overloaded type by switching on the type of the raw value and assigning it to the `StringOrInt` variable `s`. 
>> We convert `float64` and `int` values to strings using `strconv`.


# Pass me the error
Finally, I wanted to have the `apiErrorDetails` (as possibly returned by the Shakespeare translating service) implement the Golang's error interface. 
Usually, this is useful when logging or bubbling up a bespoke type as an error. So here I went again (notice the casual "i/face" thrown at the AI engine):


> make apiErrorDetails implement the Error i/face
>>
>> To make the `apiErrorDetails` struct implement the `error` interface in Go, you need to define an `Error` method on the struct. Here's an example:

```go
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

```
>> In this example, we define the `apiErrorDetails` struct with the `Code` and `Message` fields. We then implement the `Error` method on the struct, 
>> which returns a string representation of the error.
>> 
>> In the `main` function, we create an instance of `apiErrorDetails` with sample values, and then call the `Error` method on 
>> the error instance to print the error message.
>>
>> By implementing the `Error` method on the `apiErrorDetails` struct, it now satisfies the `error` interface in Go.


# Conclusions
It is difficult and not accurate to draw ultimate conclusions based on single runs and such simple tests. But here are my initial impressions:

- I consider to have had a positive outcome of this AI interaction
- almost feels like a non-judgemental peer programming buddy
- several layers of re-factoring are available.
- it is fast. It may take seconds to crunch the data but then it outputs "many" (each engine has a hard limit) lines of code like nothing. How fast can you type?
- is up to you as user to ask for iterative changes, refactoring and honing down to an acceptable (correct, to your taste and liking) level of implementation.
- empirically, it saved me between 30% to 50% of development time. Realistically, would be even more when it gets down to writing documentation, 
 unit tests and other time consuming cruft.
- needless to say but don't trust the generated output 100%. Always validate and review the produced code. 
  It seems very good at seeing "the trees for the forest" but not so good as seeing "the forest" or the whole ecosystem.
- would I use it for personal (toy) projects and learning experience? Heck yes!
- would I use this (or alternatives) for work? Absolutely, subject to the powers to be adopting such an AI provider.
- beware addictions: it seems so fast and helpful in assisting with refactoring that you may be tempted to jump straight into it and 
  attempt "changing the world" (aka your project) . Whilst the whole evolution and even the learning model is based on iterations, 
  not all refactoring is good or worthwhile... so get a life people.








