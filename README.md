# zerochi

zerochi provides a simple logging middleware for [chi](https://github.com/go-chi/chi) v5 using [zerolog](https://github.com/rs/zerolog). It will also log and recover panics.

## Usage

Grab the module:

```
go get -u github.com/Xiol/zerochi
```

Then set up the middleware:

```go
import (
    "github.com/rs/zerolog/log"
    "github.com/Xiol/zerochi"
    "github.com/go-chi/chi/v5"
)

func main() {
    r := chi.NewRouter()
    r.Use(zerochi.Logger(&log.Logger))  // To use the default logger

    // ...
}
```

If you have paths that you don't want cluttering up your logs, you can use `zerochi.SetSilentPaths()` to avoid logging requests to those paths:

```go
zerochi.SetSilentPaths("/ping", "/health")
```

The log message that accompanies each log line can also be changed:

```go
zerochi.LogMessage = "api request"
```
