package main

import (
	"docker-search/cmd"
	_ "docker-search/cmd/search"
)

func main() {
	cmd.Execute()
}
