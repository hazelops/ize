package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	var err error
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		VarEnv := getEnv("EXAMPLE_SECRET", "default-pass")
		VarSecret := getEnv("EXAMPLE_API_KEY", "default-pass")
		_, err = fmt.Fprintln(w, VarSecret, "Welcome to goblin! demo service👾")
		if err != nil {
			fmt.Printf("%s", err)
		}

		_, err = fmt.Fprintln(w, VarEnv, "Welcome to goblin! demo service👾")
		if err != nil {
			fmt.Printf("%s", err)
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
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
