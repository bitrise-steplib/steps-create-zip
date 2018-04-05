package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/ziputil"
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

	destination := cfg.Destination
	destination, err := fixDestination(destination, cfg.SourcePath)
	if err != nil {
		failf("Issue with compress: %s", err)
	}

	if err := checkAlreadyExist(destination); err != nil {
		failf("Issue with compress: %s", err)
	}

	if err := ensureZIP(cfg.SourcePath, destination); err != nil {
		failf("Issue with compress: %s", err)
	}

}

func ensureZIP(sourcePath string, destination string) error {
	info, err := os.Lstat(sourcePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return ziputil.ZipDir(sourcePath, destination, false)
	}

	return ziputil.ZipFile(sourcePath, destination)
}

// checkAlreadyExist will return an error if the zip has already exist at the destination.
func checkAlreadyExist(destination string) error {
	targetName := filepath.Base(destination)

	exist, err := pathutil.IsPathExists(destination)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("The - %s - already exists at location: %s", targetName, destination)
	}

	return nil
}

func fixDestination(destination string, sourcePath string) (string, error) {
	destination = cleanDestination(destination)

	if err := ensureDestinationPath(destination); err != nil {
		return "", err
	}

	isDir, err := checkDestinationIsDir(destination)
	if err != nil {
		return "", err
	}

	if isDir {
		destination = filepath.Join(destination, filepath.Base(sourcePath))
	}
	destination = fixDestinationExt(destination)

	return destination, nil
}

func cleanDestination(destination string) string {
	return filepath.Clean(destination)
}

func fixDestinationExt(destination string) string {
	if filepath.Ext(destination) == "" {
		destination += ".zip"
	}
	return destination
}

func ensureDestinationPath(destination string) error {
	dirOftargetPath := filepath.Dir(destination)
	return os.MkdirAll(dirOftargetPath, 0755)
}

func checkDestinationIsDir(destination string) (bool, error) {
	exist, err := pathutil.IsPathExists(destination)
	if err != nil {
		return false, err
	}

	if !exist {
		return false, nil
	}

	info, err := os.Lstat(destination)
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
