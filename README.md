# confik
[![Build and Test](https://github.com/42z-io/confik/actions/workflows/build_test.yml/badge.svg)](https://github.com/42z-io/confik/actions/workflows/build_test.yml) [![Coverage Status](https://coveralls.io/repos/github/42z-io/confik/badge.svg?branch=main)](https://coveralls.io/github/42z-io/confik?branch=main) [![GitHub Tag](https://img.shields.io/github/tag/42z-io/confik?include_prereleases=&sort=semver&color=blue)](https://github.com/42z-io/confik/releases/)
[![License](https://img.shields.io/badge/License-MIT-blue)](https://github.com/42z-io/confik/blob/main/LICENSE) [![Docs](https://img.shields.io/badge/API-docs?label=docs&color=blue&link=https%3A%2F%2Fpkg.go.dev%2Fgithub.com%2F42z-io%2Fconfik)](https://pkg.go.dev/github.com/42z-io/confik)

![Logo](logo.png)


`Confik` parses environment files and variables and loads them into a struct.

## Usage

```
go get github.com/42z-io/confik
```

```go
import (
    "os"
    "fmt"
    "github.com/42z-io/confik"
)

type ExampleConfig struct {
    Name   string
    Age    uint8 `env:"AGE,optional"`
    Height float32
}

func init() {
    os.Setenv("NAME", "Bob")
    os.Setenv("AGE", "20")
    os.Setenv("HEIGHT", "5.3")

    cfg, _ := confik.LoadFromEnv(Config[ExampleConfig]{
        UseEnvFile: false,
    })

    fmt.Println(cfg.Name)
    fmt.Println(cfg.Age)
    fmt.Println(cfg.Height)
    // Output: Bob
    // 20
    // 5.3
}
```
