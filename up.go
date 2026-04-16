package main

import (
    "fmt"
    "log"
    "os/exec"
)

func up(args []string) {
    if len(args) < 1 {
        log.Fatal("Path required.")
    }
    cmd := exec.Command("devcontainer", "up", "--workspaceFolder", args[0])
    var err error; var stdout []byte
    if stdout, err = cmd.Output(); err != nil {
        log.Fatal(err)
    } else {
        fmt.Println(string(stdout))
    }
}
