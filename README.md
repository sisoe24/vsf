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
aligned, err := vsf.Format(input, ":", ":")
if err != nil {
    log.Fatal(err)
}
fmt.Println(aligned)
```

## CLI Usage

```
vsf [-d delimiter] [-o output_delimiter] [-header lines] [-h]
```

**Flags:**
- `-d` : Input delimiter (default: ":")
- `-o` : Output delimiter (default: same as input)
- `-header` : Number of header lines to preserve without alignment (default: 0)
- `-h` : Show help

## Examples

### CLI Examples

* Basic formatting

  ```bash
  echo "name:john\nage:30" | vsf
  ```

* Git branch formatting

  ```bash
  git for-each-ref --sort=-committerdate refs/heads/ --format='%(refname:short):%(committerdate:short)' | vsf -o "|" | fzf
  ```

* CSV with headers (perfect with fzf)

  ```bash
  cat data.csv | vsf -d ',' -header 1 | fzf --header-lines 1
  ```

* Custom delimiters

  ```bash
  echo "a|b|c\nlong text|short|medium" | vsf -d '|' -o ' → '
  ```

### Library Examples

* Basic column alignment

  ```go
  data := "name:john\nage:30\ncity:new york"
  result, _ := vsf.Format(data, ":", ":")
  // Output:
  // name : john
  // age  : 30  
  // city : new york
  ```

* CSV with headers (perfect for fzf)

  ```go
  csv := "name,age,city\njohn,30,nyc\namy,25,rome"
  result, _ := vsf.FormatWithHeader(csv, ",", " | ", 1)
  // Output:
  // name,age,city
  // john | 30  | nyc
  // amy  | 25  | rome
  ```

* Parse individual lines

  ```go
  columns := vsf.ParseLine("name:'john doe':30", ":")
  // columns = ["name", "'john doe'", "30"]
  // Respects quotes - won't split on delimiters inside quotes
  ```

* Custom output delimiter

  ```go
  result, _ := vsf.Format(data, ":", " | ")
  // Output:
  // name | john
  // age  | 30
  // city | new york
  ```

## API Reference

### Functions

- `Format(input, delimiter, outputDelimiter string) (string, error)` - Standard column alignment
- `FormatWithHeader(input, delimiter, outputDelimiter string, headerLines int) (string, error)` - Column alignment preserving header lines
- `ParseLine(line, delimiter string) []string` - Parse a single line into columns (respects quotes)

## Development

```bash
make       # Build
make test  # Run tests
```

## Contributing

PRs welcome!

## License

MIT License. See [LICENSE](LICENSE) file for more info.
