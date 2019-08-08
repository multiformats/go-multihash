package multihash

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestSpec(t *testing.T) {
	file, err := os.Open("spec/multicodec/table.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = false
	reader.FieldsPerRecord = 4
	reader.TrimLeadingSpace = true

	values, err := reader.ReadAll()
	if err != nil {
		t.Error(err)
	}
	expectedFunctions := make(map[uint64]string, len(values)-1)
	if values[0][0] != "name" || values[0][1] != "tag" || values[0][2] != "code" {
		t.Fatal("table format has changed")
	}

	for _, v := range values[1:] {
		name := v[0]
		tag := v[1]
		codeStr := v[2]

		if tag != "multihash" {
			// not a multihash function
			continue
		}

		var code uint64
		if !strings.HasPrefix(codeStr, "0x") {
			t.Errorf("invalid multicodec code %q (%s)", codeStr, name)
			continue
		}

		i, err := strconv.ParseUint(codeStr[2:], 16, 64)
		if err != nil {
			t.Errorf("invalid multibase code %q (%s)", codeStr, name)
			continue
		}
		code = uint64(i)
		expectedFunctions[code] = name
	}

	for code, name := range Codes {
		expectedName, ok := expectedFunctions[code]
		if !ok {
			t.Errorf("multihash %q (%x) not defined in the spec", name, code)
			continue
		}
		if expectedName != name {
			t.Errorf("encoding %q (%x) has unexpected name %q", expectedName, code, name)
		}
	}
}

func TestSpecVectors(t *testing.T) {
	file, err := os.Open("spec/multihash/tests/values/test_cases.csv")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = false
	reader.FieldsPerRecord = 4
	reader.TrimLeadingSpace = true

	values, err := reader.ReadAll()
	if err != nil {
		t.Error(err)
	}
	if len(values) == 0 {
		t.Error("no test values")
		return
	}

	// check the header.
	if values[0][0] != "algorithm" ||
		values[0][1] != "bits" ||
		values[0][2] != "input" ||
		values[0][3] != "multihash" {
		t.Fatal("table format has changed")
	}

	for i, testCase := range values[1:] {
		function := testCase[0]
		lengthStr := testCase[1]
		input := testCase[2]
		expectedStr := testCase[3]

		t.Run(fmt.Sprintf("%d/%s/%s", i, function, lengthStr), func(t *testing.T) {
			code, ok := Names[function]
			if !ok {
				t.Skipf("skipping %s: not supported", function)
				return
			}

			length, err := strconv.ParseInt(lengthStr, 10, 64)
			if err != nil {
				t.Fatalf("failed to decode length: %s", err)
			}

			if length%8 != 0 {
				t.Fatal("expected the length to be a multiple of 8")
			}

			actual, err := Sum([]byte(input), code, int(length/8))
			if err != nil {
				t.Fatalf("failed to hash: %s", err)
			}
			actualStr := actual.HexString()
			if actualStr != expectedStr {
				t.Fatalf("got the wrong hash: expected %s, got %s", expectedStr, actualStr)
			}
		})

	}
}
