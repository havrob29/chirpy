package main

import (
	"fmt"
	"net/http"
)

func main() {
	httpHandler := http.NewServeMux()
	fmt.Println(httpHandler)
}
