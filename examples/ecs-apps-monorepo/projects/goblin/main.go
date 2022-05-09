package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		VarEnv := getEnv("EXAMPLE_SECRET", "default-pass")
		VarSecret := getEnv("EXAMPLE_API_KEY", "default-pass")
		fmt.Fprintln(w, VarSecret, "Welcome to goblin! demo service👾")
		fmt.Fprintln(w, VarEnv, "Welcome to goblin! demo service👾")
	})
	println("Goblin demo service started on port 3000")
	http.ListenAndServe(":3000", nil)
}

// getEnv get key environment variable if exist otherwise return defalutValue
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
