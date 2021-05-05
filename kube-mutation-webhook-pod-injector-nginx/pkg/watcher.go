/*
 * @Author: calmwu
 * @Date: 2021-05-04 19:41:29
 * @Last Modified by: calmwu
 * @Last Modified time: 2021-05-04 20:34:22
 */

package pkg

import (
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type OnUpdate func(configFile string) error

// Watcher watches for config updates
type Watcher interface {
	//
	SetUpdateNotify(OnUpdate)

	// Run starts the Watcher, Must call this after SetNotify
	Run(<-chan struct{})
}

type fileWatcher struct {
	watcher       *fsnotify.Watcher
	updateHandler OnUpdate
	configFile    string
}

const (
	watchDebounceDelay = 100 * time.Millisecond
)

// NewFileWatcher create a Watcher for local config file
func NewFileWatcher(configFile string) (Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		err = errors.Wrap(err, "fsnotify New Watcher failed.")
		glog.Error(err)
		return nil, err
	}

	// 配置文件实际是个symlink，只有watch父目录。
	watchDir, _ := filepath.Split(configFile)
	if err := watcher.Add(watchDir); err != nil {
		err = errors.Wrapf(err, "Could not watch %s", watchDir)
		glog.Error(err)
		return nil, err
	}

	return &fileWatcher{
		watcher:    watcher,
		configFile: configFile,
	}, nil
}

func (fw *fileWatcher) SetUpdateNotify(handler OnUpdate) {
	fw.updateHandler = handler
}

func (fw *fileWatcher) Run(stopCh <-chan struct{}) {
	defer fw.watcher.Close()
	var timeC <-chan time.Time

	for {
		select {
		case <-timeC:
			timeC = nil
			if fw.updateHandler != nil {
				fw.updateHandler(fw.configFile)
			}
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}
			glog.Infof("Inject watch update: %+v", event)
			if ((event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Create == fsnotify.Create)) && timeC == nil {
				timeC = time.After(watchDebounceDelay)
				glog.Infof("modified file: %s", event.Name)
			}
		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			glog.Errorf("Inject watch error: %v", err)
		case <-stopCh:
			return
		}
	}
}
