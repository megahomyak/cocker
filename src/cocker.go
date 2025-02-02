package main

import (
	"log"
	"os"
	"strings"
)

/*

 */

type command struct {
    name, contents string
}

func executeSource(source string) string {

}

func main() {
    data, err := os.ReadFile("Cockerfile")
    if err != nil {
        log.Fatalln(err)
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
                log.Fatalln("Contents of a command should come after the name of the command")
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
        log.Fatalln("The first line of the Cockerfile should be its heading without any content lines")
    }
    i++
    if len(commands) <= i || commands[i].name != "SOURCE" {
        log.Fatalln("The first command should be SOURCE")
    }
    currentContainerName := executeSource(commands[i].contents)
    i++
    for _, command := range commands[i:] {
        switch command.name {
        case "LAYER":
            // Create new container layer from currentContainerName
            // Execute the commands
        default:
            log.Fatalf("Unknown command name: \"%s\"\n", command.name)
        }
    }
    // Create final ("live") container with the specified name
}
