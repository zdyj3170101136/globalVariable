package main

import (
	"github.com/zdyj3170101136/globalVariable"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(globalVariable.Analyzer) }
