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
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// config ...
type config struct {
	SourcePath string `env:"source_path,required"`
	TargetDir  string `env:"target_dir"`
}

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(cfg)

	fmt.Println()
	cfg.print()

	if err := cfg.validate(); err != nil {
		fail("Issue with input: %s", err)
	}

	_, err := ensureZIPExtension(cfg.SourcePath, cfg.TargetDir)

	if err != nil {
		fail("Issue with compress: %s", err)
		return
	}
}

// fixTargetExtension Adds the file extension to the end of the targetpath if it's needed.
func fixTargetExtension(targetDir string) string {
	targetExtension := ".zip"

	if strings.HasSuffix(targetDir, targetExtension) {
		return targetDir
	}
	return fmt.Sprintf("%s%s", targetDir, targetExtension)
}

func ensureZIPExtension(sourcePath string, targetDir string) (string, error) {
	targetDir = fixTargetExtension(targetDir)

	ensureDir(targetDir)

	zipfile, err := os.Create(targetDir)
	if err != nil {
		return "", err
	}

	defer func() {
		if err = zipfile.Close(); err != nil {
			log.Errorf("%s", err)
		}
	}()

	archive := zip.NewWriter(zipfile)

	defer func() {
		if err = archive.Close(); err != nil {
			log.Errorf("%s", err)
		}
	}()

	info, err := os.Lstat(sourcePath)
	if err != nil {
		return "", err
	}

	var baseDir string
	if !info.IsDir() {
		baseDir = filepath.Base(sourcePath)
	}

	err = filepath.Walk(sourcePath, func(pth string, fileInfo os.FileInfo, err error) error {
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

		file, err := os.Open(pth)
		if err != nil {
			return err
		}

		defer func() {
			if err = file.Close(); err != nil {
				log.Errorf("%s", err)
			}
		}()

		_, err = io.Copy(writer, file)
		return err
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprint(zipfile), nil
}

// If the file is a symbolic link it will evaulate the path name.

// If the target is a symbolic link it will return true, the evaulated path, and the original path.
// If the target is  not a symbolic link it will return false
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

// ensureDir will check the target path's existence.
// If the targetDir does not exist it will create that.
func ensureDir(targetDir string) {
	dirOftargetPath := filepath.Dir(targetDir)
	fmt.Print(dirOftargetPath)
	if err := input.ValidateIfPathExists(dirOftargetPath); err != nil {
		log.Printf("targetRootPath: %s does not exist", dirOftargetPath)
		err = os.MkdirAll(dirOftargetPath, os.ModePerm)
		if err != nil {
			log.Errorf("Failed to create directory, error: %s", err)
		}
	}
}

func (configs config) print() {
	log.Infof("Create ZIP configs:")
	log.Printf("- SourcePath: %s", configs.SourcePath)
	log.Printf("- TargetDir: %s", configs.TargetDir)
}

func (configs config) validate() error {
	if err := input.ValidateIfNotEmpty(configs.SourcePath); err != nil {
		return errors.New("issue with input SourcePath: " + err.Error())
	}

	if err := input.ValidateIfNotEmpty(configs.TargetDir); err != nil {
		return errors.New("issue with input TargetDir: " + err.Error())
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

// It will warn the user if the zip has already exist at the targetDir.
// Just log a warning, still continues.
func checkAlreadyExist(configs config) {
	targetPathWithExtension := fixTargetExtension(configs.TargetDir)
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
