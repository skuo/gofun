// Package logit contains utility functions for logging for BeyondAI.
package logit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type logLevel_t int

const (
	FATAL logLevel_t = iota
	ERROR
	WARN
	INFO
	DEBUG
)

type logStats_t struct {
	fatalCount int32 // count of fatal messages output
	errorCount int32 // count of fatal messages output
	warnCount  int32 // count of warning messages output to log
	infoCount  int32 // count of informational messages output to log
	debugCount int32 // count of debug messages output to log
	lineCount  int32 // how many lines in the log file
	logSize    int64 // size of log in bytes
}

var logStats logStats_t

// logFlags_t holds all the flags loaded from the log config file
type logFlags_t struct {
	generation  int32            // the generation that this was reloaded
	siteID      string           // customer site name for log tracking
	sysID       string           // id used for url collections
	url         string           // url of log server
	useStdOut   bool             // use stdout for log messages
	logFileName string           // file name of log file
	logLevel    logLevel_t       // level = debug, info, warn, fatal
	debugAll    bool             // enable debug for everything
	dFlags      map[string]bool  // debug per package or package:file
	xFlags      map[string]int32 // granular debug for package:file
}

var allLogFlags *logFlags_t

type monitorFunc_t func(bool) (bool, int32) // your function pointer type
var monitorFunc monitorFunc_t

type DFlags_t struct {
	generation int32  // what generation of config file did the flags come from
	pkgName    string // name of this package
	fileName   string // name of this file
	dFlag      bool   // debug messages on/off
	xFlag      int32  // granular debug flags
}

var logConfigFileName string    // the configuration file name, set by 'OpenLog'
var logFileHandle *os.File      // handle to the log file itself
var logFileWriter *bufio.Writer // handle to the writer to the log file

const SHOWMONITOR int32 = 0x01 // show more debug msgs from inside the config monitor
var myFlags DFlags_t           // typical log flags just like all other files

// these are logs msgs that are placed here before the logger is fully established
type logDelayed_t struct {
	logLevel logLevel_t
	msg      string
}

var delayedLogs []logDelayed_t

/*
  OpenLog
  Open up the logconfig file to setup the logging
  either locally, or to a web service, or to stdout.
*/
func OpenLog(newLogConfigFileName string) error {
	var err error
	if len(newLogConfigFileName) == 0 {
		// use default log configuration file
		logConfigFileName = "logitcfg.json"
	} else {
		logConfigFileName = newLogConfigFileName
	}
	delayLog(INFO, fmt.Sprintf("Open logger from config file '%s'.", logConfigFileName))
	//
	var tFlags logFlags_t
	err = getConfig(logConfigFileName, &tFlags)
	if err != nil { // return if problemn
		return err
	}
	allLogFlags = &tFlags
	//
	// retrieve the package/file specific flags, just like any package
	//
	GetMyLogInfo(&myFlags)
	//
	// open the log file for writing
	//
	var fileErr error = nil
	logFileHandle = nil
	logFileWriter = nil
	if len(allLogFlags.logFileName) > 0 {
		logFileHandle, fileErr := os.Create(allLogFlags.logFileName)
		if fileErr == nil {
			logFileWriter = bufio.NewWriterSize(logFileHandle, 16384)
		} else { // then send everything to stdout
			allLogFlags.useStdOut = true
			delayLog(WARN, fmt.Sprintf("Failed to open output log: '%s'.", allLogFlags.logFileName))
			delayLog(WARN, "Logs will go to 'stdout'.")
		}
	}
	//
	// !!!The logger can be used at this point.
	// Flush any logs that were "delayed"
	//
	flushDelayLog()
	logTheFlags(allLogFlags)
	//
	// start the monitoring of log config changed
	//
	monitorFunc = watchLogConfig(logConfigFileName)
	monitorFunc(false) // get the baseline
	return fileErr
}

