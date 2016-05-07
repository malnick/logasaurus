package loga

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mgutz/ansi"
)

func highlightQuery(line string, query string) {
	// Split query into multiple parts for regex
	q := strings.Split(query, " ")
	// Match the string
	match, err := regexp.Compile(q[0])
	if err != nil {
		panic(err)
	}

	// Split our line into an ary
	lineAry := strings.Split(line, " ")
	// Iterate the ary, finding the string match
	for i, s := range lineAry {
		if match.MatchString(s) {
			// Color just the string which matches
			hlQuery := ansi.Color(s, "yellow:black")
			// Thren break down into three parts
			lpt1 := lineAry[:i]
			lpt2 := lineAry[i:]
			lpt2 = append(lpt2[:0], lpt2[1:]...)
			// Contatenate back together
			part1 := strings.Join(lpt1, " ")
			part2 := strings.Join(lpt2, " ")
			final := []string{part1, hlQuery, part2}
			finalHl := strings.Join(final, " ")
			// Print the final output
			//log.Info(finalHl)
			fmt.Println(finalHl)
		}
	}
}

func setLogger(verbose bool) {
	if verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("DEBUG Logger")
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
