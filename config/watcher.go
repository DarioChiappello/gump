package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// ConfigWatcher observe config files changes
type ConfigWatcher struct {
	config    *Config
	filePaths []string
	watcher   *fsnotify.Watcher
	interval  time.Duration
	callbacks []func(*Config)
	stop      chan struct{}
}

// NewConfigWatcher create new config observer
func NewConfigWatcher(cfg *Config, reloadInterval time.Duration, files ...string) (*ConfigWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Add folders to watcher
	dirs := make(map[string]bool)
	for _, file := range files {
		dir := filepath.Dir(file)
		if _, exists := dirs[dir]; !exists {
			if err := watcher.Add(dir); err == nil {
				dirs[dir] = true
			}
		}
	}

	return &ConfigWatcher{
		config:    cfg,
		filePaths: files,
		watcher:   watcher,
		interval:  reloadInterval,
		stop:      make(chan struct{}),
	}, nil
}

// OnReload register callback for changes
func (w *ConfigWatcher) OnReload(callback func(*Config)) {
	w.callbacks = append(w.callbacks, callback)
}

// Start observer
func (w *ConfigWatcher) Start() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				for _, file := range w.filePaths {
					if filepath.Clean(event.Name) == filepath.Clean(file) {
						w.reloadConfig()
						break
					}
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Config watcher error: %v", err)

		case <-ticker.C:
			// Verify changes
			for _, file := range w.filePaths {
				if w.fileChanged(file) {
					w.reloadConfig()
					break
				}
			}

		case <-w.stop:
			w.watcher.Close()
			return
		}
	}
}

// Stop watcher
func (w *ConfigWatcher) Stop() {
	close(w.stop)
}

func (w *ConfigWatcher) fileChanged(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return info.ModTime().After(w.config.LastModified)
}

func (w *ConfigWatcher) reloadConfig() {
	newConfig := NewConfig()
	success := true

	for _, file := range w.filePaths {
		if err := newConfig.LoadFromJSON(file); err != nil {
			log.Printf("Error reloading config: %v", err)
			success = false
			// Continue trying with other files
		}
	}

	if !success {
		return // No apply changes or invoke callbacks
	}

	// Update main config
	w.config.Merge(newConfig)
	w.config.LastModified = time.Now()

	// Invoke callbacks
	for _, callback := range w.callbacks {
		callback(w.config)
	}
}
