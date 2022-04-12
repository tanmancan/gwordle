# Gwordle - Alpha

A word game written in Go inspired by Wordle.

## Run from source

Requires `go` to be installed: https://go.dev/doc/install

Clone or download source from this repository. Then from project root run:

```bash
go run cmd/cli/main.go
```

Then follow on screen instructions and enjoy:

```bash

Enter a guess word, or multiple words separate by space.
Type /help for more option.
You have 6 tries: Apple

A P P L E
_ _ _ _ _
_ _ _ _ _
_ _ _ _ _
_ _ _ _ _
_ _ _ _ _

Enter a guess word, or multiple words separate by space.
Type /help for more option.
You have 5 tries:
```

## Feature Roadmap

- Customization of word length and number of tries
- gRPC server/client
- RESTful server/client
- Web assembly service
- Web based UI
