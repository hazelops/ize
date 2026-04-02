package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Printf("Cron job started at %s\n", time.Now().Format(time.RFC3339))
	fmt.Printf("APP_NAME=%s\n", os.Getenv("APP_NAME"))
	fmt.Println("Cron job completed successfully")
}
