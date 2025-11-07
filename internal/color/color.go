package color

import (
	"github.com/fatih/color"
)

var (
	Red    = color.New(color.FgRed)
	Green  = color.New(color.FgGreen)
	Yellow = color.New(color.FgYellow)
	Blue   = color.New(color.FgBlue)
)

// PrintError prints an error message in red
func PrintError(format string, a ...interface{}) {
	Red.Printf("Error: "+format+"\n", a...)
}

// PrintWarning prints a warning message in yellow
func PrintWarning(format string, a ...interface{}) {
	Yellow.Printf("Warning: "+format+"\n", a...)
}

// PrintSuccess prints a success message in green
func PrintSuccess(format string, a ...interface{}) {
	Green.Printf(format+"\n", a...)
}

func PrintInfo(format string, a ...interface{}) {
	Blue.Printf(format+"\n", a...)
}
