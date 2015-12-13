// Package golocal provides a driver for the local filesystem for use with
// gofile.
package golocal

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Machiel/gofile"
)

type localDriver struct {
	rootDir string
}

func (c localDriver) absPath(path string) string {
	path, err := filepath.Abs(c.rootDir + "/" + path)
	if err != nil {
		log.Println(err)
	}
	return path
}

func (c localDriver) Has(path string) bool {
	if _, err := os.Stat(c.absPath(path)); os.IsNotExist(err) {
		return false
	}

	return true
}

func (c localDriver) Read(path string) (io.ReadCloser, error) {
	b, err := os.Open(c.absPath(path))

	if err != nil {
		return b, fmt.Errorf("golocal: Unable to open file '%s' (%s)", path, err)
	}

	return b, err
}

func (c localDriver) CreateDir(path string) error {
	err := os.Mkdir(c.absPath(path), 0755)

	if err != nil {
		return fmt.Errorf("golocal: Could not create directory '%s' (%s)", path, err)
	}

	return nil
}

func (c localDriver) DeleteDir(path string) error {
	err := os.RemoveAll(c.absPath(path))

	if err != nil {
		return fmt.Errorf("golocal: Could not delete directory '%s' (%s)", path, err)
	}

	return nil
}

func (c localDriver) List(path string) ([]gofile.File, error) {

	dirname := c.absPath(path)

	files, err := ioutil.ReadDir(dirname)

	if err != nil {
		return gofile.EmptyFileSet(), fmt.Errorf("golocal: Could not read directory '%s' (%s)", path, err)
	}

	foundFiles := make([]gofile.File, len(files))

	for i, file := range files {

		fullPath, _ := filepath.Abs(dirname + "/" + file.Name())

		foundFiles[i] = gofile.File{
			Path:  fullPath,
			IsDir: file.IsDir(),
		}
	}

	return foundFiles, nil
}

func (c localDriver) Write(path string, data io.Reader) error {
	contents, err := ioutil.ReadAll(data)

	if err != nil {
		return fmt.Errorf("golocal: Could not read data for writing to path %s (%s)", path, err)
	}

	err = ioutil.WriteFile(c.absPath(path), contents, 0644)

	if err != nil {
		return fmt.Errorf("golocal: Could not write file '%s' (%s)", path, err)
	}

	return nil
}

func (c localDriver) Copy(path, newPath string) error {
	s, err := os.Open(c.absPath(path))
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(c.absPath(newPath))
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}

func (c localDriver) Update(path string, data io.Reader) error {
	return c.Write(path, data)
}

func (c localDriver) Rename(path, newPath string) error {
	err := os.Rename(c.absPath(path), c.absPath(newPath))

	if err != nil {
		return fmt.Errorf("golocal: Could not rename file '%s' (%s)", path, err)
	}

	return nil
}

func (c localDriver) Delete(path string) error {
	err := os.Remove(c.absPath(path))

	if err != nil {
		return fmt.Errorf("golocal: Could not delete file '%s' (%s)", path, err)
	}

	return nil
}

func build(c map[string]string) (gofile.Driver, error) {

	var ok bool
	var rootDir string

	if rootDir, ok = c["rootDir"]; !ok {
		return nil, fmt.Errorf("golocal: 'rootDir' not specified")
	}

	return &localDriver{rootDir: rootDir}, nil
}

func init() {
	gofile.Register("local", build)
}
