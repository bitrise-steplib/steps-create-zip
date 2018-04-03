package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	if !strings.HasSuffix(destination, ".zip") {
		destination += ".zip"
	}

	if err := checkAlreadyExist(destination); err != nil {
		failf("Issue with compress: %s", err)
	}

	if err := ensureZIP(cfg.SourcePath, destination); err != nil {
		failf("Issue with compress: %s", err)
	}

}

func ensureZIP(sourcePath string, destination string) error {
	dirOftargetPath := filepath.Dir(destination)

	if err := os.MkdirAll(dirOftargetPath, 0755); err != nil {
		return fmt.Errorf("Failed to create directory, error: %s", err)
	}

	info, err := os.Lstat(sourcePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		if err := ziputil.ZipDir(sourcePath, destination, false); err != nil {
			return err
		}

	} else {
		if err := ziputil.ZipFile(sourcePath, destination); err != nil {
			return err
		}
	}

	return nil
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

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}
