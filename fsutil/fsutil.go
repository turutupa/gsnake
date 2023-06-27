package fsutil

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

const FOLDER_NAME = "gsnake"

var (
	locks = make(map[string]*sync.RWMutex)
	m     sync.RWMutex
)

func NewCfgFile(filename string) (string, bool) {
	lock := getLockForFile(filename)
	lock.Lock()
	defer lock.Unlock()

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
	lock := getLockForFile(filename)
	lock.RLock()
	defer lock.RUnlock()

	if configDir, ok := getDir(); !ok {
		return nil, errors.New("Something went wrong retrieving user config dir")
	} else {
		fileDir := filepath.Join(configDir, filename)
		return os.ReadFile(fileDir)
	}
}

func WriteFile(filename string, data []byte) error {
	lock := getLockForFile(filename)
	lock.Lock()
	defer lock.Unlock()

	if configDir, ok := getDir(); !ok {
		return errors.New("Something went wrong retrieving user config dir")
	} else {
		fileDir := filepath.Join(configDir, filename)
		return os.WriteFile(fileDir, data, 700)
	}
}

func AppendToFile(filename string, data string) error {
	lock := getLockForFile(filename)
	lock.Lock()
	defer lock.Unlock()
	if configDir, ok := getDir(); !ok {
		return errors.New("Something went wrong retrieving user config dir")
	} else {
		fileDir := filepath.Join(configDir, filename)
		file, err := os.OpenFile(fileDir, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(data)
		return err
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

func getLockForFile(filename string) *sync.RWMutex {
	m.Lock()
	defer m.Unlock()

	lock, exist := locks[filename]
	if !exist {
		lock = &sync.RWMutex{}
		locks[filename] = lock
	}
	return lock
}
