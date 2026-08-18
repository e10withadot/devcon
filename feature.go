package main

import (
    "fmt"
    "log"
    "os"
)

func feature(args []string) {
    if len(args) < 1 {
        log.Fatal("Subcommand required.")
    }
    cargs := os.Args[2:]
    switch cmd := os.Args[1]; cmd {
    case "add":
        add(cargs)
    case "rm":
        rm(cargs)
    default:
        fmt.Println("Invalid command:", cmd)
    }
}

func add(args []string) {
}

func rm(args []string) {
}
