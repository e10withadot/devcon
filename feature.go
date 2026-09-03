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
    cargs := args[1:]
    switch cmd := args[0]; cmd {
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
    var dc_path string = args[1]
    data, err := os.ReadFile(dc_path + "/.devcontainer/devcontainer.json")
    if err != nil {
        log.Fatal("Could not read devcontainer.json.")
    }
    var dcj map[string]any
    if err = json.Unmarshal(data, &dcj); err != nil {
        log.Fatal("The devcontainer.json is invalid.")
    }
    var features map[string]any; var ok bool
    if features, ok = dcj["features"].(map[string]any); !ok {
        features = map[string]any{}
    }
    features[feature] = map[string]any{}
    dcj["features"] = features
    if data, err = json.MarshalIndent(&dcj, "", "   "); err != nil {
        log.Fatal("Devcon could not reformat into json. Please report this error.")
    }
    if err = os.WriteFile(dc_path + "/.devcontainer/devcontainer.json", data, 0644); err != nil {
        log.Fatal("Error writing to file.")
    }
}

func rm(args []string) {
    if len(args) < 2 {
        log.Fatal("Correct syntax: devcon feature rm [feature] [path]")
    }
    var feature string = args[0]
    var dc_path string = args[1]
    data, err := os.ReadFile(dc_path + "/.devcontainer/devcontainer.json")
    if err != nil {
        log.Fatal("Could not read devcontainer.json.")
    }
    var dcj map[string]any
    if err = json.Unmarshal(data, &dcj); err != nil {
        log.Fatal("The devcontainer.json is invalid.")
    }
    var features map[string]any; var ok bool
    if features, ok = dcj["features"].(map[string]any); !ok {
        log.Fatal("No features found in devcontainer.json.")
    }
    if _, exists := features[feature]; !exists {
        log.Fatal("Feature not found: " + feature)
    }
    delete(features, feature)
    dcj["features"] = features
    if data, err = json.MarshalIndent(&dcj, "", "   "); err != nil {
        log.Fatal("Devcon could not reformat into json. Please report this error.")
    }
    if err = os.WriteFile(dc_path + "/.devcontainer/devcontainer.json", data, 0644); err != nil {
        log.Fatal("Error writing to file.")
    }
}