/*
  CloseLog
  Close the open log or tell the url
*/
func CloseLog() {
	logTheLogStats()
	if monitorFunc != nil {
		monitorFunc(true) // stop it
		monitorFunc = nil
	}
	Info(&myFlags, "Log file is being closed.")
	if logFileWriter != nil {
		tempFileWriter := logFileWriter
		logFileWriter = nil // turn this off while we wait for the final flush
		tempFileWriter.Flush()
		logFileHandle.Close()
		logFileHandle = nil
	}
	//fmt.Printf("Logger closed\n")
}

/*
  delayLog
  This is a special case logger for inside the log package.
  It stores logs until the logger is fully configured.
  When the logger is fully functional, these delayed msgs will
  be posted.
*/
func delayLog(level logLevel_t, msg string) {
	var dlog logDelayed_t
	dlog.logLevel = level
	dlog.msg = msg //fmt.Sprintf("Open logger from config file '%s'\n", logConfigFileName)
	delayedLogs = append(delayedLogs, dlog)
}

/*
  flushDelayLog
  This is a special case logger for inside the log package
  When the log is ready to write to, write all the "delayed" logs
  to the log file.
*/
func flushDelayLog() {
	for _, logItem := range delayedLogs {
		switch logItem.logLevel {
		case FATAL:
			Fatal(&myFlags, logItem.msg)
		case ERROR:
			Error(&myFlags, logItem.msg)
		case WARN:
			Warn(&myFlags, logItem.msg)
		case INFO:
			Info(&myFlags, logItem.msg)
		case DEBUG:
			Debug(&myFlags, logItem.msg)
		}
	}
	delayedLogs = delayedLogs[:0] // clear out the array
}

/*
  reloadConfig
  This function is called when the log configueration had been modified.
  This will cause all the calling log functions to get there flags reloaded.
  CAVEAT: This will reload the config from the save file as the original
  call to 'OpenLog', and dump to the same log and/or URL.
*/
func reloadConfig(newGeneration int32) {
	logTheLogStats()
	Infof(&myFlags, "log config file '%s' is being reloaded, gen %d.", logConfigFileName, newGeneration)
	savedURL := allLogFlags.url
	savedLogFileName := allLogFlags.logFileName
	// Now get the new stuff
	var tFlags logFlags_t
	err := getConfig(logConfigFileName, &tFlags)
	if err != nil {
		// if problem, throw new config away and log it
		flushDelayLog()
		Warnf(&myFlags, "log config file '%s' could not be reloaded, gen %d.", logConfigFileName, newGeneration)
		Warnf(&myFlags, "Continue to log with gen %d configuration.", tFlags.generation)
		return
	}
	//
	// no error was detected so setup flags with new confiuration
	// BUT transfer over the old url and logFileName from the old,
	// these are not mutable
	//
	tFlags.generation = newGeneration // set new generation
	tFlags.url = savedURL
	tFlags.logFileName = savedLogFileName
	allLogFlags = &tFlags   // this switches the world to the new config
	getLogDXFlags(&myFlags) // this is specific to just this file
	flushDelayLog()         // using the new flags
	logTheFlags(allLogFlags)
}

