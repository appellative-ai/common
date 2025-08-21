package iox

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
)

const (
	CwdVariable = "[cwd]"
	fileScheme  = "file"
	embeddedFS  = "file:///f:/"
	//embeddedFS = "f:/"
)

var (
	basePath = ""
	win      = false
	f        embed.FS
	mounted  = false
)

// init - set the base path and windows flag
func init() {
	cwd, err := os.Getwd()
	if err != nil {
		basePath = err.Error()
	}
	if os.IsPathSeparator(uint8(92)) {
		win = true
	}
	basePath = cwd
}

func DirFS(dir string) fs.FS {
	return os.DirFS(FileName(dir))
}

func Mount(fs embed.FS) {
	if mounted {
		fmt.Printf("error: file system is already mounted\n")
		return
	}
	f = fs
	mounted = true
}

// FileName - return the OS correct file name from a URI
func FileName(uri any) string {
	if uri == nil {
		return "error: URL is nil"
	}
	if s, ok := uri.(string); ok {
		if len(s) == 0 {
			return "error: URL is empty"
		}
		u, _ := url.Parse(s)
		return fileName(u)
	}
	if u, ok := uri.(*url.URL); ok {
		return fileName(u)
	}
	return fmt.Sprintf("error: invalid URL type: %v", reflect.TypeOf(uri))
}

func fileName(u *url.URL) string {
	if u == nil || u.Scheme != fileScheme {
		return fmt.Sprintf("error: scheme is invalid [%v]", u.Scheme)
	}
	name := basePath
	if u.Host == CwdVariable {
		name += u.Path
	} else {
		name = u.Path[1:]
	}
	if win {
		name = strings.ReplaceAll(name, "/", "\\")
	}
	return name
}

// ReadFile - read a file with a Status
func ReadFile(uri any) ([]byte, error) {
	rawUri := ""
	if s, ok := uri.(string); ok {
		rawUri = s
	} else {
		if u, ok2 := uri.(*url.URL); ok2 {
			rawUri = u.String()
		}
	}
	if strings.HasPrefix(rawUri, embeddedFS) {
		rawUri = rawUri[len(embeddedFS):]
		buf, err := f.ReadFile(rawUri)
		if err == nil {
			return buf, nil
		}
		return nil, err
	}
	buf, err := os.ReadFile(FileName(uri))
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// ReadFileWithEncoding - read a file with a possible encoding and a Status
func ReadFileWithEncoding(uri string, h http.Header) ([]byte, error) {
	buf, status := ReadFile(uri)
	if status != nil {
		return nil, status
	}
	return Decode(buf, h)
}
