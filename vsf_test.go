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
			name:            "Multiple columns",
			input:           "col1:col2:col3\nvalue1:value2:value3",
			delimiter:       ":",
			outputDelimiter: "",
			want:            "col1   : col2   : col3\nvalue1 : value2 : value3",
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

func TestFormatWithSeparator(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		delimiter       string
		outputDelimiter string
		afterLine       int
		sepChar         string
		want            string
		wantErr         bool
	}{
		{
			name:            "Add separator after header",
			input:           "Index:Directory\n5:/long/path\n0:/short",
			delimiter:       ":",
			outputDelimiter: "",
			afterLine:       0,
			sepChar:         "-",
			want:            "Index : Directory\n------:----------\n5     : /long/path\n0     : /short",
			wantErr:         false,
		},
		{
			name:            "Custom output delimiter with separator",
			input:           "name,age\njohn,30\namy,25",
			delimiter:       ",",
			outputDelimiter: "|",
			afterLine:       0,
			sepChar:         "=",
			want:            "name | age\n=====|====\njohn | 30\namy  | 25",
			wantErr:         false,
		},
		{
			name:            "Separator after middle line",
			input:           "line1:data1\nline2:data2\nline3:data3",
			delimiter:       ":",
			outputDelimiter: "",
			afterLine:       1,
			sepChar:         "-",
			want:            "line1 : data1\nline2 : data2\n------:------\nline3 : data3",
			wantErr:         false,
		},
		{
			name:            "Invalid line number (too high)",
			input:           "name:john\nage:30",
			delimiter:       ":",
			outputDelimiter: "",
			afterLine:       5,
			sepChar:         "-",
			want:            "name : john\nage  : 30",
			wantErr:         false,
		},
		{
			name:            "Invalid line number (negative)",
			input:           "name:john\nage:30",
			delimiter:       ":",
			outputDelimiter: "",
			afterLine:       -1,
			sepChar:         "-",
			want:            "name : john\nage  : 30",
			wantErr:         false,
		},
		{
			name:            "Single line with separator",
			input:           "single:line",
			delimiter:       ":",
			outputDelimiter: "",
			afterLine:       0,
			sepChar:         "=",
			want:            "single : line\n=======:=====",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatWithSeparator(tt.input, tt.delimiter, tt.outputDelimiter, tt.afterLine, tt.sepChar)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatWithSeparator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FormatWithSeparator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatSkipLines(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		delimiter       string
		outputDelimiter string
		skipLines       []int
		want            string
		wantErr         bool
	}{
		{
			name:            "Skip separator line",
			input:           "name:john\n----:----\nage:30",
			delimiter:       ":",
			outputDelimiter: "",
			skipLines:       []int{1},
			want:            "name : john\n----:----\nage  : 30",
			wantErr:         false,
		},
		{
			name:            "Skip header line from width calculation",
			input:           "VERY_LONG_HEADER:DESCRIPTION\nshort:val\nmed:value",
			delimiter:       ":",
			outputDelimiter: "",
			skipLines:       []int{0},
			want:            "VERY_LONG_HEADER:DESCRIPTION\nshort : val\nmed   : value",
			wantErr:         false,
		},
		{
			name:            "Skip multiple lines",
			input:           "header:info\n----:----\nname:john\n====:====\nage:30",
			delimiter:       ":",
			outputDelimiter: "",
			skipLines:       []int{0, 1, 3},
			want:            "header:info\n----:----\nname : john\n====:====\nage  : 30",
			wantErr:         false,
		},
		{
			name:            "Skip no lines (same as Format)",
			input:           "name:john\nage:30",
			delimiter:       ":",
			outputDelimiter: "",
			skipLines:       []int{},
			want:            "name : john\nage  : 30",
			wantErr:         false,
		},
		{
			name:            "Skip invalid line numbers",
			input:           "name:john\nage:30",
			delimiter:       ":",
			outputDelimiter: "",
			skipLines:       []int{-1, 5, 10},
			want:            "name : john\nage  : 30",
			wantErr:         false,
		},
		{
			name:            "Skip all lines",
			input:           "line1:data1\nline2:data2",
			delimiter:       ":",
			outputDelimiter: "",
			skipLines:       []int{0, 1},
			want:            "line1:data1\nline2:data2",
			wantErr:         false,
		},
		{
			name:            "Custom output delimiter with skipped lines",
			input:           "HEADER,INFO\nname,john\nage,30",
			delimiter:       ",",
			outputDelimiter: "|",
			skipLines:       []int{0},
			want:            "HEADER,INFO\nname | john\nage  | 30",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatSkipLines(tt.input, tt.delimiter, tt.outputDelimiter, tt.skipLines)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatSkipLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FormatSkipLines() = %v, want %v", got, tt.want)
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

