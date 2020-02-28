package filewatcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
)

type FileWatcher struct {
	Handler     WatchHandler
	ExcludeDirs []string

	watcher *fsnotify.Watcher
}

func New() *FileWatcher {
	return &FileWatcher{}
}

// Watch will do an initial scan to find sub-directories to watch, and will then
// start listening for file updates.  This is to be started in a Go routine.
func (fw *FileWatcher) Watch() error {
	var err error
	fw.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := fw.subscribeToDirRecursively("."); err != nil {
		return err
	}

	log.Println("Scanning for files")
	for {
		select {
		case event := <-fw.watcher.Events:

			// TODO: If it's a new directory, subscribe

			if event.Op&fsnotify.Write == fsnotify.Write {
				if filepath.Ext(event.Name) == ".go" {
					log.Println("modified file:", event.Name, " Rebuilding")
					fw.Handler.OnFileModified(event.Name)
				}
			}
		case err := <-fw.watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func (fw *FileWatcher) subscribeToDirRecursively(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return fmt.Errorf("cannot open '%v': %v", dir, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("cannot stat '%v': %v", dir, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("cannot subscribe to '%v': not a directory", dir)
	}

	fInfo, err := f.Readdir(0)
	if err != nil {
		return fmt.Errorf("cannot read entries of dir '%v': %v", dir, err)
	}

	fw.subscribeToDir(dir)

fileIter:
	for _, f := range fInfo {
		// TODO: Ignore list
		if f.IsDir() {
			for _, excludeDir := range fw.ExcludeDirs {
				if match, _ := filepath.Match(excludeDir, f.Name()); match {
					continue fileIter
				}
			}
			if err := fw.subscribeToDirRecursively(filepath.Join(dir, f.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

func (fw *FileWatcher) subscribeToDir(dir string) {
	if err := fw.watcher.Add(dir); err != nil {
		log.Printf("cannot listen to '%v': %v", dir, err)
	}

	log.Printf("listening to '%v'", dir)
}
