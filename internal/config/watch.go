package config

import (
	"context"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func Watch(ctx context.Context, path string, callBack func(*Config)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	dir := filepath.Dir(path)
	file := filepath.Base(path)

	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case event := <-watcher.Events:
			if filepath.Base(event.Name) != file {
				continue
			}

			if event.Op&(fsnotify.Create|fsnotify.Write|fsnotify.Create) != 0 {
				newCfg, err := LoadConfig(path)
				if err != nil {
					return err
				}
				callBack(newCfg)
			}
		case err := <-watcher.Errors:
			return err
		}
	}
}
