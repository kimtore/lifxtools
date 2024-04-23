# LIFXTOOL

This program runs effects on LIFX bulbs.

## Getting started

1. Create a configuration file; see [example configuration](./config.example.yaml) for syntax. Save the file as `config.yaml`.
2. Compile the program (see below)
3. Run the program: `lifxtool -preset <mypreset>`

## Compiling

1. [Download and install Go 1.17 or higher](https://go.dev/doc/install)
2. Checkout the repository: `git clone https://github.com/dorkowscy/lifxtool`
3. Compile the program: `cd lifxtool; go build -o lifxtool cmd/lifxtool/main.go`
