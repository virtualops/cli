package main

import (
	"fmt"
	"github.com/virtualops/cli/cmd"
)

var OauthSecret string

func main() {
	fmt.Println("test", OauthSecret)
	cmd.Execute()
}
