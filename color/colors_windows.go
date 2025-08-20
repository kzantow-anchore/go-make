package color

import (
	"os"

	"golang.org/x/sys/windows"
)

func init() {
	fd := windows.Handle(os.Stderr.Fd()) // os.Stdout?
	var originalMode uint32
	windows.GetConsoleMode(fd, &originalMode)
	windows.SetConsoleMode(fd, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
