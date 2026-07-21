package table

import (
	"os"

	"golang.org/x/term"
)

// TerminalWidth returns the current terminal's column width, or 0 when
// stdout isn't attached to a terminal (piped, redirected, tests, ...).
// Pass it to SetWidth/Row.SetWidth to size a table/row to the terminal; 0
// is a no-op there, so callers don't need to special-case non-tty output.
func TerminalWidth() int {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return 0
	}

	width, _, err := term.GetSize(fd)
	if err != nil || width <= 0 {
		return 0
	}

	return width
}
