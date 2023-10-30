# confik

Confik parses environment files and variables and loads them into a struct.

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
