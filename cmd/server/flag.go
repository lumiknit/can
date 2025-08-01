package main

import (
	"flag"
	"log"

	"github.com/lumiknit/can/internal/config"
)

func HandleFlags() *config.Config {
	// Define flag variables
	var (
		configFile  = flag.String("config", "", "Path to JSON configuration file")
		host        = flag.String("host", "", "Host to bind the server to")
		port        = flag.Int("port", 0, "Port to bind the server to")
		basePath    = flag.String("base", "", "Base path for the application")
		releaseMode = flag.Bool("release-mode", false, "Enable release mode")
	)

	// Add short flag aliases
	flag.StringVar(configFile, "C", "", "Path to JSON configuration file (short)")
	flag.StringVar(host, "h", "", "Host to bind the server to (short)")
	flag.IntVar(port, "p", 0, "Port to bind the server to (short)")
	flag.StringVar(basePath, "b", "", "Base path for the application (short)")
	flag.BoolVar(releaseMode, "R", false, "Enable release mode (short)")

	flag.Parse()

	// Start with base configuration
	var cfg *config.Config

	if *configFile != "" {
		// Load from config file if specified
		var err error
		cfg, err = config.LoadFromFile(*configFile)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	} else {
		// Use default configuration
		cfg = config.Default()
	}

	// Override with command line flags (only if they were explicitly set)
	if *host != "" {
		cfg.Host = *host
	}
	if *port != 0 {
		cfg.Port = *port
	}
	if *basePath != "" {
		cfg.BasePath = *basePath
	}
	if *releaseMode {
		cfg.ReleaseMode = *releaseMode
	}

	return cfg
}
