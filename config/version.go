package config

import (
	"flag"
)

var (
	BuildVersion string
	BuildTime    string
	BuildName    string
	CommitID     string
	ShowVer      bool
)

func init() {
	flag.BoolVar(&ShowVer, "version", false, "show version")
}
