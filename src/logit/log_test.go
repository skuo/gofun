// Package logit contains utility functions for logging for BeyondAI.
package logit

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

/*
	TestLogConfig1
	Test the json test configuration to verify proper configuration.
*/
func TestLogConfig1(t *testing.T) {
	//println("Testing: '" + t.Name() + "'")
	const configFileName = "logtestcfg.json"
	// This string gets written to a temporary file for testing
	const jsonTest string = `
	{
		"SiteID": "BeyondAI",
		"SystemID": "local",
		"filename": "logTestFile.txt",
		"url": "",
		"stdout": true,
		"level": "DEBUG",
		"debugFlags": [
			{ "pkg": "main"},
			{ "pkg": "logit", "file": "" },
			{ "pkg": "logit", "file": "log_test" }
		],
		"xFlags": [
			{ "pkg": "main", "file": "main", "flags": "3" },
			{ "pkg": "logit", "file": "" , "flags": "0x10" },
			{ "pkg": "logit", "file": "log" , "flags": "0x1" },
			{ "pkg": "logit", "file": "log_test" , "flags": "0xff00" }
		]
	}`

	err := writeConfigFile(configFileName, jsonTest)
	if err != nil {
		t.Errorf("Logit problem %s", err.Error())
	}
	defer removeConfigFile(configFileName)
	err = OpenLog(configFileName)
	if err != nil {
		t.Errorf("Logit problem %s", err.Error())
	}
	//
	// now test the results from loading the test configuration
	//
	var myFlags DFlags_t
	GetMyLogInfo(&myFlags)
	if myFlags.pkgName != "logit" {
		t.Errorf("Logit problem: package name '%s' is incorrect", myFlags.pkgName)
	}
	if myFlags.fileName != "log_test" {
		t.Errorf("Logit problem: file name '%s' is incorrect", myFlags.fileName)
	}
	if myFlags.dFlag {
		t.Errorf("Logit problem: debug flag is incorrect")
	}
	if myFlags.xFlag != 0xff10 {
		t.Errorf("Logit problem: xflag of %d is incorrect", myFlags.xFlag)
	}
	//
	// Now write some logs
	//
	startingLineCount := GetLogStats().lineCount
	Fatal(&myFlags, "Fatal test message 1 of 02")
	Fatalf(&myFlags, "Fatal test message %d of %02d", 2, 2)
	Warn(&myFlags, "Warning test message 1 of 02")
	Warnf(&myFlags, "Warning test message %d of %02d", 2, 2)
	Info(&myFlags, "Info test message 1 of 02")
	Infof(&myFlags, "Info test message %d of %02d", 2, 2)
	// Because of the config, only the 2nd and 3rd debug message should log
	// This tests the exclusive or nature of the debug flags
	Debug(&myFlags, "Debug test message 1 of 03")
	myFlags.dFlag = true
	Debug(&myFlags, "Debug test message 2 of 03")
	Debugf(&myFlags, "Debug test message %d of %02d", 3, 3)
	// get the size of the log
	lineCount := GetLogStats().lineCount
	if (lineCount - startingLineCount) != 8 {
		t.Errorf("Logit problem: Number of logged messages (%d) is incorrect", lineCount)
	}
	CloseLog()
}

