package main

import (
	"github.com/DarioChiappello/gump/config"
)

func main() {
	cfg := config.NewConfig()

	// Load base configuration
	err := cfg.LoadFromJSON("./base_config.json")
	if err != nil {
		panic("Failed to load base config: " + err.Error())
	}

	// Load environment overrides
	// err = cfg.LoadFromEnvironment("APP")
	// if err != nil {
	// 	panic("Failed to load environment config: " + err.Error())
	// }

	// // Merge with emergency defaults
	// emergencyCfg := config.NewConfig()
	// err = emergencyCfg.LoadFromJSON("emergency.json")
	// if err != nil {
	// 	panic("Failed to load emergency config: " + err.Error())
	// }
	// cfg.Merge(emergencyCfg)

	// Validate required keys
	if err := cfg.Validate([]string{"db.host", "db.port"}); err != nil {
		panic("Configuration validation failed: " + err.Error())
	}

	// Access values
	host, err := cfg.GetString("db.host")
	if err != nil {
		panic("Failed to get db.host: " + err.Error())
	}

	port, err := cfg.GetInt("db.port")
	if err != nil {
		panic("Failed to get db.port: " + err.Error())
	}

	ssl, err := cfg.GetBool("db.ssl")
	if err != nil {
		panic("Failed to get db.ssl: " + err.Error())
	}

	// Use the variables
	println("Database host:", host)
	println("Database port:", port)
	println("Database SSL:", ssl)
}
