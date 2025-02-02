package main

import (
	"log"
	"os"
	"strings"
)

/*
Read the file by lines
Parse lines into array[command]
Execute the commands one by one, passing the new names forward, switching layers
*/

type command struct {
    name, contents string
}

func main() {
    data, err := os.ReadFile("Cockerfile")
    if err != nil {
        log.Fatal(err)
    }

    lines := strings.Split(string(data), "\n")

    commands := []command{}

    for _, line := range lines {
        if strings.TrimSpace(line) == "" {
            continue
        }
        contentLine, isSuccess := strings.CutPrefix(line, ">")
        if isSuccess {
            if len(commands) == 0 {
                log.Fatal("Contents of a command should come after the name of the command")
            }
            if commands[len(commands) - 1].contents != "" {
                commands[len(commands) - 1].contents += "\n"
            }
            commands[len(commands) - 1].contents += contentLine
        } else {
            trimmed := strings.TrimSpace(line)
            commands = append(commands, command{ name: trimmed, contents: "" })
        }
    }

    i := 0
    if len(commands) <= i || commands[i].name != "COCKERFILE v1" || commands[i].contents != "" {
        log.Fatal("The first line of the Cockerfile should be its heading without any content lines")
    }
    i++
    if len(commands) <= i || commands[i].name != "SOURCE" {
        log.Fatal("The first command should be SOURCE")
    }
    i++
}
