//go:build !finder

package viper

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Search all configPaths for any config file.
// Returns the first path that exists (and is a config file).
func (v *Viper) findConfigFile() (string, error) {
	v.logger.Info("searching for config in paths", "paths", v.configPaths)

	for _, cp := range v.configPaths {
		file := v.searchInPath(cp)
		if file != "" {
			return file, nil
		}
	}
	return "", ConfigFileNotFoundError{v.configName, fmt.Sprintf("%s", v.configPaths)}
}

func (v *Viper) searchInPath(in string) (filename string) {
	v.logger.Debug("searching for config in path", "path", in)
	for _, ext := range SupportedExts {
		v.logger.Debug("checking if file exists", "file", filepath.Join(in, v.configName+"."+ext))
		if b, _ := exists(v.fs, filepath.Join(in, v.configName+"."+ext)); b {
			v.logger.Debug("found file", "file", filepath.Join(in, v.configName+"."+ext))
			return filepath.Join(in, v.configName+"."+ext)
		}
	}

	if v.configType != "" {
		if b, _ := exists(v.fs, filepath.Join(in, v.configName)); b {
			return filepath.Join(in, v.configName)
		}
	}

	return ""
}

// exists checks if file exists.
func exists(f fs.FS, path string) (bool, error) {
	stat, err := fs.Stat(f, path)
	if err == nil {
		return !stat.IsDir(), nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

type osFS struct{}

func (o osFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func (o osFS) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func (o osFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (o osFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (o osFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
