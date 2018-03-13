package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/input"
)

const (
	targetExtension string = ".zip"
)

// ConfigsModel ...
type ConfigsModel struct {
	SourcePath string
	TargetPath string
}

func main() {
	configs := createConfigsModelFromEnvs()

	fmt.Println()
	configs.print()

	if err := configs.validate(); err != nil {
		fail("Issue with input: %s", err)
	}

	_, err := compress(configs.SourcePath, configs.TargetPath)

	if err != nil {
		fail("Issue with compress: %s", err)
		return
	}
}

func fixTargetExtension(targetPath string) string {
	if strings.HasSuffix(targetPath, targetExtension) {
		return targetPath
	}
	return fmt.Sprintf("%s%s", targetPath, targetExtension)
}

func compress(sourcePath string, targetPath string) (string, error) {
	targetPath = fixTargetExtension(targetPath)

	dirOftargetPath := filepath.Dir(targetPath)
	fmt.Print(dirOftargetPath)
	if err := input.ValidateIfPathExists(dirOftargetPath); err != nil {
		log.Printf("targetRootPath: %s does not exist", dirOftargetPath)
		os.MkdirAll(dirOftargetPath, os.ModePerm)
	}
	zipfile, err := os.Create(targetPath)

	if err != nil {
		return "", err
	}

	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Lstat(sourcePath)
	if err != nil {
		return "", err
	}

	var baseDir string
	if !info.IsDir() {
		baseDir = filepath.Base(sourcePath)
	}

	filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isSymlink, evaledPath, originalPath, err := checkSymlink(path)
		if err != nil {
			return err
		}

		if isSymlink {
			path = evaledPath
		}

		info, err = os.Lstat(path)
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
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, sourcePath))
			}

		} else {
			baseDir = "/"
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return fmt.Sprint(zipfile), nil
}

func checkSymlink(path string) (isSymlink bool, evaledPath string, originalPath string, error error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, path, path, err
	}

	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		evaledPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return true, path, path, err
		}

		return true, evaledPath, path, nil
	}
	return false, path, path, nil
}

func createConfigsModelFromEnvs() ConfigsModel {

	return ConfigsModel{
		SourcePath: os.Getenv("source_path"),
		TargetPath: os.Getenv("target_path"),
	}
}

func (configs ConfigsModel) print() {
	log.Infof("Create ZIP configs:")
	log.Printf("- SourcePath: %s", configs.SourcePath)
	log.Printf("- TargetPath: %s", configs.TargetPath)
}

func (configs ConfigsModel) validate() error {
	if err := input.ValidateIfNotEmpty(configs.SourcePath); err != nil {
		return errors.New("issue with input SourcePath: " + err.Error())
	}

	if err := input.ValidateIfNotEmpty(configs.TargetPath); err != nil {
		return errors.New("issue with input TargetPath: " + err.Error())
	}

	if err := input.ValidateIfDirExists(configs.SourcePath); err != nil {
		return errors.New("issue with input SourcePath: " + err.Error())
	}

	checkAlreadyExist(configs)
	return nil
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func checkAlreadyExist(configs ConfigsModel) {
	targetPathWithExtension := fixTargetExtension(configs.TargetPath)
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
