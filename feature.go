package main

import (
	"encoding/json"
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
    case "rm", "remove":
        rm(cargs)
    default:
        fmt.Println("Invalid command:", cmd)
    }
}

func add(args []string) {
    if len(args) < 2 {
        log.Fatal("Correct syntax: devcon feature add [feature] [path]")
    }
    var feature string = args[0]
    var path string = args[1]
    data, err := os.ReadFile(path + "/.devcontainer/devcontainer.json")
    if err != nil {
        log.Fatal("Could not read devcontainer.json.")
    }
    var dcj map[string]any
    if err = json.Unmarshal(data, &dcj); err != nil {
        log.Fatal("The devcontainer.json is invalid.")
    }
    features, ok := dcj["features"].([]any)
    if !ok {
        features = []any{}
    }
    features = append(features, map[string]any{ feature: map[string]any{}, })
    dcj["features"] = features
    if data, err = json.Marshal(&dcj); err != nil {
        log.Fatal("Devcon could not reformat into json. Please report this error.")
    }
    os.WriteFile(path + "/.devcontainer/devcontainer.json", data, 0666)
}

func rm(args []string) {
    if len(args) < 2 {
        log.Fatal("Correct syntax: devcon feature rm [feature] [path]")
    }
}
