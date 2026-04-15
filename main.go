package main

import (
    "os"
    "fmt"
)

func main() {
    cargs := os.Args[2:]
    switch cmd := os.Args[1]; cmd {
    case "create":
        create(cargs)
    case "up":
        fmt.Println("Up command called!")
    default:
        fmt.Println("Invalid command:", cmd)
    }
}
