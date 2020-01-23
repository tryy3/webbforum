package fileutils

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tryy3/webbforum/models"
)

// WriteTempFile attempts to write the content of reader into a temporary file
// returns a base64Hash of the data written, the size of the file and the location of the file
func WriteTempFile(reader io.Reader, maxFileSize int64, basePath string) (base64Hash string, size int64, path string, err error) {
	size = -1

	// create a temp file write
	tmpFileWriter, tmpFile, tmpDir, err := createTempFileWriter(basePath)
	if err != nil {
		return
	}

	// close file when function is done
	defer (func() { err = tmpFile.Close() })()

	// The amount of data read is limited to maxFileSizeBytes. At this point, if there is more data it will be truncated.
	limitedReader := io.LimitReader(reader, maxFileSize)

	// Hash the file data. The hash will be returned. The hash is useful as a
	// method of deduplicating files to save storage, as well as a way to conduct
	// integrity checks on the file data in the repository.
	hasher := sha256.New()
	teeReader := io.TeeReader(limitedReader, hasher)
	bytesWritten, err := io.Copy(tmpFileWriter, teeReader)
	if err != nil && err != io.EOF {
		return
	}

	err = tmpFileWriter.Flush()
	if err != nil {
		return
	}

	base64Hash = base64.RawURLEncoding.EncodeToString(hasher.Sum(nil)[:])
	size = bytesWritten
	path = tmpDir
	return
}

// createTempFileWriter will attempt to create a file writer of the temporary dir
func createTempFileWriter(basePath string) (*bufio.Writer, *os.File, string, error) {
	tmpDir, err := createTempDir(basePath)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create temp dir: %q", err)
	}

	writer, tmpFile, err := createFileWriter(tmpDir)
	if err != nil {
		return nil, nil, "", fmt.Errorf("failed to create file writer: %q", err)
	}

	return writer, tmpFile, tmpDir, nil
}

// createTempDir will attempt to create a base tmp folder and a tmp folder using ioutil.TempDir
func createTempDir(basePath string) (string, error) {
	if err := os.MkdirAll(basePath, 0770); err != nil {
		return "", fmt.Errorf("failed to create base temp dir: %v", err)
	}

	// create a temp folder inside the baseTmpDir
	tmpDir, err := ioutil.TempDir(basePath, "")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}

	return tmpDir, nil
}

// createFileWriter creates a buffered file writer with a new file
// The caller should flush the writer before closing the file.
// Returns the file handle as it needs to be closed when writing is complete
func createFileWriter(directory string) (*bufio.Writer, *os.File, error) {
	filePath := filepath.Join(string(directory), "content")
	file, err := os.Create(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create file: %v", err)
	}

	return bufio.NewWriter(file), file, nil
}

// MoveFile will attempt to move the file from temporary folder to real folder
func MoveFile(tmpDir, mediaFolder string, media models.File) (finalPath string, duplicate bool, err error) {
	// create media folder if it doesn't exists
	if err := os.MkdirAll(mediaFolder, 0770); err != nil {
		return "", false, fmt.Errorf("failed to create media folder: %v", err)
	}

	finalPath = filepath.Join(mediaFolder, media.Base64Hash)

	// check if file already exists
	var stat os.FileInfo
	if stat, err = os.Stat(finalPath); !os.IsNotExist(err) {
		duplicate = true

		if stat.Size() == media.FileSizeBytes {
			return finalPath, duplicate, nil
		}
		return "", duplicate, fmt.Errorf("downloaded file with hash collision but different file size (%v)", finalPath)
	}

	// move the file
	err = moveFile(filepath.Join(tmpDir, "content"), finalPath)
	if err != nil {
		return "", duplicate, fmt.Errorf("failed to move file to final destination (%v): %q", finalPath, err)
	}
	return finalPath, duplicate, nil
}

// moveFile will simply move src folder to dst folder
func moveFile(src, dst string) error {
	dstDir := filepath.Dir(dst)

	err := os.MkdirAll(dstDir, 0770)
	if err != nil {
		return fmt.Errorf("failed to make directory: %q", err)
	}

	err = os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("failed to move directory: %q", err)
	}

	return nil
}
