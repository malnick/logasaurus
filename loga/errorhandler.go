package loga

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

func CheckElasticResponse(response *ESResponse) {
	if response.Status != 200 {
		log.Errorf("Elastic search returned an error handling request, run in debug mode (-v) to see complete response. HTTP code: %d", response.Status)
		os.Exit(1)
	}
}
