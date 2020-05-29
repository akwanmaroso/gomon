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
	fmt.Println("helo")
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
					execute()
					log.Println("modified file: ", event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
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

func execute() {
	out, err := exec.Command("go run main.go").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Println("Command success executed")
	output := string(out[:])
	fmt.Println(output)
}
