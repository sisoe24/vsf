# vsf - Very Simple Formatter

A Go CLI tool and library for column alignment and text formatting. Works great with fzf.

## What it does

- Aligns text columns based on delimiters
- Integrates with fzf and Unix pipelines
- Handy for formatting Git logs, CSVs, and other structured text
- Available as both a CLI tool and Go library

## Quick start

### CLI Tool

```bash
go install github.com/sisoe24/vsf/cmd/vsf@latest
echo "name:city:age\njohnny:new york:30\namy:rome:25" | vsf | fzf --header-lines 1
```

### Go Library

```bash
go get github.com/sisoe24/vsf
```

```go
import "github.com/sisoe24/vsf"

input := "name:city:age\njohnny:new york:30\namy:rome:25"
aligned, err := vsf.AlignColumns(input, ":", ":")
if err != nil {
    log.Fatal(err)
}
fmt.Println(aligned)
```

## CLI Usage

```
vsf [-d delimiter] [-o output_delimiter] [-h]
```

## Examples

### CLI Examples

* Git branch formatting

  ```bash
  git for-each-ref --sort=-committerdate refs/heads/ --format='%(refname:short):%(committerdate:short)' | vsf -o "|" | fzf
  ```

* CSV

  ```bash
  cat data.csv | vsf -d ',' | fzf --header-lines 1
  ```

### Library Examples

* Basic column alignment

  ```go
  data := "name:john\nage:30\ncity:new york"
  result, _ := vsf.AlignColumns(data, ":", ":")
  // Output:
  // name : john
  // age  : 30  
  // city : new york
  ```

* Parse individual lines

  ```go
  columns := vsf.ParseLine("name:john doe:30", ":")
  // columns = ["name", "john doe", "30"]
  ```

* Custom output delimiter

  ```go
  result, _ := vsf.AlignColumns(data, ":", " | ")
  // Output:
  // name | john
  // age  | 30
  // city | new york
  ```

## API Reference

### Functions

- `AlignColumns(input, delimiter, outputDelimiter string) (string, error)` - Main formatting function
- `ParseLine(line, delimiter string) []string` - Parse a single line into columns

## Development

```bash
make       # Build
make test  # Run tests
```

## Contributing

PRs welcome!

## License

MIT License. See [LICENSE](LICENSE) file for more info.
