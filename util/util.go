// Package util provides some utilities
package util

import (
	"log"
	"os"

	"github.com/comail/colog"
)

var clg *colog.CoLog

// InitializeLog initialize log settings
func InitializeLog(verbose bool) {
	colog.Register()
	colog.SetOutput(os.Stderr)
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LInfo)

	if verbose {
		colog.SetMinLevel(colog.LDebug)
	}

	clg = colog.NewCoLog(os.Stderr, "", 0)
	clg.SetOutput(os.Stderr)
	clg.SetDefaultLevel(colog.LDebug)
	clg.SetMinLevel(colog.LInfo)
	if verbose {
		clg.SetMinLevel(colog.LDebug)
	}
}

func GetLogger() *log.Logger {
	return clg.NewLogger()
}
