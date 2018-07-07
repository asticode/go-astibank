package main

import (
	"flag"
	"path/filepath"

	"encoding/json"

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
	d       *data
	p       *pdfParser
	w       *astilectron.Window
)

func main() {
	// Parse flags
	flag.Parse()

	// Init logger
	astilog.FlagInit()

	// Init pdf parser
	p = newPDFParser()

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

	// Init data
	d = newData(a.Paths().DataDirectory())

	// Load data
	if err = d.load(); err != nil {
		err = errors.Wrap(err, "main: loading data failed")
		return
	}

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

	// Init window
	if w, err = a.NewWindow(filepath.Join(a.Paths().DataDirectory(), "resources", "app", "index.html"), &astilectron.WindowOptions{
		BackgroundColor: astilectron.PtrStr("#333"),
		Height:          astilectron.PtrInt(600),
		HideOnClose:     astilectron.PtrBool(true),
		Title:           astilectron.PtrStr(AppName),
		Width:           astilectron.PtrInt(600),
	}); err != nil {
		astilog.Error(errors.Wrap(err, "main: initializing window failed"))
		return
	}

	// Handle messages
	w.OnMessage(bootstrap.HandleMessages(w, handleMessages))

	// Create window
	if err = w.Create(); err != nil {
		astilog.Error(errors.Wrap(err, "main: creating window failed"))
		return
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

	// Save data
	if err = d.save(); err != nil {
		err = errors.Wrap(err, "main: saving data failed")
		return
	}
	return
}

func handleDoubleClicked(_ astilectron.Event) (deleteListener bool) {
	// Show
	if err := w.Show(); err != nil {
		astilog.Error(errors.Wrap(err, "main: showing window failed"))
		return
	}
	return
}

func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "js.import":
		// Decode body
		var paths []string
		if err = json.Unmarshal(m.Payload, &paths); err != nil {
			err = errors.Wrap(err, "main: decoding body failed")
			return
		}

		// Loop through paths
		for _, path := range paths {
			// Only process .pdf files
			if filepath.Ext(path) != ".pdf" {
				astilog.Warnf("main: %s is not a valid pdf file", path)
				continue
			}

			// Parse
			var s pdfStatement
			if s, err = p.parse(path); err != nil {
				err = errors.Wrapf(err, "main: parsing %s failed", path)
				return
			}

			// Add to data
			if err = d.addStatement(s); err != nil {
				err = errors.Wrap(err, "main: adding statement failed")
				return
			}
		}
	case "js.quit":
		// Quit
		if err = a.Quit(); err != nil {
			err = errors.Wrap(err, "main: quitting failed")
			return
		}
	}
	return
}
