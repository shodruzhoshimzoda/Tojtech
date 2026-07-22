package main

import (
	"fmt"

	"github.com/shodruzhoshimzoda/tojtech/internal/config"
)


func main() {

	
	cfg := config.MustLoadConfig()
	
	fmt.Println(cfg)


}