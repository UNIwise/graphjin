package core

import (
	"fmt"
	"time"

	"github.com/bep/debounce"
	"github.com/dosco/graphjin/core/internal/sdata"
	"github.com/fsnotify/fsnotify"
)

func (g *GraphJin) initDBWatcher() error {
	gj := g.Load().(*graphjin)

	// no schema polling in production unless allowlist is disabled
	if gj.prod && !gj.conf.DisableAllowList {
		return nil
	}

	var ps time.Duration

	switch d := gj.conf.DBSchemaPollDuration; {
	case d < 0:
		return nil
	case d < 5:
		ps = 10 * time.Second
	default:
		ps = d * time.Second
	}

	go func() {
		g.startDBWatcher(ps)
	}()
	return nil
}

func (g *GraphJin) startDBWatcher(ps time.Duration) {
	ticker := time.NewTicker(ps)
	defer ticker.Stop()

	for range ticker.C {
		gj := g.Load().(*graphjin)

		dbinfo, err := sdata.GetDBInfo(
			gj.db,
			gj.dbtype,
			gj.conf.Blocklist)

		if err != nil {
			gj.log.Println(err)
			continue
		}

		if dbinfo.Hash() == gj.dbinfo.Hash() {
			continue
		}

		gj.log.Println("database change detected. reinitializing...")

		if err := g.Reload(); err != nil {
			gj.log.Println(err)
		}
	}
}

func (g *GraphJin) initFSWatcher() error {
	gj := g.Load().(*graphjin)

	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			gj.log.Printf("Error during watcher initialization: %s", err)
			return
		}

		watchPaths := []string{
			gj.conf.ConfigPath,
			fmt.Sprintf("%s/fragments/", gj.conf.ConfigPath),
			fmt.Sprintf("%s/queries/", gj.conf.ConfigPath),
			fmt.Sprintf("%s/scripts/", gj.conf.ConfigPath),
			fmt.Sprintf("%s/keys/", gj.conf.ConfigPath),
		}

		for _, path := range watchPaths {
			err = watcher.Add(path)
			if err != nil {
				gj.log.Printf("Failed adding watch path: %s", path)
				return
			}
		}

		invokeReload := func() {
			gj.log.Println("Reloading config due to one or more file changes...")
			err = g.Reload()
			if err != nil {
				gj.log.Println("Failed to reload config during FS change: ", err)
			}
		}

		debounced := debounce.New(10 * time.Second)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				gj.log.Printf("%s %s\n", event.Name, event.Op)
				debounced(invokeReload)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				gj.log.Println("error:", err)
			}
		}
	}()

	return nil
}
