# vsf - Very Simple Formatter

A Go CLI tool for column alignment and text formatting. Works great with fzf.

## What it does

- Aligns text columns based on delimiters
- Integrates smoothly with fzf and Unix pipelines
- Handy for formatting Git logs, CSVs, and other structured text

## Quick start

```bash
go install github.com/sisoe24/vsf
echo "name:age\njohnny:30\namy:25" | vsf | fzf --header-lines 1
```

## Basic usage

```
vsf [-d delimiter] [-o output_delimiter] [-h]
```

## Some useful examples

### Git branch formatting

```bash
git for-each-ref --sort=-committerdate refs/heads/ --format='%(refname:short):%(committerdate:short)' | vsf -o "|" | fzf
```

### CSV exploration

```bash
cat data.csv | vsf -d ',' | fzf --header-lines 1
```

## Development

```bash
make  # Build
make test  # Run tests
```

## Contributing

PRs welcome!

## License

MIT License. See [LICENSE](LICENSE) file for more info.
