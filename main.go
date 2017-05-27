package main

import (
	"flag"
	"path/filepath"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilectron-bundler"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Vars
var (
	a       *astilectron.Astilectron
	AppName string
	w       *astilectron.Window
)

func main() {
	// Parse flags
	flag.Parse()

	// Init logger
	astilog.FlagInit()

	// Run app
	if err := runApp(); err != nil {
		astilog.Fatal(errors.Wrap(err, "main: running app failed"))
	}
}

func runApp() (err error) {
	// Create astilectron
	if a, err = astilectron.New(astilectron.Options{
		AppName:            AppName,
		AppIconDarwinPath:  "resources/icon.icns",
		AppIconDefaultPath: "resources/icon.png",
	}); err != nil {
		return errors.Wrap(err, "main: creating new astilectron failed")
	}
	defer a.Close()

	// Handle signals
	a.HandleSignals()

	// Set provisioner
	a.SetProvisioner(astibundler.NewProvisioner(Asset))

	// Restore resource
	if err = bootstrap.RestoreResources(a, Asset, AssetDir, RestoreAssets, "resources"); err != nil {
		return errors.Wrap(err, "main: restoring resources failed")
	}

	// Start
	if err = a.Start(); err != nil {
		return errors.Wrap(err, "main: starting astilectron failed")
	}

	// Init tray
	t := a.NewTray(&astilectron.TrayOptions{
		Image:   astilectron.PtrStr(filepath.Join(a.Paths().DataDirectory(), "resources", "tray.png")),
		Tooltip: astilectron.PtrStr(AppName),
	})

	// Handle double click
	t.On(astilectron.EventNameTrayEventDoubleClicked, handleDoubleClicked)

	// Create tray
	if err = t.Create(); err != nil {
		return errors.Wrap(err, "main: creating tray failed")
	}

	// Init tray menu
	tm := t.NewMenu([]*astilectron.MenuItemOptions{
		{Label: astilectron.PtrStr("Quit"), Role: astilectron.MenuItemRoleQuit},
	})

	// Create tray menu
	if err = tm.Create(); err != nil {
		return errors.Wrap(err, "creating tray menu failed")
	}

	// Blocking pattern
	a.Wait()
	return
}

func handleDoubleClicked(_ astilectron.Event) (deleteListener bool) {
	// Window already exists
	var err error
	if w != nil {
		// Show
		if err = w.Show(); err != nil {
			astilog.Error(errors.Wrap(err, "main: showing window failed"))
			return
		}
		return
	}

	// Init window
	if w, err = a.NewWindow(filepath.Join(a.Paths().DataDirectory(), "resources", "app", "index.html"), &astilectron.WindowOptions{
		BackgroundColor: astilectron.PtrStr("#333"),
		HideOnClose:     astilectron.PtrBool(true),
		Show:            astilectron.PtrBool(false),
		Title:           astilectron.PtrStr(AppName),
	}); err != nil {
		astilog.Error(errors.Wrap(err, "main: initializing window failed"))
		return
	}

	// Maximize on init
	w.On(astilectron.EventNameWindowEventDidFinishLoad, func(_ astilectron.Event) (deleteListener bool) {
		if err = w.Maximize(); err != nil {
			astilog.Error(errors.Wrap(err, "main: maximizing window failed"))
			return
		}
		if err = w.Show(); err != nil {
			astilog.Error(errors.Wrap(err, "main: showing window failed"))
			return
		}
		return true
	})

	// Create window
	if err = w.Create(); err != nil {
		astilog.Error(errors.Wrap(err, "main: creating window failed"))
		return
	}
	return
}
