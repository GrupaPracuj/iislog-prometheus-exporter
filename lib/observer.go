package lib

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/GrupaPracuj/iislog-prometheus-exporter/logging"

	"github.com/fsnotify/fsnotify"
)

//Observe directory dir. If there is new file send its name by chanel
func newFileCheck(dir string, logger *log.Logger) (newFileAlert <-chan string, err error) {
	logging.Info(logger, fmt.Sprintf("Observing %s directory", dir))
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logging.Error(logger, "Directory watcher fail", err)
		return
	}

	err = watcher.Add(dir)
	if err != nil {
		logging.Error(logger, fmt.Sprintf("Watch directory %s fail", dir), err)
		return
	}

	newFileAlertSender := make(chan string)
	newFileAlert = newFileAlertSender

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op == fsnotify.Create {
					logging.Info(logger, fmt.Sprintf("New file: %s", ev.Name))
					newFileAlertSender <- ev.Name
				}
			case err := <-watcher.Errors:
				logging.Error(logger, "Watcher Error:", err)
			}
		}
	}()

	return
}

func findNewestFile(dir string, logger *log.Logger) (filename string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logging.Error(logger, fmt.Sprintf("Error with reading filenames in directory %s", dir), err)
		return ""
	}

	if len(files) > 0 {
		newestFile := files[0]
		for _, file := range files {
			if file.ModTime().After(newestFile.ModTime()) {
				newestFile = file
			}
		}
		return newestFile.Name()
	}

	return ""
}
