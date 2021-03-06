package converter

import (
	"fmt"
	"os"
	"os/exec"
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

func TestGenCsvFiles(t *testing.T) {
	// 100x
	os.Remove("/Users/steve/git/gofun/src/glue/mt_bruno_elevation_100.csv")
	for i := 0; i < 100; i++ {
		cmd := exec.Command("bash", "-c", "cat /Users/steve/git/gofun/src/glue/mt_bruno_elevation.csv >> /Users/steve/git/gofun/src/glue/mt_bruno_elevation_100.csv")
		err := cmd.Run()
		checkError(err)
	}
	// 10,000x
	os.Remove("/Users/steve/git/gofun/src/glue/mt_bruno_elevation_10K.csv")
	for i := 0; i < 100; i++ {
		cmd := exec.Command("bash", "-c", "cat /Users/steve/git/gofun/src/glue/mt_bruno_elevation_100.csv >> /Users/steve/git/gofun/src/glue/mt_bruno_elevation_10K.csv")
		err := cmd.Run()
		checkError(err)
	}
	// 100,000x
	os.Remove("/Users/steve/git/gofun/src/glue/mt_bruno_elevation_100K.csv")
	for i := 0; i < 10; i++ {
		cmd := exec.Command("bash", "-c", "cat /Users/steve/git/gofun/src/glue/mt_bruno_elevation_10K.csv >> /Users/steve/git/gofun/src/glue/mt_bruno_elevation_100K.csv")
		_, err := cmd.Output()
		checkError(err)
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func TestGenPbFiles(t *testing.T) {
	fmt.Println("TestGenPbFiles")
	// read csv and write pb
	Csv2Pb("/Users/steve/git/gofun/src/glue/mt_bruno_elevation.csv", "/Users/steve/git/gofun/src/glue/mt_bruno_elevation.pb")
	Csv2Pb("/Users/steve/git/gofun/src/glue/mt_bruno_elevation_100.csv", "/Users/steve/git/gofun/src/glue/mt_bruno_elevation_100.pb")
	Csv2Pb("/Users/steve/git/gofun/src/glue/mt_bruno_elevation_10K.csv", "/Users/steve/git/gofun/src/glue/mt_bruno_elevation_10K.pb")
	Csv2Pb("/Users/steve/git/gofun/src/glue/mt_bruno_elevation_100K.csv", "/Users/steve/git/gofun/src/glue/mt_bruno_elevation_100K.pb")
	//Csv2Pb("/Users/steve/git/gofun/src/glue/mt_bruno_elevation_1M.csv", "/Users/steve/git/gofun/src/glue/mt_bruno_elevation_1M.pb")
}
