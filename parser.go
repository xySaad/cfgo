package main

import (
	"fmt"
	"regexp"
)

var envPattern = regexp.MustCompile(`@ENV.(\S+)`)

func parseImports(path string) (string, bool) {
	matches := envPattern.FindStringSubmatch(path)
	if len(matches) > 1 {
		return fmt.Sprintf(`os.Getenv("%s")`, matches[1]), true
	}
	return fmt.Sprintf(`"%s"`, path), false
}
