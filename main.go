package main

import (
	"fmt"

	"github.com/Gergenus/GoMockServer/src/config"
)

func main() {
	cfg, err := config.LoadConfig("./conf.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
