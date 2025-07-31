package main

import "github.com/alecthomas/kong"
var Cli struct {
	Number int `name:"number" short:"n" default:"1" help:"Number of UUIDs to generate."`
	NoHyphons bool `name:"no-hyphens" short:"H" help:"Do not include hyphens in the UUID."`
}

func main()  {
	kong.Parse(&Cli,
		kong.Name("go-uuidv7"),
		kong.Description("Generate UUIDv7 (draft)"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}),
		kong.Vars{
			"version": "0.1.0",
			"commit":  "unknown",
			"date":    "unknown",
		},
	)	
}