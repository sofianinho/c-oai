package utils

import (
	"archive/zip"
	"os"
	"path/filepath"
	"io"
	"fmt"
)

//Unzip unzips a zip archive given in src into the location dst
func Unzip(src, dst string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return fmt.Errorf("error zip reader: %s", err)
    }
    defer func() error{
        if err := r.Close(); err != nil {
            return fmt.Errorf("error zip close: %s", err)
		}
		return nil
    }()

    os.MkdirAll(dst, 0755)

    // Closure to address file descriptors issue with all the deferred .Close() methods
    extractAndWriteFile := func(f *zip.File) error {
        rc, err := f.Open()
        if err != nil {
            return fmt.Errorf("error zip reader: %s", err)
        }
        defer func() error{
            if err := rc.Close(); err != nil {
                return fmt.Errorf("error zip close: %s", err)
			}
			return nil
        }()

        path := filepath.Join(dst, f.Name)

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, f.Mode())
        } else {
            os.MkdirAll(filepath.Dir(path), f.Mode())
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return fmt.Errorf("failed to create directory %s: %s", path, err)
            }
            defer func() error {
                if err := f.Close(); err != nil {
                    return fmt.Errorf("error close dir: %s", err)
				}
				return nil
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return fmt.Errorf("failed to write: %s", err)
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return fmt.Errorf("failed to extract: %s", err)
        }
    }

    return nil
}