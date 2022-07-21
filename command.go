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
			Name: "input-multicolor",
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
		Command{
			Name: "apt-check-toggle",
			Do: func(params []string, rawInput []byte) {
				AptPackageCheck = !AptPackageCheck
			},
		},
		Command{
			Name: "timo",
			Do: func(params []string, rawInput []byte) {

			},
		},
		Command{
			Name: "plreload", // reload all the plugins
			Do: func(params []string, rawInput []byte) {
				// unload the plugins => undo all the changes did by the plugins
				UnloadPlugins()
				LoadPlugins()
			},
		},
		Command{
			Name: "plload", // load all the plugins
			Do: func(params []string, rawInput []byte) {
				LoadPlugins()
			},
		},
		Command{
			Name: "plunload", // unload all the plugins
			Do: func(params []string, rawInput []byte) {
				UnloadPlugins()
			},
		},
	)
	return true
}()
