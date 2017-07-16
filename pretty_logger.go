package main

import "time"
import "os"
import "fmt"
import "github.com/briandowns/spinner"
import "github.com/mgutz/ansi"

type PrettyLogger struct {
	spinner *spinner.Spinner
}

func NewPrettyLogger() *PrettyLogger {
	return &PrettyLogger{
		spinner: spinner.New(spinner.CharSets[14], 100*time.Millisecond),
	}
}

func (logger *PrettyLogger) Info(message string) {
	fmt.Printf("%s›%s %s\n", ansi.ColorCode("yellow"), ansi.ColorCode("reset"), message)
}

func (logger *PrettyLogger) Success(message string) {
	fmt.Printf("%s✔%s %s\n", ansi.ColorCode("green"), ansi.ColorCode("reset"), message)
}

func (logger *PrettyLogger) Error(message string) {
	fmt.Fprintf(os.Stderr, "%s✖%s %s\n", ansi.ColorCode("red"), ansi.ColorCode("reset"), message)
	os.Exit(1)
}

func (logger *PrettyLogger) Loading() {
	logger.spinner.Start()
}

func (logger *PrettyLogger) Done() {
	logger.spinner.Stop()
}
