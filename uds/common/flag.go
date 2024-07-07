package common

import (
	"emperror.dev/emperror"
	"flag"
	"os"
)

var socketPathPtr = flag.String("p", "", "Unix Domain Socket path")

var SocketPath string

func Init(genNew bool) {
	flag.Parse()
	SocketPath = *socketPathPtr

	if SocketPath == "" && genNew {
		f, err := os.CreateTemp("", "vinci_uds_")
		emperror.Panic(err)

		SocketPath = f.Name()
		emperror.Panic(f.Close())
	}
	_ = os.Remove(SocketPath)
}
