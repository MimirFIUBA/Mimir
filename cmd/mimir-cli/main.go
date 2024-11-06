// client.go
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Conectar al servidor TCP en localhost:8080
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run client.go <comando>")
		return
	}

	conn, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		fmt.Println("Error al conectarse al servidor:", err)
		return
	}
	defer conn.Close()

	command := strings.Join(os.Args[1:], " ") + "\n"

	_, err = conn.Write([]byte(command))
	if err != nil {
		fmt.Println("Error al enviar comando:", err)
		return
	}
}
