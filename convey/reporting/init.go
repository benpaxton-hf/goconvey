package reporting

import (
	"bytes"
	"os/exec"
	"runtime"
)

func init() {
	redColor = tput("setaf", "1")
	greenColor = tput("setaf", "2")
	yellowColor = tput("setaf", "3")
	resetColor = tput("sgr0")

	if runtime.GOOS == "windows" {
		success, failure, error_ = dotSuccess, dotFailure, dotError
	}
}

func tput(args ...string) string {
	buf := &bytes.Buffer{}

	cmd := exec.Command("tput", args...)
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		// If there was an error, ignore it and proceed in monochrome
		return ""
	}

	return buf.String()
}

func BuildJsonReporter() Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(),
		NewJsonReporter(out))
}
func BuildDotReporter() Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(),
		NewDotReporter(out),
		NewProblemReporter(out),
		consoleStatistics)
}
func BuildStoryReporter() Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(),
		NewStoryReporter(out),
		NewProblemReporter(out),
		consoleStatistics)
}
func BuildSilentReporter() Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(),
		NewSilentProblemReporter(out))
}

var (
	newline         = "\n"
	success         = "✔"
	failure         = "✘"
	error_          = "🔥"
	skip            = "⚠"
	dotSuccess      = "."
	dotFailure      = "x"
	dotError        = "E"
	dotSkip         = "S"
	errorTemplate   = "* %s \nLine %d: - %v \n%s\n"
	failureTemplate = "* %s \nLine %d:\n%s\n%s\n"
)

var (
	greenColor, yellowColor, redColor, resetColor string
)

var consoleStatistics = NewStatisticsReporter(NewPrinter(NewConsole()))

func SuppressConsoleStatistics() { consoleStatistics.Suppress() }
func PrintConsoleStatistics()    { consoleStatistics.PrintSummary() }

// QuiteMode disables all console output symbols. This is only meant to be used
// for tests that are internal to goconvey where the output is distracting or
// otherwise not needed in the test output.
func QuietMode() {
	success, failure, error_, skip, dotSuccess, dotFailure, dotError, dotSkip = "", "", "", "", "", "", "", ""
}

// This interface allows us to pass the *testing.T struct
// throughout the internals of this tool without ever
// having to import the "testing" package.
type T interface {
	Fail()
}
