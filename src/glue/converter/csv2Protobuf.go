package converter

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/golang/protobuf/proto"
)

// Csv2Protobuf reads input csvFname and output the data to protobufFname
func Csv2Protobuf(csvFname string, protobufFname string) {
	csvFile, err := os.Open(csvFname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record)

		// generate protobuf
		dataRow := new(DataRow)
		for i, rec := range record {
			value, _ := strconv.ParseInt(rec, 10, 64)
			val32 := int32(value)
			if i == 0 && rec != "" {
				dataRow.Key = val32
			} else {
				dataRow.Data = append(dataRow.Data, val32)
			}
		}
		dataRowProtobuf, err := proto.Marshal(dataRow)
		fmt.Println(dataRowProtobuf)
	}
}