/*
  TestLogConfig2
  Test to make sure debug logging is turned off when the level is set to INFO
*/
func TestLogConfig2(t *testing.T) {
	const configFileName = "logtestcfg.json"
	// This string gets written to a temporary file for testing
	const jsonTest string = `
	{
		"SiteID": "BeyondAI",
		"SystemID": "local",
		"filename": "logTestFile.txt",
		"url": "",
		"stdout": true,
		"level": "INFO",
		"debugFlags": [
			{ "pkg": "main"},
			{ "pkg": "logit", "file": "" },
			{ "pkg": "logit", "file": "log_test" }
		],
		"xFlags": [
			{ "pkg": "main", "file": "main", "flags": "3" },
			{ "pkg": "logit", "file": "" , "flags": "0x10" },
			{ "pkg": "logit", "file": "log" , "flags": "0x1" },
			{ "pkg": "logit", "file": "log_test" , "flags": "0xff00" }
		]
	}`

	err := writeConfigFile(configFileName, jsonTest)
	if err != nil {
		t.Errorf("Logit problem %s", err.Error())
	}
	defer removeConfigFile(configFileName)
	err = OpenLog(configFileName)
	if err != nil {
		t.Errorf("Logit problem %s", err.Error())
	}
	//
	// now test the results from loading the test configuration
	//
	var myFlags DFlags_t
	GetMyLogInfo(&myFlags)
	if myFlags.pkgName != "logit" {
		t.Errorf("Logit problem: package name '%s' is incorrect", myFlags.pkgName)
	}
	if myFlags.fileName != "log_test" {
		t.Errorf("Logit problem: file name '%s' is incorrect", myFlags.fileName)
	}
	if myFlags.dFlag {
		t.Errorf("Logit problem: debug flag is incorrect")
	}
	if myFlags.xFlag != 0xff10 {
		t.Errorf("Logit problem: xflag of %d is incorrect", myFlags.xFlag)
	}
	//
	// Now write some logs
	//
	startingLineCount := GetLogStats().lineCount
	Fatal(&myFlags, "Fatal test message 1 of 02")
	Fatalf(&myFlags, "Fatal test message %d of %02d", 2, 2)
	Error(&myFlags, "Error test message 1 of 02")
	Errorf(&myFlags, "Error test message %d of %02d", 2, 2)
	Warn(&myFlags, "Warning test message 1 of 02")
	Warnf(&myFlags, "Warning test message %d of %02d", 2, 2)
	Info(&myFlags, "Info test message 1 of 02")
	Infof(&myFlags, "Info test message %d of %02d", 2, 2)
	// Because of the config, none of these should output
	Debug(&myFlags, "Debug test message 1 of 03")
	Debug(&myFlags, "Debug test message 2 of 03")
	Debugf(&myFlags, "Debug test message %d of %02d", 3, 3)
	// get the size of the log
	lineCount := GetLogStats().lineCount
	if (lineCount - startingLineCount) != 8 {
		t.Errorf("Logit problem: Number of logged messages (%d) is incorrect", lineCount)
	}
	CloseLog()
}

/*
  TestLogConfig2
  Test to make sure all debug logging is turned oon when
  the level is set to INFO, and the package is set to "all"
*/
func TestLogConfig3(t *testing.T) {
	const configFileName = "logtestcfg.json"
	// This string gets written to a temporary file for testing
	const jsonTest string = `
	{
		"SiteID": "BeyondAI",
		"SystemID": "local",
		"filename": "logTestFile.txt",
		"url": "",
		"stdout": true,
		"level": "DEBUG",
		"debugFlags": [
			{ "pkg": "all"}
		],
		"xFlags": [
			{ "pkg": "main", "file": "main", "flags": "3" },
			{ "pkg": "logit", "file": "" , "flags": "0x10" },
			{ "pkg": "logit", "file": "log" , "flags": "0x1" },
			{ "pkg": "logit", "file": "log_test" , "flags": "0xff00" }
		]
	}`

	err := writeConfigFile(configFileName, jsonTest)
	if err != nil {
		t.Errorf("Logit problem %s", err.Error())
	}
	defer removeConfigFile(configFileName)
	err = OpenLog(configFileName)
	if err != nil {
		t.Errorf("Logit problem %s", err.Error())
	}
	//
	// now test the results from loading the test configuration
	//
	var myFlags DFlags_t
	GetMyLogInfo(&myFlags)
	if myFlags.pkgName != "logit" {
		t.Errorf("Logit problem: package name '%s' is incorrect", myFlags.pkgName)
	}
	if myFlags.fileName != "log_test" {
		t.Errorf("Logit problem: file name '%s' is incorrect", myFlags.fileName)
	}
	if !myFlags.dFlag {
		t.Errorf("Logit problem: debug flag is incorrect")
	}
	if myFlags.xFlag != 0xff10 {
		t.Errorf("Logit problem: xflag of %d is incorrect", myFlags.xFlag)
	}
	//
	// Now write some logs
	//
	startingLineCount := GetLogStats().lineCount
	Fatal(&myFlags, "Fatal test message 1 of 02")
	Fatalf(&myFlags, "Fatal test message %d of %02d", 2, 2)
	Warn(&myFlags, "Warning test message 1 of 02")
	Warnf(&myFlags, "Warning test message %d of %02d", 2, 2)
	Info(&myFlags, "Info test message 1 of 02")
	Infof(&myFlags, "Info test message %d of %02d", 2, 2)
	// Because of the config, none of these should output
	Debug(&myFlags, "Debug test message 1 of 03")
	Debug(&myFlags, "Debug test message 2 of 03")
	Debugf(&myFlags, "Debug test message %d of %02d", 3, 3)
	// get the size of the log
	lineCount := GetLogStats().lineCount
	if (lineCount - startingLineCount) != 9 {
		t.Errorf("Logit problem: Number of logged messages (%d) is incorrect", lineCount)
	}
	CloseLog()
}

