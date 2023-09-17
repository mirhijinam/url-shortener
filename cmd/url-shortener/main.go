package main

import (
	"fmt"

	"github.com/mirhijinam/url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
