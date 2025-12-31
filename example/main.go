//go:generate cfgo
package main

import (
	"fmt"

	"example/config"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println(cfg.Graphql.Password)
}
