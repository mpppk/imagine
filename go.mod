module github.com/mpppk/imagine

go 1.14

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/blang/semver/v4 v4.0.0
	github.com/comail/colog v0.0.0-20160416085026-fba8e7b1f46c
	github.com/gen2brain/dlgs v0.0.0-20201118155338-03fe7f81ad25
	github.com/golang/mock v1.4.4
	github.com/google/wire v0.4.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.3.3
	github.com/mpppk/lorca-fsa/lorca-fsa v0.0.0-20200916170540-145bd67d1a8e
	github.com/rhysd/go-github-selfupdate v1.2.2
	github.com/spf13/afero v1.4.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	go.etcd.io/bbolt v1.3.5
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
)

// replace github.com/zserge/lorca => ../lorca-fsa/lorca

//replace github.com/mpppk/lorca-fsa/lorca-fsa => ../lorca-fsa/lorca-fsa
