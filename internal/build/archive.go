// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package build

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractArchive(archive string, dest string) error {
	ar, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer ar.Close()

	switch {
	case strings.HasSuffix(archive, ".tar.gz"):
		return extractTarball(ar, dest)
	case strings.HasSuffix(archive, ".zip"):
		return extractZip(ar, dest)
	default:
		return fmt.Errorf("unhandled archive type %s", archive)
	}
}

// extractTarball unpacks a .tar.gz file.
func extractTarball(ar io.Reader, dest string) error {
	gzr, err := gzip.NewReader(ar)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		// Move to the next file header.
		header, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		// We only care about regular files, directory modes
		// and special file types are not supported.
		if header.Typeflag == tar.TypeReg {
			armode := header.FileInfo().Mode()
			err := extractFile(header.Name, armode, tr, dest)
			if err != nil {
				return fmt.Errorf("extract %s: %v", header.Name, err)
			}
		}
	}
}

// extractZip unpacks the given .zip file.
func extractZip(ar *os.File, dest string) error {
	info, err := ar.Stat()
	if err != nil {
		return err
	}
	zr, err := zip.NewReader(ar, info.Size())
	if err != nil {
		return err
	}

	for _, zf := range zr.File {
		if !zf.Mode().IsRegular() {
			continue
		}

		data, err := zf.Open()
		if err != nil {
			return err
		}
		err = extractFile(zf.Name, zf.Mode(), data, dest)
		data.Close()
		if err != nil {
			return fmt.Errorf("extract %s: %v", zf.Name, err)
		}
	}
	return nil
}

// extractFile extracts a single file from an archive.
func extractFile(arpath string, armode os.FileMode, data io.Reader, dest string) error {
	// Check that path is inside destination directory.
	target := filepath.Join(dest, filepath.FromSlash(arpath))
	if !strings.HasPrefix(target, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("path %q escapes archive destination", target)
	}

	// Remove the preivously-extracted file if it exists
	if err := os.RemoveAll(target); err != nil {
		return err
	}

	// Recreate the destination directory
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	// Copy file data.
	file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, armode)
	if err != nil {
		return err
	}
	if _, err = io.Copy(file, data); err != nil {
		file.Close()
		os.Remove(target)
		return err
	}
	return file.Close()
}
