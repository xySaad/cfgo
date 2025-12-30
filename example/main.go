//go:generate cfgo config/config.json config/config.go
package main

import (
	"fmt"

	"example/config"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println(cfg.Graphql.Password)
}
