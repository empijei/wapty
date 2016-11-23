package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/empijei/Wapty/intercept"
	"github.com/empijei/Wapty/ui"
	"golang.org/x/net/websocket"
)

var ws *websocket.Conn
var serverChannel chan ui.Command
var cliChannel chan string
var stdin *bufio.Scanner
var suppCommands map[string]dir

type dir func([]string) ui.Command

func init() {
	serverChannel = make(chan ui.Command)
	cliChannel = make(chan string)
	stdin = bufio.NewScanner(os.Stdin)
	suppCommands = make(map[string]dir)
	suppCommands[intercept.EDIT.String()] = edit
	suppCommands[intercept.FORWARD.String()] = forward
	suppCommands[intercept.FETCH.String()] = fetch
	suppCommands["intercept"] = setIntercept
}

func main() {
	go wsLoop()
	go interact()
	mainLoop()
}

func mainLoop() {
	for {
		prompt()
		select {
		case cmd := <-serverChannel:
			switch cmd.Channel {
			case intercept.EDITORCHANNEL:
				_ = ioutil.WriteFile("tmp.swp", cmd.Payload, 0644)
				fmt.Println("\nPayload intercepted, edit it and press enter to continue.")
			case intercept.HISTORYCHANNEL:
				handleHistory(cmd)
			case intercept.SETTINGSCHANNEL:
				handleSettings(cmd)
			}

		case input := <-cliChannel:
			cmd, err := parseInput(input)
			if err != nil {
				fmt.Println(err.Error())
				printHelp()
			}
			err = websocket.JSON.Send(ws, cmd)
			if err != nil {
				panic(err)
			}
		}
	}
}

func printHelp() {
	fmt.Println("Supported commands:")
	for cmd, _ := range suppCommands {
		fmt.Println(cmd)
	}
}

func parseInput(in string) (ui.Command, error) {
	commands := strings.Split(in, " ")
	directive, ok := suppCommands[commands[0]]
	if !ok {
		for cmd, dir := range suppCommands {
			if strings.HasPrefix(cmd, commands[0]) {
				directive = dir
				if ok == true {
					return ui.Command{}, fmt.Errorf("Ambiguous command")
				}
				ok = true
			}
		}
		if directive != nil {
			return directive(commands), nil
		}
	} else {
		return directive(commands), nil
	}
	return ui.Command{}, fmt.Errorf("Command not found")
}

func edit(commands []string) ui.Command {
	var payload []byte
	payload, _ = ioutil.ReadFile("tmp.swp") //TODO chech this error
	return ui.Command{Action: intercept.EDIT.String(), Channel: intercept.EDITORCHANNEL, Payload: payload}
}

func forward(commands []string) ui.Command {
	return ui.Command{Action: intercept.FORWARD.String(), Channel: intercept.EDITORCHANNEL}
}

func fetch(commands []string) ui.Command {
	return ui.Command{Action: intercept.FETCH.String(), Channel: intercept.HISTORYCHANNEL}
}

func setIntercept(commands []string) ui.Command {
	if len(commands) == 1 {
		return ui.Command{Action: "intercept", Channel: intercept.SETTINGSCHANNEL}
	}
	value := "false"
	if strings.HasPrefix("true", commands[1]) || strings.HasPrefix("on", commands[1]) {
		value = "true"
	}
	return ui.Command{Action: "intercept", Channel: intercept.SETTINGSCHANNEL, Args: []string{value}}
}

func interact() {
	for stdin.Scan() {
		cliChannel <- stdin.Text()
	}
}

func handleSettings(cmd ui.Command) {
	switch cmd.Action {
	case "intercept":
		fmt.Println("Intercept is " + cmd.Args[0])
	}
}

func handleHistory(cmd ui.Command) {
}

func prompt() {
	fmt.Printf("wapty-cli >> ")
}

func wsLoop() {
	var url = "ws://localhost:8081/ws"
	var origin = "http://localhost/"
	var err error
	ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}
	for {
		var msg ui.Command
		err := websocket.JSON.Receive(ws, &msg)
		if err != nil {
			panic(err)
		}
		serverChannel <- msg
	}

}
