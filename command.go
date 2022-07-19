package main

import (
	"fmt"
	"os"
	"strconv"
)

type Command struct {
	Name string
	Do   func(params []string, rawInput []byte)
}

var TIMO_COMMANDS = []Command{}
var TIMO_COMMANDS_LOADER = func() bool {
	TIMO_COMMANDS = append(
		TIMO_COMMANDS,
		Command{
			Name: "test",
			Do: func(params []string, rawInput []byte) {
				if len(params) < 1 {
					os.Stdout.Write(append([]byte("no params"), 0x0a))
				} else {
					fmt.Println(params)
				}
			},
		},
		Command{
			Name: "exit",
			Do: func(params []string, rawInput []byte) {
				if len(params) < 1 {
					os.Exit(0)
				}
				i, err := strconv.Atoi(params[0])
				if err != nil {
					os.Exit(-1)
				}
				os.Exit(i)
			},
		},
		Command{
			Name: "motd",
			Do: func(params []string, rawInput []byte) {
				GenMOTD()
			},
		},
	)
	return true
}()
