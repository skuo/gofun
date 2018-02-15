package converter

import (
	"encoding/binary"
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

	protoFile, err := os.Create(protobufFname)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer protoFile.Close()

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
		// https://developers.google.com/protocol-buffers/docs/techniques
		// write size of pb msg for a file containing multiple msgs
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(len(dataRowProtobuf)))
		_, err = protoFile.Write(b)
		_, err = protoFile.Write(dataRowProtobuf)
		fmt.Println(len(dataRowProtobuf), dataRowProtobuf)
	}
}

// ReadCsvProtobuf returns [][]string
func ReadCsvProtobuf(protobufFname string) [][]string {
	protoFile, err := os.Open(protobufFname)
	if err != nil {
		fmt.Println("Error opening file", err)
		os.Exit(1)
	}
	defer protoFile.Close()

	sizeBuf := make([]byte, 4)
	for {
		// first read in the size
		n, err := protoFile.Read(sizeBuf)
		if n != 4 || err != nil {
			fmt.Println("Error reading size", err)
			os.Exit(1)
		}
		size := int(binary.LittleEndian.Uint32(sizeBuf))
		// read the data row
		dataRowProtobuf := make([]byte, size)
		n, err = protoFile.Read(dataRowProtobuf)
		if err != nil {
			fmt.Println("Error reading file", err)
			os.Exit(1)
		}
		fmt.Println("num of bytes read", n)
		dataRow := new(DataRow)
		if err = proto.Unmarshal(dataRowProtobuf, dataRow); err != nil {
			fmt.Println("Error unmarshalling dataRowProtobuf")
			os.Exit(1)
		}
		fmt.Println("dataRowPb", dataRow)
		fmt.Println(dataRow)
	}
}
