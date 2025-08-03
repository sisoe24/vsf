package main

import "testing"

// Test helper functions
func TestParseLineNumbers(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{
			name:    "Valid numbers",
			input:   "1,3,5",
			want:    []int{1, 3, 5},
			wantErr: false,
		},
		{
			name:    "Single number",
			input:   "2",
			want:    []int{2},
			wantErr: false,
		},
		{
			name:    "Empty string",
			input:   "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "With spaces",
			input:   " 1 , 3 , 5 ",
			want:    []int{1, 3, 5},
			wantErr: false,
		},
		{
			name:    "Invalid number",
			input:   "1,abc,3",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLineNumbers(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLineNumbers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("parseLineNumbers() = %v, want %v", got, tt.want)
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("parseLineNumbers() = %v, want %v", got, tt.want)
					break
				}
			}
		})
	}
}