func TestLogConfigChanges(t *testing.T) {
	const configFileName = "logtestcfg.json"
	// This string gets written to a temporary file for testing
	const jsonTest string = `
	{
		"SiteID": "BeyondAI",
		"SystemID": "local",
		"filename": "logTestFile.txt",
		"url": "",
		"stdout": true,
		"level": "INFO",
		"xFlags": [
			{ "pkg": "logit", "file": "log" , "flags": "0x01" }
		]
	}`

	const jsonTest2 string = `
	{
		"SiteID": "BeyondAI",
		"SystemID": "local",
		"filename": "logTestFile.txt", 
		"url": "",
		"stdout": true,
		"level": "DEBUG",
		"debugFlags": [
			{ "pkg": "all"}
		]
	}`

	err := writeConfigFile(configFileName, jsonTest)
	if err != nil {
		t.Errorf("Logit problem %s", err.Error())
	}
	defer removeConfigFile(configFileName)
	err = OpenLog(configFileName)
	if err != nil {
		t.Errorf("Logit problem %s", err.Error())
	}
	//
	// now test the results from loading the test configuration
	//
	var myFlags DFlags_t
	GetMyLogInfo(&myFlags)
	startingLineCount := GetLogStats().lineCount
	// These next 3 debug messages should NOT print to log
	Debug(&myFlags, "Debug test message 1 of 03")
	Debug(&myFlags, "Debug test message 2 of 03")
	Debugf(&myFlags, "Debug test message %d of %02d", 3, 3)
	// get the size of the log
	lineCount := GetLogStats().lineCount
	delta := lineCount - startingLineCount
	if delta != 0 {
		t.Errorf("Logit problem: Number of logged messages (%d) is incorrect", delta)
	}
	// remove then rewrite the config file again, to generate a modification event
	time.Sleep(1 * time.Second)
	Info(&myFlags, "Removing log config file")
	removeConfigFile(configFileName)
	Info(&myFlags, "Rewriting log config file")
	err = writeConfigFile(configFileName, jsonTest2)
	if err != nil {
		t.Errorf("Logit problem %s", err.Error())
	}
	// wait for more than the monitor loop of 5 seconds
	// to generate the event
	startingLineCount2 := GetLogStats().lineCount
	Info(&myFlags, "With next reconfig, no more DBGX messages")
	time.Sleep(22 * time.Second)
	Debug(&myFlags, "Sleep in main test loop just returned")
	// get the size of the log
	lineCount2 := GetLogStats().lineCount
	delta2 := lineCount2 - startingLineCount2
	if delta2 != 12 {
		t.Errorf("Logit problem: Number of logged messages (%d) is incorrect", delta2)
	}
	CloseLog()
}

//
// writeConfigFile
// Remove old file if present
// Write a new test config file of "contents"
//
func writeConfigFile(configFileName string, contents string) error {
	// remove old config file if present
	err := removeConfigFile(configFileName)
	if err != nil {
		return nil
	}
	// now create the test config file
	fh, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer fh.Close()
	// Now populate the test config file
	_, err = io.Copy(fh, strings.NewReader(contents))
	return err
}

//
// removeConfigFile
// If file exists, try to remove it. Return error if problem
//
func removeConfigFile(configFileName string) error {
	_, err := os.Stat(configFileName)
	if !os.IsNotExist(err) { // IsExist did not work properly
		err := os.Remove(configFileName)
		return err
	}
	return nil // no error
}
