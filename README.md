# Markdown Simple Static Site Generator

This is a very simple markdown to HTML generator for creating static sites. It can be used for anything but was created with personal notes & snippets and documentation in mind.

Yes there are MANY other static site generators, but this one is mine, it does what I need

Goals:

- Convert markdown to HTML
- Index folders, creating index.html as it goes
- Use custom themes, wow!
- Single static binary
- Minimal config

Use cases & key features:

- Log personal work
- Notekeeping
- Informal docs
- Code snippets

Supporting technologies and libraries:

- https://github.com/gomarkdown

![](https://img.shields.io/github/license/benc-uk/site-generator)
![](https://img.shields.io/github/last-commit/benc-uk/site-generator)
![](https://img.shields.io/github/release/benc-uk/site-generator)
![](https://img.shields.io/github/checks-status/benc-uk/site-generator/main)

# Screenshots

<img src="https://user-images.githubusercontent.com/14982936/181744683-61925fdf-62de-432a-9db6-8dffe741deb5.png" width="500px">

# Getting Started

## Installing

You can install straight from GitHub using `go install` pick a tag or use `latest`

```bash
go install github.com/benc-uk/site-generator@latest
```

## Running locally

- Have Go SDK installed
- Clone this repo `git clone github.com/benc-uk/site-generator`
- Use make (see below)

```text
help                 💬 This help message :)
lint                 🌟 Lint & format, will not fix but sets exit code on error
lint-fix             🔍 Lint & format, will try to fix errors and modify code
build                🔨 Run a local build without a container
run                  🚀 Run application, used for local development
```

## Command Usage

```text
$ site-generator --help
🧵 Simple Site Generator v0.0.1 (Manual build)

Usage:
  -o string
        Output HTML and site content here (default "./html")
  -s string
        Source directory containing Markdown files (default "./src")
  -t string
        Optional, custom template file
```

# Templates

See the [example template](./templates/example.html) for the basics of how to build your own. [Go templates](https://pkg.go.dev/html/template) are used, and the file is expected to be a single HTML file, preferably with an embedded stylesheet.

The `default.html` template provides more details on how to style the output, and is embedded into the binary by default

# Known Issues

Do not have your source and output directories overlapping!

# License

This project uses the MIT software license. See [full license file](./LICENSE)
