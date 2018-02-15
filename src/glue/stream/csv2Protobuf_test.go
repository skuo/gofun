package converter

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestIntToBytes(t *testing.T) {
	fmt.Println("TestIntToBytes")

	size := 12958723
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(size))
	bRead := int(binary.LittleEndian.Uint32(b))
	fmt.Println("bRead", bRead)
	if bRead != size {
		t.Error("bRead != size")
	}
}

func TestCsv2Protocol(t *testing.T) {
	fmt.Println("TestCsv2Protocol")
	// read csv and write pb
	Csv2Protobuf("/Users/steve/git/gofun/src/glue/csv2Protobuf_test.csv", "/Users/steve/git/gofun/src/glue/csv2Protobuf_test.pb")
	// read and verify
	ReadCsvProtobuf("/Users/steve/git/gofun/src/glue/csv2Protobuf_test.pb")
}
