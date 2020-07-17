module github.com/mpppk/imagine

go 1.14

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/blang/semver/v4 v4.0.0
	github.com/comail/colog v0.0.0-20160416085026-fba8e7b1f46c
	github.com/gen2brain/dlgs v0.0.0-20200211102745-b9c2664df42f
	github.com/go-playground/validator/v10 v10.3.0
	github.com/google/wire v0.4.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/gotk3/gotk3 v0.4.0 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mpppk/lorca-fsa v0.0.0-00010101000000-000000000000
	github.com/rhysd/go-github-selfupdate v1.2.2
	github.com/spf13/afero v1.3.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/sqweek/dialog v0.0.0-20200601143742-43ea34326190
	github.com/valyala/fasttemplate v1.1.0 // indirect
	go.etcd.io/bbolt v1.3.5
	golang.org/x/crypto v0.0.0-20200128174031-69ecbb4d6d5d // indirect
)

replace github.com/zserge/lorca => ../lorca-fsa/lorca

replace github.com/mpppk/lorca-fsa => ../lorca-fsa/lorca-fsa
