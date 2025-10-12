package cut

import (
	"os"
	"strings"
	"testing"
)

type testCase struct {
	name      string
	inputFile string
	wantFile  string
	opts      Options
}

func TestSortFromFiles(t *testing.T) {
	tests := []testCase{
		{
			name:      "simple test",
			inputFile: "test_cases/1_simple/input.txt",
			wantFile:  "test_cases/1_simple/simple_test.txt",
			opts: Options{
				Fields:    "1,3",
				Delimiter: ",",
			},
		},
		{
			name:      "sepatared flag test",
			inputFile: "test_cases/2_separated/input.txt",
			wantFile:  "test_cases/2_separated/separated_test.txt",
			opts: Options{
				Fields:    "1,3",
				Delimiter: ",",
				Separated: true,
			},
		},
		{
			name:      "range test",
			inputFile: "test_cases/3_range/input.txt",
			wantFile:  "test_cases/3_range/range_test.txt",
			opts: Options{
				Fields:    "2-3",
				Delimiter: ",",
			},
		},
		{
			name:      "more fields test",
			inputFile: "test_cases/3_range/input.txt",
			wantFile:  "test_cases/3_range/more_fields_test.txt",
			opts: Options{
				Fields:    "1,2-3",
				Delimiter: ",",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inBytes, err := os.ReadFile(tc.inputFile)
			if err != nil {
				t.Fatalf("can't read input: %v", err)
			}

			lines := string(inBytes)
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}

			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()
			os.Stdin = r

			_, err = w.Write([]byte(lines))
			if err != nil {
				t.Fatal(err)
			}
			_ = w.Close()

			resultChannel, err := GetResultChannel(tc.opts)
			if err != nil {
				t.Errorf("error: %v", err)
			}

			wantBytes, err := os.ReadFile(tc.wantFile)
			if err != nil {
				t.Fatalf("can't read want: %v", err)
			}
			wantLines := strings.Split(strings.ReplaceAll(string(wantBytes), "\r", ""), "\n")

			for i := range wantLines {
				got := <-resultChannel
				if strings.Join(got, tc.opts.Delimiter) != wantLines[i] {
					t.Errorf("[%d] got %q, want %q", i, strings.Join(got, tc.opts.Delimiter), wantLines[i])
				}
			}
			select {
			case got, ok := <-resultChannel:
				if ok {
					t.Errorf("got %q, want \"\"", got)
				}
			default:
			}
		})
	}
}
