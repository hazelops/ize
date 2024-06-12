package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Data struct {
	message       string
	exampleSecret string
	exampleApiKey string
}

func main() {
	var err error
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		exampleSecret := getEnv("EXAMPLE_SECRET", "default-pass")
		exampleApiKey := getEnv("EXAMPLE_API_KEY", "default-pass")

		d := Data{
			message:       "Welcome to goblin demo service 👾",
			exampleSecret: exampleSecret,
			exampleApiKey: exampleApiKey,
		}

		jData, err := json.MarshalIndent(d, "", "  ")
		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(jData)
		if err != nil {
			fmt.Println(err)
		}
	})

	fmt.Printf("Goblin demo service started on port 3000")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

// getEnv get key environment variable if exist otherwise return defalutValue
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
