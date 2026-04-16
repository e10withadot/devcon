package main

import (
    "fmt"
    "log"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Subcommand required.")
    }
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
