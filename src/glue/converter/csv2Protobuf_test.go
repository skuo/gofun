package converter

import (
	"fmt"
	"testing"
)

func TestCsv2Protocol(t *testing.T) {
	fmt.Println("TestCsv2Protocol")
	Csv2Protobuf("/Users/steve/git/gofun/src/glue/mt_bruno_elevation.csv", "/Users/steve/git/gofun/src/glue/mt_bruno_elevation.pb")
}
