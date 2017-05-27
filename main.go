package main

import (
	"flag"

	"os"
	"path/filepath"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron/bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Vars
var data *Data

//go:generate go-bindata -pkg $GOPACKAGE -o resources.go resources/...
func main() {
	// Init
	flag.Parse()
	astilog.SetLogger(astilog.New(astilog.FlagConfig()))

	// Fetch executable path
	var p string
	var err error
	if p, err = os.Executable(); err != nil {
		astilog.Fatal(errors.Wrap(err, "fetching executable path failed"))
	}
	p = filepath.Dir(p)

	// Import data
	if data, err = NewData(p); err != nil {
		astilog.Fatal(errors.Wrap(err, "importing data failed"))
	}
	defer data.Close()

	// Run bootstrap
	if err = bootstrap.Run(bootstrap.Options{
		AdaptRouter: adaptRouter,
		AstilectronOptions: astilectron.Options{
			AppName: "Astibank",
		},
		Homepage: "/templates/index",
		// RestoreAssets: RestoreAssets,
		WindowOptions: &astilectron.WindowOptions{
			BackgroundColor: astilectron.PtrStr("#333"),
			Center:          astilectron.PtrBool(true),
			Height:          astilectron.PtrInt(600),
			Width:           astilectron.PtrInt(600),
		},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}
