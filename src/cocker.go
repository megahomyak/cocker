package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type command struct {
	name, contents string
}

type context struct {
	cacheDirectory, finalContainerName, programInterpreterPath string
}

func executeCommand()

func (ctx *context) executeSource(program string) string {
	// TODO: create descriptor, run the program, get the container name (don't forget to .TrimSpace()) or panic if there's none
}

func (ctx *context) overlay(sourceContainerName string, newContainerName string) {
	// TODO: just one command to lxc here
}

func (ctx *context) executeLayer(currentContainerName string, program string) {
	// TODO: format the command with the current container name, run the program
}

func main() {
	// TODO: add layer caching
	data, err := os.ReadFile("Cockerfile")
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(data), "\n")

	commands := []command{}

	ctx := context{
		cacheDirectory:         *flag.String("cache-directory", "/var/lib/lxc", "The directory where Cockerfiles are cached"),
		finalContainerName:     flag.Arg(0),
		programInterpreterPath: *flag.String("program-interpreter-path", "/usr/bin/bash", "A path to the binary that will execute all of your Cockerfile commands"),
	}
	if ctx.finalContainerName == "" {
		log.Fatalln("Please, provide the final container name as the first positional argument")
	}

	if err := os.MkdirAll(ctx.cacheDirectory, 0660); err != nil {
		log.Fatalln(err)
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		contentLine, isSuccess := strings.CutPrefix(line, ">")
		if isSuccess {
			if len(commands) == 0 {
				log.Fatalln("Contents of a command should come after the name of the command")
			}
			if commands[len(commands)-1].contents != "" {
				commands[len(commands)-1].contents += "\n"
			}
			commands[len(commands)-1].contents += contentLine
		} else {
			commands = append(commands, command{name: trimmed, contents: ""})
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
	currentContainerName := ctx.executeSource(commands[i].contents)
	i++
	layerCount := 0
	for _, command := range commands[i:] {
		switch command.name {
		case "LAYER":
			layerCount++
			newContainerName := fmt.Sprintf("%s-layer-%s", ctx.finalContainerName, layerCount)
			ctx.overlay(currentContainerName, newContainerName)
			currentContainerName = newContainerName
			ctx.executeLayer(command.contents, currentContainerName)
		default:
			log.Fatalf("Unknown command name: \"%s\"\n", command.name)
		}
	}
	ctx.overlay(currentContainerName, ctx.finalContainerName)
}
