package core

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/Leakageonthelamp/go-leakage-core/utils"
	ziplib "github.com/mholt/archiver/v4"
	"github.com/pkg/errors"
)

type IArchiver interface {
	FromURLs(fileName string, urls []string, options *ArchiverOptions) ([]byte, IError)
	FromBytes(fileName string, body []ArchiveByteBody, options *ArchiverOptions) ([]byte, IError)
}

type ArchiverOptions struct {
}

type ArchiveByteBody struct {
	File []byte
	Name string
}
type archiver struct {
	ctx IContext
}

var archiverFailDownLoad = Error{
	Status:  http.StatusInternalServerError,
	Code:    "ARCHIVER_DOWNLOAD_ERROR",
	Message: "failed to download file"}
var archiverFailCreateFile = Error{
	Status:  http.StatusInternalServerError,
	Code:    "ARCHIVER_CREATE_FILE_ERROR",
	Message: "failed to create file",
}

var archiverFailWriteFile = Error{
	Status:  http.StatusInternalServerError,
	Code:    "ARCHIVER_WRITE_FILE_ERROR",
	Message: "failed to write downloaded data to file",
}
var archiverFailZip = Error{
	Status:  http.StatusInternalServerError,
	Code:    "ARCHIVER_ZIP_ERROR",
	Message: "failed to zip file",
}

var archiverReadFileError = Error{
	Status:  http.StatusInternalServerError,
	Code:    "ARCHIVER_READ_FILE_ERROR",
	Message: "failed to read file",
}

func (s archiver) FromURLs(fileName string, urls []string, options *ArchiverOptions) ([]byte, IError) {
	tmpDir := "tmp" + utils.NewSha256(utils.GetCurrentDateTime().String())
	os.MkdirAll(tmpDir, os.ModePerm)

	// Cleanup the temporary folder
	defer os.RemoveAll(tmpDir)

	var wg sync.WaitGroup
	errs := make(chan IError, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			filename := filepath.Base(url)
			dest := filepath.Join(tmpDir, filename)

			err := s.downloadFile(url, dest)
			errs <- err
		}(url)
	}

	wg.Wait()
	close(errs)

	for ierr := range errs {
		if ierr != nil {
			return nil, s.ctx.NewError(ierr, ierr)
		}
	}

	return s.createZipFile(fileName, tmpDir)
}

func (s archiver) FromBytes(fileName string, body []ArchiveByteBody, options *ArchiverOptions) ([]byte, IError) {
	tmpDir := "tmp" + utils.NewSha256(utils.GetCurrentDateTime().String())
	os.MkdirAll(tmpDir, os.ModePerm)

	// Cleanup the temporary folder
	defer os.RemoveAll(tmpDir)

	for _, file := range body {
		dest := filepath.Join(tmpDir, file.Name)
		out, err := os.Create(dest)
		if err != nil {
			return nil, s.ctx.NewError(errors.Wrap(err, "failed to create destination file"), archiverFailCreateFile)
		}

		reader := bytes.NewReader(file.File)
		_, err = io.Copy(out, reader)
		if err != nil {
			return nil, s.ctx.NewError(errors.Wrap(err, "failed to write downloaded data to file"), archiverFailWriteFile)
		}

		out.Close()
	}

	return s.createZipFile(fileName, tmpDir)
}

func NewArchiver(ctx IContext) IArchiver {
	return &archiver{
		ctx: ctx,
	}
}

func (s archiver) createZipFile(fileName string, fromDir string) ([]byte, IError) {
	zipPath := fromDir + "/" + fileName + ".zip"
	// Get the list of files to archive
	files, err := ziplib.FilesFromDisk(nil, map[string]string{
		fromDir: fileName,
	})
	if err != nil {
		return nil, s.ctx.NewError(err, archiverFailZip)
	}

	// Create the output file
	out, err := os.Create(zipPath)
	if err != nil {
		return nil, s.ctx.NewError(err, archiverFailZip)
	}
	defer out.Close()

	// Archive the files
	format := ziplib.CompressedArchive{
		Archival: ziplib.Zip{},
	}
	err = format.Archive(context.Background(), out, files)
	if err != nil {
		return nil, s.ctx.NewError(err, archiverFailZip)
	}

	file, err := os.Open(zipPath)
	if err != nil {
		return nil, s.ctx.NewError(err, archiverReadFileError)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, s.ctx.NewError(err, archiverFailZip)
	}

	return data, nil
}

func (s archiver) downloadFile(url string, dest string) IError {
	resp, err := http.Get(url)
	if err != nil {
		return s.ctx.NewError(errors.Wrap(err, "failed to download file"), archiverFailDownLoad)
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return s.ctx.NewError(errors.Wrap(err, "failed to create destination file"), archiverFailCreateFile)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return s.ctx.NewError(errors.Wrap(err, "failed to write downloaded data to file"), archiverFailWriteFile)
	}

	return nil
}
