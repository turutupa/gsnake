package fsutil

import (
	"errors"
	"os"
	"path/filepath"
)

const FOLDER_NAME = "gsnake"

func NewCfgFile(filename string) (string, bool) {
	configDir, ok := mkDirApp()
	if !ok {
		return "", false
	}
	scoreboardFileDir := filepath.Join(configDir, filename)
	if _, err := os.Stat(scoreboardFileDir); os.IsNotExist(err) {
		_, err := os.Create(scoreboardFileDir)
		if err != nil {
			return "", false
		}
	}
	return scoreboardFileDir, true
}

func ReadFile(filename string) ([]byte, error) {
	if configDir, ok := getDir(); !ok {
		return nil, errors.New("Something went wrong retrieving user config dir")
	} else {
		fileDir := filepath.Join(configDir, filename)
		return os.ReadFile(fileDir)
	}
}

func WriteFile(filename string, data []byte) error {
	if configDir, ok := getDir(); !ok {
		return errors.New("Something went wrong retrieving user config dir")
	} else {
		fileDir := filepath.Join(configDir, filename)
		return os.WriteFile(fileDir, data, 700)
	}
}

func OpenFile(filename string) (*os.File, error) {
	if configDir, ok := getDir(); !ok {
		return nil, errors.New("Something went wrong retrieving user config dir")
	} else {
		fileDir := filepath.Join(configDir, filename)
		return os.OpenFile(fileDir, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}
}

func getDir() (string, bool) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", false
	}
	return filepath.Join(configDir, FOLDER_NAME), true
}

func mkDirApp() (string, bool) {
	configDir, ok := getDir()
	if !ok {
		return "", false
	}
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0700)
		if err != nil {
			return "", false
		}
	}
	return configDir, true
}
