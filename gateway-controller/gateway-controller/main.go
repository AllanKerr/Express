package main

import (
	"gateway-controller/cmd"
	"fmt"
)

func main() {
	fmt.Println("Execute")
	cmd.RootCmd.Execute()
}