/*
  getConfig
  Get the logger configuration out of the local file
  The configuration is in json format
  Everything in this function writes to the calling parameter "tFlags"
  so if this configuration is flawed, it can simply be discarded.
*/
func getConfig(logConfigFileName string, tFlags *logFlags_t) error {
	type dflags_t struct {
		Pkg  string `json:"pkg"`
		File string `json:"file"`
	}

	type xflags_t struct {
		Pkg   string `json:"pkg"`
		File  string `json:"file"`
		Flags string `json:"flags"`
	}

	type configJason_t struct {
		SiteID     string     `json:"SiteID"`
		SysID      string     `json:"SystemID"`
		FileName   string     `json:"filename"`
		Url        string     `json:"url"`
		StdOut     bool       `json:"stdout"`
		Level      string     `json:"level"`
		Debugflags []dflags_t `json:"debugFlags"`
		XFlags     []xflags_t `json:"xFlags"`
	}
	//
	// setup defaults if the log configuration is not present
	//
	tFlags.generation = 0
	tFlags.siteID = "??"
	tFlags.url = "" // just log to local file
	tFlags.logFileName = "logfile.txt"
	tFlags.useStdOut = true
	tFlags.logLevel = INFO // info level
	//
	// Try to read the configuration file
	// TBD: if read error, write the production JSON and retry
	//
	raw, err_rf := ioutil.ReadFile(logConfigFileName)
	if err_rf != nil { // problem reading config
		delayLog(WARN, fmt.Sprintf("Error loading '%s', Error:'%s'.", logConfigFileName, err_rf.Error()))
		return err_rf
	}
	//
	// parse the json format into local jason structs
	// TBD: if JSON error, write defaults to file and try again
	//
	var res configJason_t
	err_um := json.Unmarshal(raw, &res)
	if err_um != nil {
		delayLog(WARN, fmt.Sprintf("JSON error: '%s'.", err_um.Error()))
		return err_um
	}
	//
	// get the general configuration flags out of the structs
	//
	tFlags.siteID = res.SiteID // label to tie the logs to a customer site
	tFlags.sysID = res.SysID   // id used for url collections
	delayLog(INFO, fmt.Sprintf("SiteID '%s', SystemID '%s'.", tFlags.siteID, tFlags.sysID))
	tFlags.url = res.Url
	if len(tFlags.url) != 0 {
		delayLog(INFO, fmt.Sprintf("Logs going to log server at '%s'.", tFlags.url))
	} else {
		delayLog(INFO, fmt.Sprintf("No log server url was specified in configuration."))
	}
	tFlags.useStdOut = res.StdOut       // true == output to stdout
	tFlags.logFileName = res.FileName   // file name of log file
	tFlags.logLevel = WARN              // default is log FATALs, ERRORs, and WARNs
	llstr := strings.ToUpper(res.Level) // level = trace, debug, info, warn, fatal
	switch llstr[0] {
	case 'W':
		tFlags.logLevel = WARN
	case 'I':
		tFlags.logLevel = INFO
	case 'D':
		tFlags.logLevel = DEBUG
	default:
		tFlags.logLevel = WARN
	}

	//
	// get the package and/or file specific flags out of json structs
	// populate dFlags maps with content
	//
	tFlags.dFlags = make(map[string]bool)
	tFlags.xFlags = make(map[string]int32)
	tFlags.debugAll = false // maybe set to true below
	if tFlags.logLevel >= DEBUG {
		//
		// If the "DEBUG" flag is not on, don't bother to
		// load any of the debug specifiers
		//
		for _, flag := range res.Debugflags {
			if len(flag.Pkg) == 0 { //no pkg name specified
				continue // skip entry
			}
			// turn on all if set
			if strings.ToLower(flag.Pkg) == "all" {
				tFlags.debugAll = true
				continue
			}
			// there is a pkg name specified
			if len(flag.File) == 0 { //no file name specified
				tFlags.dFlags[flag.Pkg] = true
			} else { // file name is specified
				tFlags.dFlags[flag.Pkg+":"+flag.File] = true
			}
		}
	}
	// populate xFlags maps with content
	for _, xflag := range res.XFlags {
		if len(xflag.Pkg) == 0 { //no pkg name specified
			continue // skip entry
		}
		// there is a pkg name specified
		value, _ := strconv.ParseInt(xflag.Flags, 0, 32)
		ivalue := int32(value)
		if len(xflag.File) == 0 { //no file name specified
			tFlags.xFlags[xflag.Pkg] = int32(ivalue)
		} else { // file name is specified
			tFlags.xFlags[xflag.Pkg+":"+xflag.File] = int32(value)
		}
	}
	return nil
}

