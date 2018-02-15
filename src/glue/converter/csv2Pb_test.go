package converter

import (
	"fmt"
	"testing"
)

func TestCsv2Pb(t *testing.T) {
	fmt.Println("TestCsv2Pb")
	// read csv and write pb
	Csv2Pb("/Users/steve/git/gofun/src/glue/csv2Protobuf_test.csv", "/Users/steve/git/gofun/src/glue/csv2Pb_test.pb")
	// read and verify
	dataFile := ReadCsvPb("/Users/steve/git/gofun/src/glue/csv2Pb_test.pb")
	for i, row := range dataFile.Rows {
		if i == 0 {
			validateKey(row.GetKey(), -1, t)
		}
		fmt.Print("[", row.GetKey())
		for j, val := range row.GetData() {
			if i == 0 && j == 0 {
				validateRow(val, 0, t)
			} else if i == 1 && j == 0 {
				validateRow(val, 27.80985, t)
			}
			fmt.Print(" ", val)
		}
		fmt.Println("]")
	}

}

func validateKey(key int32, expKey int32, t *testing.T) {
	if key != expKey {
		t.Error("key:", key, " != expKey:", expKey)
	}
}

func validateRow(val float32, expVal float32, t *testing.T) {
	if val != expVal {
		t.Error("val:", val, " != expVal:", expVal)
	}
}
