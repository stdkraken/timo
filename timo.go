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
		return strings.Split(PATH, ":")
	}
	return []string{}
}()

var CURRENT_COMMAND *exec.Cmd

func Sep() string { // get system default seperator
	switch runtime.GOOS {
	case "windows":
		return "\\"
	default:
		return "/"
	}
}

func FloatStrLimit(s string, l int) string {
	start := false
	a := ""
	b := ""
	d := ""
	for _, c := range s {
		// check if c is non numeric
		isNumeric := func(b rune) bool {
			return b == '0' ||
				b == '1' ||
				b == '2' ||
				b == '3' ||
				b == '4' ||
				b == '5' ||
				b == '6' ||
				b == '7' ||
				b == '8' ||
				b == '9'
		}
		if c == '.' {
			start = true
			d += string(c)
			continue
		}
		if !isNumeric(c) {
			b += string(c)
			continue
		}
		if start && len(b)+len(a) < l {
			b += string(c)
		} else {
			if len(b)+len(a) < l {
				a += string(c)
			}
		}
	}
	return a + d + b
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
			/*
				if err := CURRENT_COMMAND.Process.Kill(); err != nil {
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

func main() {
	go CommandCleaner()
	log.SetOutput(&NullOutput{})
	tm.Clear() // Clear current screen
	tm.MoveCursor(1, 1)
	log.Println(SPLITTED_PATH)
	GenMOTD()
	for {
		color.Style{color.FgLightMagenta, color.OpBold}.Print("⋙ ")

		inputStyle := color.Style{color.FgLightCyan}
		PrintColorMod(inputStyle, " ")

		buffer := make([]byte, BUFFER_SIZE)
		// read from os.Stdin, until err != nil or n != 0 or
		n, err := os.Stdin.Read(buffer)
		if err != nil {
			log.Fatalln(err.Error())
		}
		buffer = BufferCut(buffer, n-1) // remove last byte, because this is 0xa, and not very helpful
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

		log.Println("Name:", cmdName)
		log.Println("Args:", cmdArgs)
		log.Println("Raw:", cmdRaw)

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
					CURRENT_COMMAND = cmd
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
			cmpm := NewCmpMap(names, cmdName)
			nearCmdName, prob := cmpm.Nearest()
			if nearCmdName != "" {
				print("Did you mean \"")
				color.Style{color.FgGreen, color.OpBold}.Print(nearCmdName)
				formattedProb := FloatStrLimit(fmt.Sprint(prob*100), 3)
				if len(formattedProb) != 0 {
					if formattedProb[len(formattedProb)-1] == '.' {
						formattedProb = formattedProb[1:]
					}
				}
				println("\"? (Probably " + formattedProb + "%).")
			}
		} else {
			// stop timer
			e := time.Now()

			dur := e.Sub(s)
			// fmt.Println(dur.String())
			color.Style{color.FgGreen, color.OpBold}.Print("\n✓ ")
			fmt.Println("Done in", FloatStrLimit(dur.String(), 6)+".")
		}
	}
}
