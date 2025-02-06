# GUMP
## Go Unified Managment Package

Gump is a lightweight, flexible configuration management package for Go. It allows you to load, merge, and validate JSON configuration files, as well as easily retrieve configuration values as strings, integers, or booleans. This package is designed to simplify the management of configuration data for your applications.

### Features
- Load JSON Configurations: Easily load configuration files from disk.

- Merge Configurations: Combine multiple configuration files (e.g., base configuration, emergency overrides, and custom overrides) into a single configuration object.

- Validation: Verify that required configuration keys are present.

- Flexible Data Retrieval: Retrieve configuration values as strings, integers, or booleans.

- Dot Notation Access: Use dot notation (e.g., db.host) to access nested configuration values.

## Getting Started

### Prerequisites

- Go 1.18 or higher

```bash
git clone https://github.com/DarioChiappello/gump.git
cd gump
```

### Running the Example

```bash
cd example
go run main.go
```

### Running Tests
```bash
go test ./config
```

### Usage
Below is a short snippet showing how to use Gump in your Go application:

``` golang
package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/DarioChiappello/gump/config"
)

func main() {
    wd, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    // Adjust the path to your configuration files as needed.
    baseConfigPath := filepath.Join(wd, "..", "testdata", "base_config.json")
    emergencyConfigPath := filepath.Join(wd, "..", "testdata", "emergency.json")

    cfg := config.NewConfig()
    if err := cfg.LoadFromJSON(baseConfigPath); err != nil {
        panic(err)
    }

    emergencyCfg := config.NewConfig()
    if err := emergencyCfg.LoadFromJSON(emergencyConfigPath); err != nil {
        panic(err)
    }
    cfg.Merge(emergencyCfg)

    // Validate required keys.
    if err := cfg.Validate([]string{"db.host", "db.port"}); err != nil {
        panic(err)
    }

    // Retrieve values.
    host, _ := cfg.GetString("db.host")
    port, _ := cfg.GetInt("db.port")
    ssl, _ := cfg.GetBool("db.ssl")

    fmt.Println("Database host:", host)
    fmt.Println("Database port:", port)
    fmt.Println("Database SSL:", ssl)
}

```


### Contributing
Contributions are welcome! Feel free to open issues or submit pull requests to improve the package.

### License
This project is licensed under the MIT License.