/*
  logFlags
  log the flags so the verification can be done with support bundles
*/
func logTheFlags(flags *logFlags_t) {
	// log the debug flags for support bundle verification
	if flags.debugAll {
		Info(&myFlags, "All debug flags are enabled.")
	} else if len(flags.dFlags) == 0 {
		Info(&myFlags, "No debug flags are enabled.")
	} else {
		keys := make([]string, len(flags.dFlags)) // size it properly
		i := 0
		for key, _ := range flags.dFlags {
			keys[i] = key
			i++
		}
		sort.Strings(keys) // a sorted list of keys
		//
		dflagMsg := "Debug flags:"
		for _, key := range keys {
			dflagMsg += fmt.Sprintf("\n  flag:'%s', value:%t", key, flags.dFlags[key])
		}
		Info(&myFlags, dflagMsg)
	}
	//
	if len(flags.xFlags) > 0 {
		keys := make([]string, len(flags.xFlags)) // size it properly
		i := 0
		for key, _ := range flags.xFlags {
			keys[i] = key
			i++
		}
		sort.Strings(keys) // a sorted list of keys
		//
		xflagMsg := "Expert flags:"
		for _, key := range keys {
			xflagMsg += fmt.Sprintf("\n  xflag:'%s', value:'%x'", key, flags.xFlags[key])
		}
		Info(&myFlags, xflagMsg)
	}
}

/*
  logTheLogStats
  Log the collected logstats to this point
*/
func logTheLogStats() {
	Infof(&myFlags, "Fatal=%d, Error=%d, Warn=%d, Info=%d, Debug=%d",
		logStats.fatalCount, logStats.errorCount, logStats.warnCount, logStats.infoCount, logStats.debugCount)
	Infof(&myFlags, "line count = %d, log size = %d", logStats.lineCount, logStats.logSize)
}

/*
  GetMyLogInfo
  Lookup the callers info to get the relevant labels
  and levels into the dFlags struct
*/
func GetMyLogInfo(flags *DFlags_t) {
	pc, file, _, _ := runtime.Caller(1)
	_, fullFileName := path.Split(file)
	fileName := strings.Split(fullFileName, ".") // take off the ".go" extension
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)

	packageName := ""
	if parts[pl-2][0] == '(' {
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}
	flags.pkgName = packageName
	flags.fileName = fileName[0]
	getLogDXFlags(flags)
}

/*
	ifOldReloadDXFlags
	Check the generation of the flags. If older than the current
	generation of flags, reload the d and x flags
*/
func ifOldReloadDXFlags(flags *DFlags_t) {
	if flags.generation == allLogFlags.generation {
		return
	}
	getLogDXFlags(flags)
}

/*
	getLogDXFlags
	This uses the package and file names from the DFlags_t
	to lookup the associated flags from the config
*/
func getLogDXFlags(flags *DFlags_t) {
	flags.generation = allLogFlags.generation
	// lookup the xflag for the package and name
	packageName := flags.pkgName
	fileName := flags.fileName
	dp := allLogFlags.dFlags[packageName]
	df := allLogFlags.dFlags[packageName+":"+fileName]
	if allLogFlags.debugAll {
		df = true
	} else {
		if allLogFlags.logLevel >= DEBUG {
			if df && dp { // exclusive or, turn it off if both are on
				df = false
			} else if dp {
				df = true // then all files in the package are true
			}
		}
	}
	xp := allLogFlags.xFlags[packageName]
	xf := allLogFlags.xFlags[packageName+":"+fileName]
	xf = xp | xf // "or" the bitwise flags together
	//
	flags.dFlag = df
	flags.xFlag = xf
}

/*
  GetLogStats
  return the stats of the current log file
*/
func GetLogStats() logStats_t {
	return logStats
}

/*
  Fatal
  log the fatal messages, which are always enabled
*/
func Fatal(flags *DFlags_t, str string) {
	logStats.fatalCount += 1
	ifOldReloadDXFlags(flags)
	logMsg("FATAL[" + flags.pkgName + ":" + flags.fileName + "] " + str)
}

/*
  Fatalf
  Build the message and log the fatal messages
*/
func Fatalf(flags *DFlags_t, str string, args ...interface{}) {
	logStats.fatalCount += 1
	ifOldReloadDXFlags(flags)
	message := fmt.Sprintf(str, args...)
	logMsg("FATAL[" + flags.pkgName + ":" + flags.fileName + "] " + message)

}

