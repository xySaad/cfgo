//go:generate ../build/cfgo
package main

import (
	config "example/config/generated"
	"fmt"
)

func main() {
	cfg := config.GetConfig()
	fmt.Println(cfg.Graphql.Password)
}
