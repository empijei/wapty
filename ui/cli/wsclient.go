package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
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
	suppCommands[intercept.EDITED.String()] = edit
	suppCommands[intercept.FORWARDED.String()] = forward
}

func main() {
	go wsLoop()
	go interact()
	mainLoop()
}

func mainLoop() {
	for {
		select {
		case cmd := <-serverChannel:
			_ = ioutil.WriteFile("tmp.swp", *cmd.Payload, 0644)
			log.Println("Payload intercepted, edit it and press enter to continue.")

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
	fmt.Println("TODO implement help")
}

func parseInput(in string) (ui.Command, error) {
	commands := strings.Split(in, " ")
	directive, ok := suppCommands[commands[0]]
	if !ok {
		for cmd, dir := range suppCommands {
			if strings.HasPrefix(cmd, commands[0]) {
				directive = dir
				return directive(commands), nil
			}
		}
	} else {
		return directive(commands), nil
	}
	return ui.Command{}, fmt.Errorf("Command not found")
}

func edit(commands []string) ui.Command {
	var payload []byte
	payload, _ = ioutil.ReadFile("tmp.swp") //TODO chech this error
	args := ui.Args(map[string]string{"action": intercept.EDITED.String()})
	return ui.Command{Args: args, Channel: intercept.EDITORCHANNEL, Payload: &payload}
}

func forward(commands []string) ui.Command {
	args := ui.Args(map[string]string{"action": intercept.FORWARDED.String()})
	return ui.Command{Args: args, Channel: intercept.EDITORCHANNEL}
}

func interact() {
	for stdin.Scan() {
		cliChannel <- stdin.Text()
	}
}

func cli() {
	for cmd := range serverChannel {
		_ = ioutil.WriteFile("tmp.swp", *cmd.Payload, 0644)
		log.Println("Payload intercepted, edit it and press enter to continue.")
		var payload []byte
		var args ui.Args
		for args == nil {
			stdin.Scan()
			switch stdin.Text() {
			case intercept.EDITED.String():
				payload, _ = ioutil.ReadFile("tmp.swp") //TODO chech this error
				args = ui.Args(map[string]string{"action": intercept.EDITED.String()})
			case intercept.FORWARDED.String():
				args = ui.Args(map[string]string{"action": intercept.FORWARDED.String()})
			default:
				log.Println("Unknown command")
				log.Println("Try with ", intercept.EDITED, " ", intercept.FORWARDED)
			}
		}
		log.Println("Continued")
		err := websocket.JSON.Send(ws, ui.Command{Args: args, Channel: intercept.EDITORCHANNEL, Payload: &payload})
		if err != nil {
			panic(err)
		}
	}
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
