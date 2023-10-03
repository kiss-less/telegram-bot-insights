package handlers

import (
	"flag"
	"log"
)

type Arguments struct {
	Parse          bool
	Debug          bool
	ParseDirectory string
	DateFormat     string
	Regex          string
	BotsApiEnvVar  string
}

func ParseArgs(defaultArgs Arguments) Arguments {
	flag.BoolVar(&defaultArgs.Parse, "parse", defaultArgs.Parse, "Parse Json files")
	flag.BoolVar(&defaultArgs.Debug, "debug", defaultArgs.Debug, "Verbose output")
	parseDirPtr := flag.String("dir", defaultArgs.ParseDirectory, "Should be passed if parse is true")
	customDateFormatPtr := flag.String("custom-date-fmt", defaultArgs.DateFormat, "Custom date format")
	customRegexPtr := flag.String("custom-regex", defaultArgs.Regex, "Custom regex to parse chat_id")
	customEnvVarPtr := flag.String("custom-env-var", defaultArgs.BotsApiEnvVar, "Custom Env Var")
	flag.Parse()

	parseDir := *parseDirPtr
	customDateFormat := *customDateFormatPtr
	customRegex := *customRegexPtr
	customEnvVar := *customEnvVarPtr

	if defaultArgs.Parse && parseDir == "" {
		log.Fatal("Please provide --dir argument when setting --parse to true")
	} else if defaultArgs.Parse && parseDir != "" {
		defaultArgs.ParseDirectory = parseDir
	}

	if customDateFormat != "" {
		defaultArgs.DateFormat = customDateFormat
	}

	if customRegex != "" {
		defaultArgs.Regex = customRegex
	}

	if customEnvVar != "" {
		defaultArgs.BotsApiEnvVar = customEnvVar
	}

	return defaultArgs
}
