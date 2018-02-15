package converter

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/golang/protobuf/proto"

	pb "glue/protobuf"
)

// Csv2Pb reads input csvFname and output the data to pbFname
func Csv2Pb(csvFname string, pbFname string) {
	csvFile, err := os.Open(csvFname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return
	}
	defer csvFile.Close()

	protoFile, err := os.Create(pbFname)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer protoFile.Close()

	dataFile := &pb.PbDataFile{}
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
		dataRow := &pb.PbDataRow{}
		for i, rec := range record {
			if i == 0 {
				value, _ := strconv.ParseInt(rec, 10, 32)
				val32 := int32(value)
				dataRow.Key = val32
			} else {
				value, _ := strconv.ParseFloat(rec, 32)
				val32 := float32(value)
				dataRow.Data = append(dataRow.Data, val32)
			}
		}
		dataFile.Rows = append(dataFile.Rows, dataRow)
	}
	// Write data file to disk.
	out, err := proto.Marshal(dataFile)
	if err != nil {
		log.Fatalln("Failed to encode data file:", err)
	}
	if err := ioutil.WriteFile(pbFname, out, 0644); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}
}

// ReadCsvPb returns *pb.PbDataFile
func ReadCsvPb(pbFname string) *pb.PbDataFile {
	// Read the existing data file
	in, err := ioutil.ReadFile(pbFname)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	dataFile := &pb.PbDataFile{}
	if err := proto.Unmarshal(in, dataFile); err != nil {
		log.Fatalln("Failed to parse data file:", err)
	}
	return dataFile
}
