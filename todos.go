package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Todos struct {
	filename string
	watcher  *fsnotify.Watcher
	ctx      context.Context
}

// NewTodos attempts to create a new Todo list
func NewTodos() (*Todos, error) {
	// Create new Todos instance
	result := &Todos{}
	// Return it
	return result, nil
}
func (a *Todos) startup(ctx context.Context) {
	a.ctx = ctx
	err := a.WailsInit()
	if err != nil {
		runtime.LogInfof(ctx, "init err: %s", err)
	}
}

func (t *Todos) startWatcher() error {
	runtime.LogInfo(t.ctx, "Starting Watcher")
	watcher, err := fsnotify.NewWatcher()
	t.watcher = watcher
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					runtime.LogInfof(t.ctx, "modified file: %s", event.Name)
					runtime.EventsEmit(t.ctx, "filemodified")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				runtime.LogError(t.ctx, err.Error())
			}
		}
	}()

	err = watcher.Add(t.filename)

	if err != nil {
		return err
	}
	return nil
}

func (t *Todos) LoadList() (string, error) {
	runtime.LogInfof(t.ctx, "Loading list from: %s", t.filename)
	bytes, err := ioutil.ReadFile(t.filename)
	if err != nil {
		err = fmt.Errorf("Unable to open list: %s", t.filename)
	}
	// runtime.WindowSetTitle(t.ctx, t.filename)
	runtime.WindowSetTitle(t.ctx, "Todo List")
	runtime.LogInfo(t.ctx, "List content:"+string(bytes))
	return string(bytes), err
}

func (t *Todos) saveListByName(todos string, filename string) error {
	return ioutil.WriteFile(filename, []byte(todos), 0600)
}

func (t *Todos) SaveList(todos string) error {
	if t == nil {
		log.Println("todo instance is nil !!")
		return fmt.Errorf("todo instance is nil: %#v !!", t)
	}
	runtime.LogInfof(t.ctx, "Saving list: %s", todos)
	return t.saveListByName(todos, t.filename)
}

func (t *Todos) setFilename(filename string) error {
	var err error
	// Stop watching the current file and return any error
	err = t.watcher.Remove(t.filename)
	if err != nil {
		return err
	}

	// Set the filename
	t.filename = filename

	// Add the new file to the watcher and return any errors
	err = t.watcher.Add(filename)
	if err != nil {
		return err
	}
	runtime.LogInfo(t.ctx, "Now watching: "+filename)
	runtime.WindowSetTitle(t.ctx, t.filename)
	return nil
}

func (t *Todos) SaveAs(todos string) error {
	filename, err := runtime.SaveFileDialog(t.ctx, runtime.SaveDialogOptions{})
	runtime.LogInfo(t.ctx, "Save As: "+filename)
	err = t.saveListByName(todos, filename)
	if err != nil {
		return err
	}
	return t.setFilename(filename)
}

func (t *Todos) LoadNewList() {
	filename, _ := runtime.OpenFileDialog(t.ctx, runtime.OpenDialogOptions{})
	if len(filename) > 0 {
		t.setFilename(filename)
		runtime.EventsEmit(t.ctx, "filemodified")
	}
}

func (t *Todos) ensureFileExists() {
	// Check status of file
	_, err := os.Stat(t.filename)
	// If it doesn't exist
	if os.IsNotExist(err) {
		// Create it with a blank list
		ioutil.WriteFile(t.filename, []byte("[]"), 0600)
	}
}

func (t *Todos) WailsInit() error {
	runtime.LogDebug(t.ctx, "I'm here")

	// Set the default filename to $HOMEDIR/mylist.json
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	t.filename = path.Join(homedir, "mylist.json")
	runtime.WindowSetTitle(t.ctx, t.filename)
	t.ensureFileExists()
	return t.startWatcher()
}
