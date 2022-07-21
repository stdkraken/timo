package main

import (
	"encoding/json"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore" // more use cases lol
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type PluginInfo struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Version         string   `json:"version"`
	Main            string   `json:"main"`
	Author          string   `json:"author"`
	Authors         []string `json:"authors"`
	RunsInGoRoutine bool     `json:"runs-in-go-routine"`
	addedCommands   []Command
}

var VMs map[*PluginInfo]*otto.Otto = make(map[*PluginInfo]*otto.Otto)
var loadedPlugins = []*PluginInfo{}

func UnloadPlugins() {
	VMs = make(map[*PluginInfo]*otto.Otto) // unset all functions, added by the plugins
	for _, plugin := range loadedPlugins {
		for _, cmd := range plugin.addedCommands {
			newCMDs := []Command{}
			for _, oldCmd := range TIMO_COMMANDS {
				if oldCmd.Name != cmd.Name {
					newCMDs = append(newCMDs, oldCmd)
				} else {
					fmt.Println("Removed function (" + oldCmd.Name + ")")
				}
			}
			TIMO_COMMANDS = newCMDs
		}
	}
}

func PluginInfoByVM(oVM *otto.Otto) *PluginInfo {
	for info, vm := range VMs {
		if oVM == vm {
			return info
		}
	}
	return nil
}

func InitVM(plugin *PluginInfo) *otto.Otto {
	vm := otto.New()
	err := vm.Set("osExit", func(call otto.FunctionCall) otto.Value {
		i, err := call.Argument(0).ToInteger()
		if err != nil {
			log.Println(err.Error())
		}
		os.Exit(int(i))
		return otto.Value{}
	})

	if err != nil {
		fmt.Println("Couldn't set osExit")
	}

	err = vm.Set("downloadFile", func(call otto.FunctionCall) otto.Value {
		AutoPath := func(URL string) string {
			return strings.Split(URL, "/")[len(strings.Split(URL, "/"))-1]
		}
		URL, err := call.Argument(0).ToString()
		if err != nil || URL == "undefined" {
			return otto.Value{}
		}
		Path, err := call.Argument(1).ToString()
		if err != nil || Path == "undefined" {
			Path = AutoPath(URL)
		}
		Agent, err := call.Argument(2).ToString()
		if err != nil || Agent == "undefined" {
			Agent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
		}
		c := &http.Client{}
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			return otto.Value{}
		}
		req.Header.Set("User-Agent", Agent)
		resp, err := c.Do(req)
		if err != nil {
			return otto.Value{}
		}
		outputFile, err := os.Create(Path)
		if err != nil {
			return otto.Value{}
		}
		_, err = io.Copy(outputFile, resp.Body)
		if err != nil {
			return otto.Value{}
		}
		return otto.Value{}
	})

	err = vm.Set("setInputFg", func(call otto.FunctionCall) otto.Value {
		i, err := call.Argument(0).ToInteger()
		if err == nil {
			ShellInputFg = int(i)
		}
		return otto.Value{}
	})

	err = vm.Set("sleep", func(call otto.FunctionCall) otto.Value {
		i, err := call.Argument(0).ToInteger()
		if err == nil {
			time.Sleep(time.Duration(i) * time.Millisecond)
		}
		return otto.Value{}
	})

	err = vm.Set("addCustomCommand", func(call otto.FunctionCall) otto.Value {
		cmd := Command{
			Name: call.Argument(0).String(),
			Do: func(params []string, rawInput []byte) {
				raw, _ := otto.ToValue(string(rawInput))
				call.Argument(1).Call(otto.Value{}, raw)
			},
		}
		PluginInfoByVM(vm).addedCommands = append(PluginInfoByVM(vm).addedCommands, cmd)
		TIMO_COMMANDS = append(TIMO_COMMANDS, cmd)

		return otto.Value{}
	})

	err = vm.Set("execute", func(call otto.FunctionCall) otto.Value {
		if call.Argument(0).IsString() && call.Argument(1).IsNumber() {
			args := []string{}
			// script path is given
			i := 2
			for call.Argument(i).String() != "undefined" {
				// new arg
				args = append(args, call.Argument(i).String())
				i++
			}
			v, err := call.Argument(1).ToInteger()
			if err != nil {
				fmt.Println("Converting otto.Value to Integer failed.")
				return otto.Value{}
			}
			rawBytes := func() (p []byte) {
				first := true
				i := 2
				for call.Argument(i).String() != "undefined" {
					if first {
						first = false
					} else {
						p = append(p, append([]byte{' '}, []byte(call.Argument(i).String())...)...)
					}
					i++
				}
				return
			}

			// check local timo commands
			wasLocal := false
			for _, localCmd := range TIMO_COMMANDS {
				if localCmd.Name == call.Argument(0).String() {
					localCmd.Do(args, rawBytes())
					wasLocal = true
				}
			}
			if !wasLocal {
				for _, d := range SPLITTED_PATH {
					binPath := d + Sep() + call.Argument(0).String()
					if f, err := os.Open(binPath); err == nil {
						f.Close()
						cmd := exec.Command(binPath, args...)
						if v&1 == 0 { // not first bit set <=> 0/2/4
							cmd.Stdout = os.Stdout
						}
						if v&2 == 0 { // not second bit set <=> 0/1/4
							cmd.Stdin = os.Stdin
						}
						cmd.Stderr = os.Stdout
						CurrentCommand = cmd
						// cmd.Dir = "./"
						err = cmd.Run()
						if err != nil {
							fmt.Println("Couldn't start command command.")
						}
						break
					}
				}
			}
		}
		return otto.Value{}
	})
	VMs[plugin] = vm
	return vm
}

func LoadPlugins() {
	stuff, err := ioutil.ReadDir("./plugins/")
	if err != nil {
		fmt.Println("Couldn't open plugins folder.")
		return
	}
	for _, fi := range stuff {
		if fi.IsDir() {
			// check for plugin.json
			if f, err := os.Open("./plugins/" + fi.Name() + "/plugin.json"); err == nil {
				data, err := ioutil.ReadAll(f)
				if err != nil {
					fmt.Println("Read of plugin.json failed:", err.Error())
				}
				pInfo := &PluginInfo{}
				err = json.Unmarshal(data, pInfo)
				if err != nil {
					fmt.Println("Failed running json.Unmarshal:", err.Error())
				}
				if srcFile, err := os.Open("./plugins/" + fi.Name() + "/" + pInfo.Main); err == nil {
					fmt.Print("Running ./plugins/" + fi.Name() + "/" + pInfo.Main)
					data, _ = ioutil.ReadAll(srcFile)
					vm := InitVM(pInfo)
					loadedPlugins = append(loadedPlugins, pInfo)

					if pInfo.RunsInGoRoutine {
						tm.Println(" in routine.")
						go vm.Run(string(data))
					} else {
						tm.Println(" in function.")
						_, err = vm.Run(string(data))
					}
					srcFile.Close()
				}
				f.Close()
			}
		}
	}
	time.Sleep(1 * time.Second)

}

var PluginLoader = func() bool {
	LoadPlugins()
	return true
}()
