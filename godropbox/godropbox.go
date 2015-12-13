// Package godropbox provides a driver for use with gofile, which allows you to
// use dropbox as a filesystem.
package godropbox

import (
	"fmt"
	"io"

	"github.com/Machiel/gofile"
	"github.com/stacktic/dropbox"
)

type dropboxDriver struct {
	db *dropbox.Dropbox
}

type nopCloser struct {
	io.Reader
}

func (nc nopCloser) Close() error {
	return nil
}

func (d dropboxDriver) Read(path string) (io.ReadCloser, error) {
	content, _, err := d.db.Download(path, "", 0)

	if err != nil {
		return nil, fmt.Errorf("godropbox: Could not read file on path '%s' (%s)", path, err)
	}

	return content, nil
}

func (d dropboxDriver) Has(path string) bool {
	_, err := d.db.Metadata(path, false, false, "", "", 1)

	if err != nil {
		return false
	}

	return true
}

func (d dropboxDriver) write(path string, data io.Reader, overwrite bool) error {
	content := nopCloser{data}
	_, err := d.db.FilesPut(content, 0, path, overwrite, "")
	return err
}

func (d dropboxDriver) delete(path string) error {
	_, err := d.db.Delete(path)
	return err
}

func (d dropboxDriver) Write(path string, data io.Reader) error {
	err := d.write(path, data, false)

	if err != nil {
		return fmt.Errorf("godropbox: Could not create file '%s' (%s)", path, err)
	}

	return nil
}

func (d dropboxDriver) Update(path string, data io.Reader) error {
	err := d.write(path, data, true)

	if err != nil {
		return fmt.Errorf("godropbox: Could not update file '%s' (%s)", path, err)
	}

	return nil
}

func (d dropboxDriver) List(path string) ([]gofile.File, error) {
	files, err := d.db.Metadata(path, true, false, "", "", 25000)

	if err != nil {
		/** @TODO Handle this properly **/
		return gofile.EmptyFileSet(), err
	}

	if !files.IsDir {
		return gofile.EmptyFileSet(), fmt.Errorf("godropbox: Cannot list files of non directory")
	}

	results := len(files.Contents)
	foundFiles := make([]gofile.File, results)

	for i, file := range files.Contents {
		foundFiles[i] = gofile.File{
			Path:  file.Path,
			IsDir: file.IsDir,
		}
	}

	return foundFiles, nil
}

func (d dropboxDriver) CreateDir(path string) error {
	if _, err := d.db.CreateFolder(path); err != nil {
		return fmt.Errorf("godropbox: Could not create folder '%s' (%s)", path, err)
	}

	return nil
}

func (d dropboxDriver) DeleteDir(path string) error {
	if err := d.delete(path); err != nil {
		return fmt.Errorf("godropbox: Could not delete folder '%s' (%s)", path, err)
	}

	return nil
}

func (d dropboxDriver) Rename(path, newPath string) error {
	if _, err := d.db.Move(path, newPath); err != nil {
		return fmt.Errorf("godropbox: Could not rename file '%s' (%s)", path, err)
	}

	return nil
}

func (d dropboxDriver) Copy(path, newPath string) error {
	if _, err := d.db.Copy(path, newPath, false); err != nil {
		return fmt.Errorf(
			"godropbox: Could not copy file '%s' to '%s' (%s)",
			path,
			newPath,
			err,
		)
	}

	return nil
}

func (d dropboxDriver) Delete(path string) error {
	if err := d.delete(path); err != nil {
		return fmt.Errorf("godropbox: Could not delete file '%s' (%s)", path, err)
	}

	return nil
}

func build(config map[string]string) (gofile.Driver, error) {

	var clientID, clientSecret, token string
	var ok bool

	if clientID, ok = config["client_id"]; !ok {
		return nil, fmt.Errorf("godropbox: Client ID not specified")
	}

	if clientSecret, ok = config["client_secret"]; !ok {
		return nil, fmt.Errorf("godropbox: Client Secret not specified")
	}

	if token, ok = config["token"]; !ok {
		return nil, fmt.Errorf("godropbox: Token not specified")
	}

	db := dropbox.NewDropbox()
	db.SetAppInfo(clientID, clientSecret)
	db.SetAccessToken(token)

	return dropboxDriver{db: db}, nil
}

func init() {
	gofile.Register("dropbox", build)
}
