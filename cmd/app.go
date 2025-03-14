package main

import (
	"fmt"
	"github.com/sqc157400661/jobx/cmd/server"
)

func main() {
	err := server.StartServer()
	fmt.Println(err)
}
