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
	ReadCsvPb("/Users/steve/git/gofun/src/glue/csv2Pb_test.pb")
}
