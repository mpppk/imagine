module github.com/mpppk/imagine

go 1.14

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/blang/semver/v4 v4.0.0
	github.com/comail/colog v0.0.0-20160416085026-fba8e7b1f46c
	github.com/gen2brain/dlgs v0.0.0-20200211102745-b9c2664df42f
	github.com/google/wire v0.4.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mpppk/lorca-fsa v0.0.0-00010101000000-000000000000
	github.com/rhysd/go-github-selfupdate v1.2.2
	github.com/spf13/afero v1.3.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	go.etcd.io/bbolt v1.3.5
	golang.org/x/crypto v0.0.0-20200128174031-69ecbb4d6d5d // indirect
)

replace github.com/zserge/lorca => ../lorca-fsa/lorca

replace github.com/mpppk/lorca-fsa => ../lorca-fsa/lorca-fsa
