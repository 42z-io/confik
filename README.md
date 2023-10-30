# confik

Confik parses environment files and variables and loads them into a struct.

## Usage

```go
package main

import (
    "github.com/42z-io/confik"
)


type MyEnv struct {
    Name        string                                          // NAME=Jean-Luc Picard
    Age         uint8                                           // AGE=25
    Height      float32                                         // HEIGHT=6.2
    Aliases     []string `env:"ALIAS_NAMES"`                    // ALIAS_NAMES=Captain,Picard,Jean-Luc
    Scores      []int                                           // SCORES=10,50,99
    Birthday    time.Time                                       // BIRTHDAY=2008-10-15T00:42:42Z
    Website     url.URL                                         // WEBSITE=https://42z.io
    TimeLived   time.Duration                                   // TIME_LIVED=25yr
    EmailOptin  bool `env:"EMAIL_OPTIN,optional,default=no"`    // EMAIL_OPTIN=true
    HomeFolder  string `env:"HOME_FOLDER,validate=dir"`         // HOME_FOLDER=/home/jean
    SecretFile  string `env:"SECRET_FILE,validate=file,unset"`  // SECRET_FILE=/home/jean/.ssh/private.key
}

var MyEnvConfig MyEnv

func init() {
    cfg, err := LoadFromEnv(Config[MyEnv]{
        // use a ".env" file
        UseEnvFile: true,               // use an env file, not just environment variables (dotenv)
        EnvFileOveride: false,          // use environment variables (if they exist) over .env file
        EnvFilePath: "/home/jean/.env"  // optional - will recursively look for ".env" if not provided
    })
    if err != nil {
        panic(err)
    }
    MyEnvConfig = cfg
}

func main() {
    fmt.Println(cfg.Name)
    // Output: Jean-Luc Picard
    fmt.Println(os.Getenv("AGE")
    // Output: 25
}
```
