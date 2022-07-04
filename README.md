# gotest

Like `go test` but with colors.

## Installation

### Pre-compiled binary

Download it from the [OSS release page](https://github.com/wesleimp/gotest/releases)

### Manually

**clone the repository**

```sh
$ git clone https://github.com/wesleimp/gotest
```

**build**

```sh
$ go build -o gotest main.go
```

Now you can copy the binary `gotest` to your `bin` folder or any other desired location

## Usage

Accepts all the arguments and flags `go test` works with.

Examples:

```sh
$ gotest ./...
```

```sh
$ gotest -v github.com/jonasbn/go-test-demo
```

## Customization

You can customize the color of the output by setting the following env variable with pattern `<PASS>,<FAIL>,<SKIP>`

Example:

```
$ GOTEST_COLORS="#b8bb26,#fb4934,#fabd2f"
```

See [lipgloss][1] for more info about colors

## Acknowledgements

This package is a modified version of the great [rakyll/gotest][2].

Changes are mostly related to customization, distribution and some bug fixes.

Make sure to check out the [original work by rakyll][1]!

[1]: https://github.com/charmbracelet/lipgloss
[2]: https://github.com/rakyll/gotest
