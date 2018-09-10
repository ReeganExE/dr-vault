package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

type Fate struct {
	client VaultClientInterface
}

type Watcher interface {
	Close() error
	Add(name string) error
	Events() chan fsnotify.Event
	Errors() chan error
}

type FsWatcher struct {
	watcher *fsnotify.Watcher
}

func (w *FsWatcher) Close() error                { return w.watcher.Close() }
func (w *FsWatcher) Add(name string) error       { return w.watcher.Add(name) }
func (w *FsWatcher) Events() chan fsnotify.Event { return w.watcher.Events }
func (w *FsWatcher) Errors() chan error          { return w.watcher.Errors }

func (v *Fate) read(file string) {
	v.delayRead(file, 0)
}

func (v *Fate) delayRead(file string, d time.Duration) {
	if !strings.HasSuffix(file, ".yml") {
		log.Printf("Not a yaml file. Ignored %s", file)
		return
	}

	if d != 0 {
		time.Sleep(d)
	}

	path := filepath.Base(file)

	dest := make(map[string]interface{})

	str, _ := readFile(file)

	if err := yaml.Unmarshal(str, &dest); err != nil {
		log.Printf("invalid YAML %s, %v", file, err)
		return
	}

	data := Flatten(dest)
	path = strings.TrimSuffix(path, filepath.Ext(path))

	log.Printf("Writing path %s", path)

	if e := v.client.Write("secret/"+path, data); e != nil {
		log.Fatal(e)
	}
}

func scanDir(dir string) {
	log.Print("Scanning ...")
	files, err := readDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		log.Println(f.Name())
		fate.read(fmt.Sprintf("%s/%s", dir, f.Name()))
	}
	log.Print("Scan completed.")
}

func watchDir(dir string, watcher Watcher) {
	defer watcher.Close()

	done := make(chan os.Signal)
	signalNotify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events():
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Modified file:", event.Name)
					go fate.delayRead(event.Name, delayReadDuration)
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("Created file:", event.Name)
					go fate.delayRead(event.Name, delayReadDuration)
				}

			case err, ok := <-watcher.Errors():
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err := watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Watching dir: %s", dir)
	<-done
	log.Println("Exiting")
}
