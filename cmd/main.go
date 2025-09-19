package main

import (
	"fmt"
	"light-defender-client/pkg/config"
)

func main() {
	_, err := config.NewConfig()
	if err != nil {
		fmt.Println(err)
	}
}
