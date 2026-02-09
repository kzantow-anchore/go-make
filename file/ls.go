package file

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/anchore/go-make/config"
	"github.com/anchore/go-make/lang"
	"github.com/anchore/go-make/log"
)

// Ls returns an ls-like string listing the contents of the given directory
func Ls(dir string) string {
	dir = lang.Return(filepath.Abs(dir))
	entries := lang.Return(os.ReadDir(dir))

	buf := bytes.Buffer{}
	for _, f := range entries {
		s := lang.Return(os.Stat(filepath.Join(dir, f.Name())))
		uid := -1
		gid := -1
		if stat, ok := s.Sys().(*syscall.Stat_t); ok {
			uid = int(stat.Uid)
			gid = int(stat.Gid)
		}
		_, _ = fmt.Fprintf(&buf, "%v %4v %4v %8v %v\n", s.Mode(), uid, gid, HumanizeBytes(s.Size()), f.Name())
	}
	return buf.String()
}

// LsAll returns an ls-like string listing the contents of the given directory and subdirectories
func LsAll(dir string) string {
	out := Ls(dir)
	entries := lang.Return(os.ReadDir(dir))
	for _, f := range entries {
		if f.IsDir() {
			subDir := filepath.Join(dir, f.Name())
			out += "\n\n -- " + subDir + ":"
			out += "\n" + LsAll(subDir)
		}
	}
	return out
}

// LogWorkdir logs an ls of the current working directory
func LogWorkdir() {
	if !config.Debug {
		return
	}
	cwd := Cwd()
	log.Info("CWD: %s", cwd)
	for _, line := range strings.Split(Ls(cwd), "\n") {
		log.Info(line)
	}
}

func HumanizeBytes[T int | int64](size T) any {
	units := ""
	value := size
	switch {
	case value < 1024:
	case value < 1024*1024:
		units = "KB"
		value /= 1024
	case value < 1024*1024*1024:
		units = "MB"
		value /= 1024 * 1024
	case value < 1024*1024*1024*1024:
		units = "GB"
		value /= 1024 * 1024 * 1024
	case value < 1024*1024*1024*1024*1024:
		units = "TB"
	}
	return fmt.Sprintf("%v%s", value, units)
}
