package data

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
)

const (
	AppName = "ribs"
)

func MakeDirIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0775)
		if err != nil {
			panic("Could not create directory: " + err.Error())
		} else {
			log.Print(fmt.Sprintf("Created cache directory: %s", path))
		}
	}
}

func CacheRoot() string {
	if runtime.GOOS == "windows" {
		return path.Join(os.Getenv("USERPROFILE"), "AppData", AppName)
	} else {
		return path.Join(os.Getenv("HOME"), fmt.Sprintf(".%s", AppName))
	}
}

func DataPath() string {
	return path.Join(CacheRoot(), "data")
}
