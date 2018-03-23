package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/input"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

type config struct {
	SourcePath  string `env:"source_path,file"`
	Destination string `env:"destination"`
}

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	stepconf.Print(cfg)

	err := ensureZIPExtension(cfg.SourcePath, cfg.Destination)

	if err != nil {
		failf("Issue with compress: %s", err)
		return
	}
}

func ensureZIPExtension(sourcePath string, destination string) (err error) {
	if !strings.HasPrefix(destination, ".zip") {
		destination += ".zip"
	}

	dirOftargetPath := filepath.Dir(destination)

	err = os.MkdirAll(dirOftargetPath, 0755)
	if err != nil {
		log.Errorf("Failed to create directory, error: %s", err)
	}

	checkAlreadyExist(destination)

	zipfile, err := os.Create(destination)
	if err != nil {
		return err
	}

	defer func() {
		if err = zipfile.Close(); err != nil {
			log.Errorf("%s", err)
		}
	}()

	zipWriter := zip.NewWriter(zipfile)

	defer func() {
		if cerr := zipWriter.Close(); err == nil {
			err = cerr
		}
	}()

	info, err := os.Lstat(sourcePath)
	if err != nil {
		return err
	}

	var baseDir string
	if !info.IsDir() {
		baseDir = filepath.Base(sourcePath)
	}

	if err := filepath.Walk(sourcePath, func(pth string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isSymlink, evaledPath, err := checkSymlink(pth)
		if err != nil {
			return err
		}

		originalPath := pth

		if isSymlink {
			pth = evaledPath
		}

		info, err = os.Lstat(pth)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			if isSymlink {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(originalPath, sourcePath))
			} else {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(pth, sourcePath))

			}

		} else {
			baseDir = "/"
		}

		if info.IsDir() {
			header.Name, err = filepath.Rel(sourcePath, pth)
			if err != nil {
				return err
			}

			header.Name += info.Name() + "/"
		}

		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fContent, err := ioutil.ReadFile(pth)
		if err != nil {
			return err
		}

		_, err = writer.Write(fContent)
		return err

	}); err != nil {
		return err
	}

	return nil
}

// checkSymlink If the file is a symbolic link it will evaulate the path name.

// If the target is a symbolic link it will return true, the evaulated path, and the original path.
// If the target is  not a symbolic link it will return false.
// If there is an error it will return the error.
func checkSymlink(path string) (isSymlink bool, evaledPath string, error error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, "", err
	}

	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		evaledPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return true, "", err
		}

		return true, evaledPath, nil
	}
	return false, "", nil
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

// It will warn the user if the zip has already exist at the destination.
// Just log a warning, still continues.
func checkAlreadyExist(destination string) {
	targetPathWithExtension := destination

	if !strings.HasPrefix(targetPathWithExtension, ".zip") {
		targetPathWithExtension += ".zip"
	}
	splittedTargetPathWithExtension := strings.Split(targetPathWithExtension, "/")

	targetName := targetPathWithExtension
	if len(splittedTargetPathWithExtension) > 0 {
		targetName = splittedTargetPathWithExtension[len(splittedTargetPathWithExtension)-1]
	}

	if err := input.ValidateIfPathExists(targetPathWithExtension); err == nil {
		targetAlreadyExist := fmt.Sprintf("The %s already exists at location: %s ", targetName, targetPathWithExtension)

		log.Warnf(targetAlreadyExist)
		fmt.Println()
	}
}
