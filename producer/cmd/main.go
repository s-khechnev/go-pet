package main

import (
	"fmt"
	"producer/internal/config"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println("Config:", cfg)
}
