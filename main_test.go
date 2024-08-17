package main

import (
	"testing"
)

func TestFormatInput(t *testing.T) {
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
			got, err := formatInput(tt.input, tt.delimiter, tt.outputDelimiter)
			if (err != nil) != tt.wantErr {
				t.Errorf("formatInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("formatInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
