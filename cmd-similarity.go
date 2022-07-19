package main

import (
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
)

type CmpMap map[string]float64 // file name - score

func SimilarityJaro(Name string, cmdName string) float64 {
	return strutil.Similarity(Name, cmdName, metrics.NewJaro())
}

func NewCmpMap(fileNames []string, cmdName string) CmpMap {
	cmpm := make(CmpMap)
	for _, name := range fileNames {
		cmpm[name] = SimilarityJaro(name, cmdName)
	}
	return cmpm
}

func (CM CmpMap) Nearest() (string, float64) {
	highest := 0.0
	highestName := ""
	for name, value := range CM {
		if value > highest {
			highestName = name
			highest = value
		}
	}
	return highestName, highest
}
