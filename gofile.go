// Package gofile provides a generic interface around filesystems.
//
// The gofile must be used in conjunction with a filesystem driver.
package gofile

import (
	"fmt"
	"io"
	"sync"
)

var (
	builders = make(map[string]Builder)
	mu       sync.Mutex
)

// File represents a file in the filesystem
type File struct {
	Path  string
	IsDir bool
}

// EmptyFileSet returns an empty slice of files
func EmptyFileSet() []File {
	return make([]File, 0)
}

// Reader provides all the operations on the filesystem that read
type Reader interface {
	Has(path string) bool
	Read(path string) (io.ReadCloser, error)
	List(path string) ([]File, error)
}

// Writer provides all the operations on the filesystem that changes the
// filesystem
type Writer interface {
	Write(path string, data io.Reader) error
	CreateDir(path string) error
	DeleteDir(path string) error
	Update(path string, data io.Reader) error
	Rename(path, newPath string) error
	Copy(path, target string) error
	Delete(path string) error
}

// Driver is the interface that must be implemented by the driver
type Driver interface {
	Reader
	Writer
}

// Builder is the builds the driver based on the configuration passed to the
// builder
type Builder func(config map[string]string) (Driver, error)

// New creates a new driver based on the driver selected and the configuration
// passed
func New(name string, config map[string]string) (Driver, error) {

	mu.Lock()
	builder, ok := builders[name]
	mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("gofile: Unknown implementation '%s'", name)
	}

	return builder(config)
}

// Register allows a driver to register itself making it available by the given
// name. If register is called multiple times with the same name, it will panic.
func Register(name string, builder Builder) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := builders[name]; ok {
		panic("Tried to register a driver under the same name twice:" + name)
	}

	builders[name] = builder
}