/*
  Error
  log the error messages if enabled
*/
func Error(flags *DFlags_t, str string) {
	ifOldReloadDXFlags(flags)
	if allLogFlags.logLevel >= ERROR {
		logStats.errorCount += 1
		logMsg("ERR[" + flags.pkgName + ":" + flags.fileName + "] " + str)
	}
}

/*
  Warnf
  Build the message and log the error messages if enabled
*/
func Errorf(flags *DFlags_t, str string, args ...interface{}) {
	ifOldReloadDXFlags(flags)
	if allLogFlags.logLevel >= ERROR {
		logStats.errorCount += 1
		message := "ERR[" + flags.pkgName + ":" + flags.fileName + "] "
		message += fmt.Sprintf(str, args...)
		logMsg(message)
	}
}

/*
  Warn
  log the warning messages if enabled
*/
func Warn(flags *DFlags_t, str string) {
	ifOldReloadDXFlags(flags)
	if allLogFlags.logLevel >= WARN {
		logStats.warnCount += 1
		logMsg("WARN[" + flags.pkgName + ":" + flags.fileName + "] " + str)
	}
}

/*
  Warnf
  Build the message and log the warning messages if enabled
*/
func Warnf(flags *DFlags_t, str string, args ...interface{}) {
	ifOldReloadDXFlags(flags)
	if allLogFlags.logLevel >= WARN {
		logStats.warnCount += 1
		message := "WARN[" + flags.pkgName + ":" + flags.fileName + "] "
		message += fmt.Sprintf(str, args...)
		logMsg(message)
	}
}

/*
  Info
  log the info messages if enabled
*/
func Info(flags *DFlags_t, str string) {
	ifOldReloadDXFlags(flags)
	if allLogFlags.logLevel >= INFO {
		logStats.infoCount += 1
		logMsg("INFO[" + flags.pkgName + ":" + flags.fileName + "] " + str)
	}
}

/*
  Infof
  Build the message and log the info messages if enabled
*/
func Infof(flags *DFlags_t, str string, args ...interface{}) {
	ifOldReloadDXFlags(flags)
	if allLogFlags.logLevel >= INFO {
		logStats.infoCount += 1
		message := "INFO[" + flags.pkgName + ":" + flags.fileName + "] "
		message += fmt.Sprintf(str, args...)
		logMsg(message)
	}
}

/*
  Debug
  log the debug message if enabled
  This allows individual packages or files to log debug statements.
*/
func Debug(flags *DFlags_t, str string) {
	ifOldReloadDXFlags(flags)
	if flags.dFlag {
		logStats.debugCount += 1
		logMsg("DBUG[" + flags.pkgName + ":" + flags.fileName + "] " + str)
	}
}

/*
  Debugx
  log the debug message if enabled by an xflag
  This allows individual packages or files to log debug statements depending
  on very specific criteria established by the developer.
*/
func Debugx(xflag int32, flags *DFlags_t, str string) {
	ifOldReloadDXFlags(flags)
	if (flags.xFlag & xflag) != 0 {
		logStats.debugCount += 1
		logMsg("DBGX[" + flags.pkgName + ":" + flags.fileName + "] " + str)
	}
}

/*
  Debugf
  Build the message and log the debug message if enabled
  This allows individual packages or files to log debug statements.
*/
func Debugf(flags *DFlags_t, str string, args ...interface{}) {
	ifOldReloadDXFlags(flags)
	if flags.dFlag {
		logStats.debugCount += 1
		message := "DBUG[" + flags.pkgName + ":" + flags.fileName + "] "
		message += fmt.Sprintf(str, args...)
		logMsg(message)
	}
}

/*
  Debugfx
  Build and log the debug message if enabled by an xflag
  This allows individual packages or files to log debug statements depending
  on very specific criteria established by the developer.
*/
func Debugfx(xflag int32, flags *DFlags_t, str string, args ...interface{}) {
	ifOldReloadDXFlags(flags)
	if (flags.xFlag & xflag) != 0 {
		logStats.debugCount += 1
		message := "DBGX[" + flags.pkgName + ":" + flags.fileName + "] "
		message += fmt.Sprintf(str, args...)
		logMsg(message)
	}
}

