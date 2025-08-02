package vsf

import (
	"testing"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		delimiter       string
		outputDelimiter string
		want            string
		wantErr         bool
	}{
		{
			name:            "Basic test",
			input:           "name:john\nage:30\ncity:new york",
			delimiter:       ":",
			outputDelimiter: "",
			want:            "name : john\nage  : 30\ncity : new york",
			wantErr:         false,
		},
		{
			name:            "Different lengths",
			input:           "short:value\nvery_long_key:short",
			delimiter:       ":",
			outputDelimiter: "",
			want:            "short         : value\nvery_long_key : short",
			wantErr:         false,
		},
		{
			name:            "Custom delimiter",
			input:           "field1,value1\nfield2,value2",
			delimiter:       ",",
			outputDelimiter: "",
			want:            "field1 , value1\nfield2 , value2",
			wantErr:         false,
		},
		{
			name:            "Custom output delimiter",
			input:           "field1,value1\nfield2,value2",
			delimiter:       ",",
			outputDelimiter: "|",
			want:            "field1 | value1\nfield2 | value2",
			wantErr:         false,
		},
		{
			name:            "Space delimiter",
			input:           "field1 value1\nfield2 value2",
			delimiter:       " ",
			outputDelimiter: "|",
			want:            "field1 | value1\nfield2 | value2",
			wantErr:         false,
		},
		{
			name:            "Space delimiter with quotes",
			input:           "'long field' value\n\"long field2\" value2",
			delimiter:       " ",
			outputDelimiter: "|",
			want:            "'long field'  | value\n\"long field2\" | value2",
			wantErr:         false,
		},
		{
			name:            "Single line",
			input:           "key:value",
			delimiter:       ":",
			outputDelimiter: "",
			want:            "key : value",
			wantErr:         false,
		},
		{
			name:            "Multiple delimiters",
			input:           "col1:col2:col3\nvalue1:value2:value3",
			delimiter:       ":",
			outputDelimiter: "",
			want:            "col1   : col2   : col3\nvalue1 : value2 : value3",
			wantErr:         false,
		},
		{
			name:            "Input with empty lines",
			input:           "key1:value1\n\nkey2:value2",
			delimiter:       ":",
			outputDelimiter: "",
			want:            "key1 : value1\n\nkey2 : value2",
			wantErr:         false,
		},
		{
			name:            "Input with spaces",
			input:           "  key1  :  value1  \n  key2  :  value2  ",
			delimiter:       ":",
			outputDelimiter: "",
			want:            "key1 : value1\nkey2 : value2",
			wantErr:         false,
		},
		{
			name:            "Delimiter not in input",
			input:           "line1\nline2\nline3",
			delimiter:       ":",
			outputDelimiter: "",
			want:            "line1\nline2\nline3",
			wantErr:         false,
		},
		{
			name:            "Empty input",
			input:           "",
			delimiter:       ":",
			outputDelimiter: "",
			want:            "",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(tt.input, tt.delimiter, tt.outputDelimiter)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatWithHeader(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		delimiter       string
		outputDelimiter string
		headerLines     int
		want            string
		wantErr         bool
	}{
		{
			name:            "CSV with single header",
			input:           "name,age,city\njohn,30,nyc\namy,25,rome",
			delimiter:       ",",
			outputDelimiter: "|",
			headerLines:     1,
			want:            "name,age,city\njohn | 30 | nyc\namy  | 25 | rome",
			wantErr:         false,
		},
		{
			name:            "Header preserved as-is",
			input:           "NAME | AGE | CITY\njohn:30:nyc\namy:25:rome",
			delimiter:       ":",
			outputDelimiter: "|",
			headerLines:     1,
			want:            "NAME | AGE | CITY\njohn | 30 | nyc\namy  | 25 | rome",
			wantErr:         false,
		},
		{
			name:            "Multiple headers",
			input:           "Data Report\nname,age,city\njohn,30,nyc\namy,25,rome",
			delimiter:       ",",
			outputDelimiter: "|",
			headerLines:     2,
			want:            "Data Report\nname,age,city\njohn | 30 | nyc\namy  | 25 | rome",
			wantErr:         false,
		},
		{
			name:            "Zero headers (same as Format)",
			input:           "name:john\nage:30",
			delimiter:       ":",
			outputDelimiter: "",
			headerLines:     0,
			want:            "name : john\nage  : 30",
			wantErr:         false,
		},
		{
			name:            "Headers equal to total lines",
			input:           "header1\nheader2",
			delimiter:       ":",
			outputDelimiter: "",
			headerLines:     2,
			want:            "header1\nheader2",
			wantErr:         false,
		},
		{
			name:            "Headers more than total lines",
			input:           "line1\nline2",
			delimiter:       ":",
			outputDelimiter: "",
			headerLines:     5,
			want:            "line1\nline2",
			wantErr:         false,
		},
		{
			name:            "Negative header count",
			input:           "name:john\nage:30",
			delimiter:       ":",
			outputDelimiter: "",
			headerLines:     -1,
			want:            "name : john\nage  : 30",
			wantErr:         false,
		},
		{
			name:            "Single data line with header",
			input:           "name,age\njohn,30",
			delimiter:       ",",
			outputDelimiter: "|",
			headerLines:     1,
			want:            "name,age\njohn | 30",
			wantErr:         false,
		},
		{
			name:            "Header with different delimiters",
			input:           "name|age|city\njohn:30:nyc\namy:25:rome",
			delimiter:       ":",
			outputDelimiter: "→",
			headerLines:     1,
			want:            "name|age|city\njohn → 30 → nyc\namy  → 25 → rome",
			wantErr:         false,
		},
		{
			name:            "Empty input with headers",
			input:           "",
			delimiter:       ":",
			outputDelimiter: "",
			headerLines:     1,
			want:            "",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatWithHeader(tt.input, tt.delimiter, tt.outputDelimiter, tt.headerLines)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatWithHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FormatWithHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		delimiter string
		want      []string
	}{
		{
			name:      "Basic parsing",
			line:      "name:age:city",
			delimiter: ":",
			want:      []string{"name", "age", "city"},
		},
		{
			name:      "With spaces",
			line:      "john doe:30:new york",
			delimiter: ":",
			want:      []string{"john doe", "30", "new york"},
		},
		{
			name:      "Single quotes",
			line:      "'john:doe':30:'new:york'",
			delimiter: ":",
			want:      []string{"'john:doe'", "30", "'new:york'"},
		},
		{
			name:      "Double quotes",
			line:      `"hello,world",foo,"bar,baz"`,
			delimiter: ",",
			want:      []string{`"hello,world"`, "foo", `"bar,baz"`},
		},
		{
			name:      "Mixed quotes",
			line:      `'single':30:"double:quoted"`,
			delimiter: ":",
			want:      []string{"'single'", "30", `"double:quoted"`},
		},
		{
			name:      "No delimiter",
			line:      "single line",
			delimiter: ":",
			want:      []string{"single line"},
		},
		{
			name:      "Empty line",
			line:      "",
			delimiter: ":",
			want:      []string{},
		},
		{
			name:      "Multi-character delimiter",
			line:      "field1::field2::field3",
			delimiter: "::",
			want:      []string{"field1", "field2", "field3"},
		},
		{
			name:      "Spaces around values",
			line:      "  name  :  value  :  city  ",
			delimiter: ":",
			want:      []string{"name", "value", "city"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLine(tt.line, tt.delimiter)
			if len(got) != len(tt.want) {
				t.Errorf("ParseLine() = %v, want %v", got, tt.want)
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("ParseLine() = %v, want %v", got, tt.want)
					break
				}
			}
		})
	}
}
