package file

import "testing"

func TestExtractLocations(t *testing.T) {
	tests := []struct {
		name   string
		output string
		path   string
		want   []string
	}{
		{
			name: "runtime error stack trace",
			path: "/private/tmp/runtimebad.js",
			output: "/private/tmp/runtimebad.js:1\n" +
				"console.log(z);\n" +
				"            ^\n\n" +
				"ReferenceError: z is not defined\n" +
				"    at Object.<anonymous> (/private/tmp/runtimebad.js:1:13)\n" +
				"    at Module._compile (node:internal/modules/cjs/loader:1934:14)\n",
			want: []string{"/private/tmp/runtimebad.js:1:13"},
		},
		{
			name: "syntax error with no stack frame",
			path: "/private/tmp/syntaxbad.js",
			output: "/private/tmp/syntaxbad.js:1\n" +
				"console.log(x\n" +
				"            ^\n\n" +
				"SyntaxError: missing ) after argument list\n" +
				"    at wrapSafe (node:internal/modules/cjs/loader:1866:18)\n",
			want: []string{"/private/tmp/syntaxbad.js:1"},
		},
		{
			name:   "no location in output",
			path:   "/tmp/clean.js",
			output: "hi\n",
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			locs := extractLocations(tt.output, tt.path)

			if len(locs) != len(tt.want) {
				t.Fatalf("got %d locations %v, want %d %v", len(locs), locs, len(tt.want), tt.want)
			}
			for i, loc := range locs {
				if loc.text != tt.want[i] {
					t.Errorf("location %d = %q, want %q", i, loc.text, tt.want[i])
				}
			}
		})
	}
}

func TestNearestLocation(t *testing.T) {
	locs := []location{
		{text: "file.js:1:1", idx: 10},
		{text: "file.js:5:1", idx: 100},
		{text: "file.js:9:1", idx: 200},
	}

	tests := []struct {
		name string
		pos  int
		want string
	}{
		{name: "closest to first", pos: 12, want: "file.js:1:1"},
		{name: "closest to middle", pos: 90, want: "file.js:5:1"},
		{name: "closest to last", pos: 500, want: "file.js:9:1"},
		{name: "exact match", pos: 100, want: "file.js:5:1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nearestLocation(locs, tt.pos); got != tt.want {
				t.Errorf("nearestLocation(locs, %d) = %q, want %q", tt.pos, got, tt.want)
			}
		})
	}

	t.Run("no locations", func(t *testing.T) {
		if got := nearestLocation(nil, 0); got != "" {
			t.Errorf("nearestLocation(nil, 0) = %q, want empty string", got)
		}
	})
}

func TestAbs(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}

	for _, tt := range tests {
		if got := abs(tt.n); got != tt.want {
			t.Errorf("abs(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}
