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

// Driver is the interface that must be implemented by the driver
type Driver interface {
	ContentChecker
	Reader
	Lister
	DirectoryCreator
	DirectoryDeleter
	Writer
	Updater
	Copier
	Deleter
	Renamer
}

// File represents a file in the filesystem
type File struct {
	Path  string
	IsDir bool
}

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

// Builder is the builds the driver based on the configuration passed to the
// builder
type Builder func(config map[string]string) (Driver, error)

// Reader provides all the operations on the filesystem that read
//
// Read returns a io.ReadCloser containing the contents of the requested file
type Reader interface {
	Read(path string) (io.ReadCloser, error)
}

// ContentChecker wraps the Contains method.
//
// Contains checks whether or not a file is available for a given filesystem
type ContentChecker interface {
	Contains(path string) bool
}

// Lister is the interface that wraps the List method
//
// Given a path it returns a list of all the objects available for that path
type Lister interface {
	List(path string) ([]File, error)
}

// Writer is the interface that wraps the Write method
//
// Write allows you to write to the file system, giving a path and an io.Reader.
type Writer interface {
	Write(path string, data io.Reader) error
}

// Updater is the interface that wraps the Update method
//
// Update takes a path that requires an update, and it passes an io.Reader with
// the content that needs to be updated.
type Updater interface {
	Update(path string, data io.Reader) error
}

// Copier is the interface that wraps the Copy method
//
// Copy copies the contents of one file to the target
type Copier interface {
	Copy(path, target string) error
}

// Deleter is the interface that wraps the Delete method.
//
// Delete removes a file based on a given path
type Deleter interface {
	Delete(path string) error
}

// Renamer is the interface that wraps the Rename method.
//
// Rename moves one file to another using the path and the newPath arguments.
type Renamer interface {
	Rename(path, newPath string) error
}

// DirectoryDeleter is the interface that wraps the DeleteDir method.
//
// DeleteDir removes an entiry directory of files
type DirectoryDeleter interface {
	DeleteDir(path string) error
}

// DirectoryCreator is the interface that wraps the CreateDir method.
//
// CreateDir creates a directory
type DirectoryCreator interface {
	CreateDir(path string) error
}

// EmptyFileSet returns an empty slice of files
func EmptyFileSet() []File {
	return make([]File, 0)
}
