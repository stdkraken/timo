package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gookit/color"

	tm "github.com/buger/goterm"
)

func BufferCut(buffer []byte, Size int) []byte { // removes elements from buffer until size is reached
	for len(buffer) > Size {
		buffer = buffer[:len(buffer)-1]
	}
	return buffer
}

var BUFFER_SIZE = 1024 * 64 // longest command
var PATH = os.Getenv("PATH")
var SPLITTED_PATH = func() []string {
	switch runtime.GOOS {
	case "windows":
		// idk how
	default:
		return append(strings.Split(PATH, ":"), "./") // add current directory to the paths
	}
	return []string{}
}()
var VERSION = "0.2.2"

var CurrentCommand *exec.Cmd

func Sep() string { // get system default seperator
	switch runtime.GOOS {
	case "windows":
		return "\\"
	default:
		return "/"
	}
}

type NullOutput struct{}

func (NO NullOutput) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func CommandCleaner() {
	for {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		go func() {
			if CurrentCommand != nil {
				// a process is running
				err := CurrentCommand.Process.Signal(os.Interrupt)
				if err != nil {
					println("Couldn't send the interrupt signal to the process.")
				}
				err = CurrentCommand.Wait()
				if err != nil {
					println("Couldn't wait for the process to stop.")
				}
				CurrentCommand = nil
			} else {
				// no process is running
				println()
				os.Exit(0)
			}
			/*
				if err := CurrentCommand.Process.Kill(); err != nil {
					log.Fatal("failed to kill process: ", err)
				}
			*/
		}()
	}
}

func GenMOTD() {
	motds := []func(){
		func() {
			color.Style{color.FgYellow, color.OpBold}.Println("Did you know, that our name was inspired by the german football player Timo Werner?")
		},
		func() {
			color.Style{color.FgLightRed, color.OpBold}.Println("Did you know, that you can fully customize this terminal?")
		},
	}
	seed := int64(time.Now().Day())
	seed += int64(time.Now().Month() << 8)
	seed += int64(time.Now().Year() << 8)
	motd := motds[rand.Intn(len(motds))]
	motd()
}

var ShellInputFg = 0x5F
var ShellIconFg = 0x5F

func main() {
	go CommandCleaner()
	log.SetOutput(&NullOutput{})
	tm.Clear() // Clear current screen
	tm.MoveCursor(1, 1)
	log.Println(SPLITTED_PATH)
	GenMOTD()
	log.Println("Launching shell.")
	for {
		color.Style{color.Color(ShellIconFg), color.OpBold}.Print("⋙ ")

		inputStyle := color.Style{color.Color(ShellInputFg)}
		PrintColorMod(inputStyle, " ")

		buffer := ReadStdin()
		log.Println(buffer)

		fmt.Print("\033[0m") // reset color

		// start timer
		s := time.Now()

		// run command

		cmdName := ""
		cmdArgs := []string{}
		cmdRaw := buffer

		if len(buffer) != 0 {
			hasSpace := false
			// check, if it has space
			for _, c := range buffer {
				if c == 0x20 {
					hasSpace = true
					break
				}
			}
			if !hasSpace {
				cmdName = string(buffer)
			} else {
				cmdName = strings.Split(string(buffer), " ")[0]
				cmdArgs = strings.Split(string(buffer), " ")[1:] // remove first element
			}
		} else {
			continue
		}

		// check local timo commands
		wasLocal := false
		wasCommand := false
		for _, localCmd := range TIMO_COMMANDS {
			if localCmd.Name == cmdName {
				localCmd.Do(cmdArgs, cmdRaw)
				wasLocal = true
				wasCommand = true
			}
		}
		if !wasLocal {
			for _, d := range SPLITTED_PATH {
				binPath := d + Sep() + cmdName
				if f, err := os.Open(binPath); err == nil {
					wasCommand = true
					f.Close()
					cmd := exec.Command(binPath, cmdArgs...)
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stderr = os.Stdout
					CurrentCommand = cmd
					// cmd.Dir = "./"
					err = cmd.Run()
					if err != nil {
						log.Println(err.Error())
					}
					break
				}
			}
		}
		if !wasCommand {
			print("The command \"")
			color.Style{color.FgRed, color.OpBold}.Print(cmdName)
			println("\" wasn't found.")
			names := []string{}
			for _, d := range SPLITTED_PATH {
				if f, err := os.Open(d); err == nil {
					f.Close()
					if dir, err := ioutil.ReadDir(d); err == nil {
						for _, fi := range dir {
							if !fi.IsDir() {
								// is file
								isExecutable := func(Name string) bool {
									switch runtime.GOOS {
									case "windows":
										return strings.HasSuffix(Name, ".exe")
									default:
										return !strings.Contains(Name, ".")
									}
								}
								// is binary
								if isExecutable(fi.Name()) {
									names = append(names, fi.Name())
								}
							}
						}
					}
				}
			}
			for _, tCommand := range TIMO_COMMANDS {
				names = append(names, tCommand.Name)
			}
			cmpm := NewCmpMap(names, cmdName)
			nearCmdName, prob := cmpm.Nearest()
			if nearCmdName != "" {
				print("Did you mean \"")
				color.Style{color.FgGreen, color.OpBold}.Print(nearCmdName)
				formattedProb := fmt.Sprint(prob * 100)
				if len(formattedProb) != 0 {
					if formattedProb[len(formattedProb)-1] == '.' {
						formattedProb = formattedProb[1:]
					}
				}
				println("\"? (Probably " + formattedProb + "%).")
				if HasAptPackage(cmdName) {
					print("You can install the package \"")
					color.Style{color.FgBlue, color.OpBold}.Print(cmdName)
					println("\" with apt.")
				}
			}
		} else {
			// stop timer
			e := time.Now()

			dur := e.Sub(s)
			// fmt.Println(dur.String())
			color.Style{color.FgGreen, color.OpBold}.Print("\n✓ "+"Done in ", DurationIntFormat(dur)+".\n")
		}
	}
}
