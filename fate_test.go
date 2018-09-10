package main

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"os"
	"strings"
	"testing"
	"time"
)

type MockVaultClient struct {
	count int
}

func (m *MockVaultClient) Write(path string, data interface{}) error {
	m.count++
	return nil
}

type MockFileInfo struct {
	isDir bool
	name  string
}

func (m *MockFileInfo) Name() string {
	return m.name
}

func (m *MockFileInfo) Size() int64 {
	panic("not implemented")
}

func (m *MockFileInfo) Mode() os.FileMode {
	panic("not implemented")
}

func (m *MockFileInfo) ModTime() time.Time {
	panic("not implemented")
}

func (m *MockFileInfo) IsDir() bool {
	return m.isDir
}

func (m *MockFileInfo) Sys() interface{} {
	panic("not implemented")
}

type FakeWatcher struct {
	event chan fsnotify.Event
	error chan error
}

func (w *FakeWatcher) Close() error {
	close(w.event)
	return nil
}
func (w *FakeWatcher) Add(name string) error       { return nil }
func (w *FakeWatcher) Events() chan fsnotify.Event { return w.event }
func (w *FakeWatcher) Errors() chan error          { return w.error }

var sampleYml = `
good:
  boy:
    - test
    - JavaScript:
      front-end:
        react: 1000
        redux: 69.22
      back-end:
        node: .12
        next: true
        Java: false
      Java.fake:
`

func TestScanAndWrite(t *testing.T) {
	client, _ := mockStuff()

	scanDir("fake-dir")

	if client.count != 1 {
		t.Fatalf("Expected one hit to vault, but got %d", client.count)
	}
}

func mockStuff() (*MockVaultClient, <-chan os.Signal) {
	readDir = func(dirname string) ([]os.FileInfo, error) {
		return []os.FileInfo{
			&MockFileInfo{true, "dir1"},
			&MockFileInfo{false, "empty-file"},
			&MockFileInfo{false, "good.yml"},
		}, nil
	}
	readFile = func(filename string) ([]byte, error) {
		if strings.HasSuffix(filename, "good.yml") {
			return []byte(sampleYml), nil
		}
		return nil, errors.New("failed")
	}

	var exitChn = make(chan os.Signal)

	signalNotify = func(c chan<- os.Signal, sig ...os.Signal) {
		go func() {
			time.Sleep(time.Second * 1)
			c <- sig[0]
			exitChn <- sig[0]
		}()
	}

	client := new(MockVaultClient)
	fate.client = client
	return client, exitChn
}

func TestWatchDir(t *testing.T) {
	client, exited := mockStuff()
	delayReadDuration = 0
	watcher = &FakeWatcher{
		event: make(chan fsnotify.Event),
	}

	go watchDir("", watcher)
	events := watcher.Events()
	events <- fsnotify.Event{Name: "good.yml", Op: fsnotify.Create}
	events <- fsnotify.Event{Name: "good.yml", Op: fsnotify.Write}
	events <- fsnotify.Event{Name: "good.yml", Op: fsnotify.Chmod}
	events <- fsnotify.Event{Name: "good.yml", Op: fsnotify.Remove}

	<-exited

	if client.count != 2 {
		t.Fatalf("Expected 2 hits to vault, but got %d", client.count)
	}
}
