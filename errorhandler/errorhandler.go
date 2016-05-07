package errorhandler

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func BasicCheckOrExit(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
