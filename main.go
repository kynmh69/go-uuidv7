package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/kynmh69/go-uuidv7/utils"
)

var Cli struct {
	Number    int  `name:"number" short:"n" default:"1" help:"Number of UUIDs to generate."`
	NoHyphons bool `name:"no-hyphens" short:"H" default:"false" help:"Do not include hyphens in the UUID."`
}

func main() {
	kong.Parse(&Cli,
		kong.Name("go-uuidv7"),
		kong.Description("Generate UUIDv7 (draft)"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
		kong.Vars{
			"version": "1.0.0",
			"date":    "2025-07-31",
		},
	)

	// 負の数のバリデーション
	if Cli.Number <= 0 {
		fmt.Fprintf(os.Stderr, "Error: Number must be a positive integer, got %d\n", Cli.Number)
		os.Exit(1)
	}

	uuidList := utils.GenerateMultipleUUIDs(Cli.Number, Cli.NoHyphons)
	utils.PrintUUIDs(uuidList)
}
