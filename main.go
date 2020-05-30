package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher

func main() {
	// create a new watcher
	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := filepath.Walk(pwd, watchDir); err != nil {
		fmt.Println("ERROR", err)
	}
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				log.Println("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}

				cmd := exec.Command("go", "run", "main.go")
				if err := cmd.Start(); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("reset completed\n")

				cmd = exec.Command("sleep", "4")
				cmd.Run()
				break

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error =========")
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}

func watchDir(path string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}
	return nil
}

// func execCmd(cmd string, args ...string) {
// 	command := exec.Command(cmd, args...)
// 	var out bytes.Buffer
// 	var stderr bytes.Buffer
// 	command.Stdout = &out
// 	command.Stderr = &stderr
// 	err := command.Run()
// 	if err != nil {
// 		errstring := fmt.Sprintf(fmt.Sprint(err) + ": " + stderr.String())
// 		io.WriteString(nil, errstring)
// 	}
// 	io.WriteString(nil, out.String())
// 	fmt.Println(out.String())
// }