/*
  logMsg
  The log message is almost built except for the time fields
  Send to the log file, and/or stdout, and/or the log server.
*/
func logMsg(msg string) {
	t := time.Now()
	msgt := t.Format(time.RFC3339) + " " + msg
	//  update stats
	logStats.lineCount += 1
	logStats.logSize += int64(len(msgt))
	// lock it down
	mutex := &sync.Mutex{} // Protect the channels at an EOL boundry
	mutex.Lock()
	defer mutex.Unlock()
	// write to standard out
	if allLogFlags.useStdOut { // write to stdout
		println(msgt)
	}
	// write to local file
	if logFileWriter != nil {
		logFileWriter.WriteString(msgt + "\n")
	}
	// wrtie to logserver
	if len(allLogFlags.url) > 0 {
		msgi := "{" + allLogFlags.sysID + "} " + msgt + "\n"
		logFileWriter.WriteString(msgi)
	}
	//
	// check if the log configuration file has changed
	// If so reload and set the next generation
	//
	if monitorFunc != nil {
		changed, newGeneration := monitorFunc(false)
		if changed {
			reloadConfig(newGeneration)
		}
	}
}

/*
	watchLogConfig
	monitor the log configuration file for changes
*/
func watchLogConfig(logConfigFileName string) func(bool) (bool, int32) {
	//GetMyLogInfo(&lwFlags)
	// get the baseline of this file
	baseLine, err := os.Stat(logConfigFileName)
	if err != nil {
		Warnf(&myFlags, "log config file '%s' cannot be monitored, error: %s", logConfigFileName, err)
		return nil
	}
	baseTime := baseLine.ModTime()
	Infof(&myFlags, "log config file '%s' is monitored", logConfigFileName)
	//
	ticker := time.NewTicker(time.Second * 5)
	stopChan := make(chan bool)              // unbuffered going into go function
	doneOrChangeChan := make(chan string, 2) // buffered coming from go function
	changed := false
	var generation int32
	monitorFunc := func() {
	DONE:
		for {
			select {
			case <-ticker.C:
				Debugx(SHOWMONITOR, &myFlags, "log reconfiguration check.")
				newStat, err := os.Stat(logConfigFileName)
				if err != nil {
					continue
				}
				modTime := newStat.ModTime()
				if baseTime != modTime {
					Debugfx(SHOWMONITOR, &myFlags, "log config file '%s' was modified.", logConfigFileName)
					baseTime = modTime
					changed = true
					//doneOrChangeChan <- "change"
				}
			case <-stopChan:
				Debugx(SHOWMONITOR, &myFlags, "log reconfiguration check terminated.")
				ticker.Stop() // remove ticker channel
				break DONE    // very obscure way to exit nested select statement
			}
		}
		Debugx(SHOWMONITOR, &myFlags, "Exiting log reconfiguration monitor.")
		doneOrChangeChan <- "done" // signal main loop to continue
	}

	monitorConfig := func(stopFlag bool) (bool, int32) {
		if stopFlag {
			stopChan <- true // tell the goroutine to shutdown
		DONE:
			for {
				select { // a blocking select statement
				case msg := <-doneOrChangeChan:
					if msg == "done" {
						Infof(&myFlags, "log config file '%s' is no longer monitored.", logConfigFileName)
						break DONE // exit for loop
					}
				}
			}
		} else {
			// check for change
			/*
				select { // a non-blocking select statement
				case msg := <-doneOrChangeChan:
					if msg == "change" {
						Infof(&lwFlags, "log config file '%s' was modified", logConfigFileName)
						return true // tell caller of change
					} else {
						Warnf(&lwFlags, "unexpected msg '%s' from log config monitor", msg)
					}
				default:
				}
			*/
			if changed {
				changed = false // reset the flag for next time
				generation += 1
				Infof(&myFlags, "log config file '%s' was modified, gen %d.", logConfigFileName, generation)
				return true, generation // file changed flag, and new generation
			}
		}
		return false, generation // file did not change and same old generation
	}
	go monitorFunc() // start the background loop
	return monitorConfig
}
