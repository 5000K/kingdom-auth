package main

import (
	"os"

	"github.com/dimiro1/banner"
)

const (
	templ = `{{ .Title "kingdom-auth" "" 2 }}
`
)

func printBanner() {
	banner.InitString(os.Stdout, true, true, templ)
}
