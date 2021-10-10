package output

import (
	"fmt"
	"io"
	"os"
)

var (
	// Stdout points to the output buffer to send screen output
	Stdout io.Writer = os.Stdout
	// Stderr points to the output buffer to send errors to the screen
	Stderr io.Writer = os.Stderr
)

// Printf is just like fmt.Printf except that it send the output to Stdout. It
// is equal to fmt.Fprintf(util.Stdout, format, args)
func Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(Stdout, format, args...)
}

// Eprintf prints the errors to the output buffer Stderr. It is equal to
// fmt.Fprintf(util.Stderr, format, args)
func Eprintf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(Stderr, format, args...)
}
