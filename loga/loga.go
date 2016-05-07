package loga

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/malnick/logasaurus/config"
)

func Start() {
	fmt.Println(`                        .       .                             `)
	fmt.Println(`                       / '.   .' \                            `)
	fmt.Println(`               .---.  <    > <    >  .---.                    `)
	fmt.Println(`               |    \  \ - ~ ~ - /  /    |                    `)
	fmt.Println(`               ~-..-~             ~-..-~                     `)
	fmt.Println(`            \~~~\.'                    './~~~/                `)
	fmt.Println(`  .-~~^-.    \__/                        \__/                 `)
	fmt.Println(`.'  O    \     /               /       \  \                   `)
	fmt.Println(`(_____'    \._.'              |         }  \/~~~/             `)
	fmt.Println(`  ----.         /       }     |        /    \__/              `)
	fmt.Println(`      \-.      |       /      |       /      \.,~~|           `)
	fmt.Println(`          ~-.__|      /_ - ~ ^|      /- _     \..-'   f: f:   `)
	fmt.Println(`               |     /        |     /     ~-.     -. _||_||_  `)
	fmt.Println(`               |_____|        |_____|         ~ - . _ _ _ _ _>`)
	fmt.Println(`██╗      ██████╗  ██████╗  █████╗ ███████╗ █████╗ ██╗   ██╗██████╗ ██╗   ██╗███████╗`)
	fmt.Println(`██║     ██╔═══██╗██╔════╝ ██╔══██╗██╔════╝██╔══██╗██║   ██║██╔══██╗██║   ██║██╔════╝`)
	fmt.Println(`██║     ██║   ██║██║  ███╗███████║███████╗███████║██║   ██║██████╔╝██║   ██║███████╗`)
	fmt.Println(`██║     ██║   ██║██║   ██║██╔══██║╚════██║██╔══██║██║   ██║██╔══██╗██║   ██║╚════██║`)
	fmt.Println(`███████╗╚██████╔╝╚██████╔╝██║  ██║███████║██║  ██║╚██████╔╝██║  ██║╚██████╔╝███████║`)
	fmt.Println(`╚══════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝`)
	fmt.Println()
	config := config.ParseArgsReturnConfig()
	setLogger(config.LogVerbose)
	query, err := config.GetDefinedQuery()
	BasicCheckOrExit(err)
	log.WithFields(log.Fields{
		"Query":        query,
		"Elastic Host": config.ElasticsearchURL,
		"Elastic Port": config.ElasticsearchPort,
	}).Info("Elastic Runner")
	// Roll into the query loop
	elasticRunner(query, config)
}
