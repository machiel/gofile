package golocal

import "testing"

func TestAbsPath(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"folder", "/tmp/folder"},
		{"/folder", "/tmp/folder"},
		{"/folder/", "/tmp/folder"},
		{"folder/file", "/tmp/folder/file"},
		{"/folder/file", "/tmp/folder/file"},
		{"/folder/file/", "/tmp/folder/file"},
	}

	rootDirs := []string{"/tmp", "/tmp/"}

	for _, rootDir := range rootDirs {
		driver := localDriver{
			rootDir: rootDir,
		}

		for _, test := range tests {
			actual := driver.absPath(test.in)

			if actual != test.out {
				t.Errorf("Expected '%s', got '%s'", test.out, actual)
			}
		}
	}
}

func TestBuild(t *testing.T) {
	_, err := build(map[string]string{})

	if err == nil {
		t.Error("Expected error for not passing rootDir")
	}

	_, err = build(map[string]string{"rootDir": "/tmp"})

	if err != nil {
		t.Error("Error returned while creating Driver")
	}
}
