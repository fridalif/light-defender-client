package hashscaner

import (
	"light-defender-client/pkg/config"
	"math/rand"
	"os"
	"sync"
	"time"
)

type HashScanerI interface {
	RunScheduler()
	RunManual()
	Scan()
	ScanFolder(string) []FolderFile
	ScanFile(string) FolderFile
}

type hashScaner struct {
	hsConfig    *config.HashScanerConfig
	hsMutex     *sync.Mutex
	logFunction func(string, ...interface{})
}

func NewHashScaner(hsConfig *config.HashScanerConfig, logFunction func(string, ...interface{}), hsMutex *sync.Mutex) HashScanerI {
	return &hashScaner{hsConfig: hsConfig, logFunction: logFunction, hsMutex: hsMutex}
}

func (hs *hashScaner) RunScheduler() {
	for {
		source := rand.NewSource(time.Now().UnixNano())
		rng := rand.New(source)
		randomNum := rng.Uint64() % hs.hsConfig.Interval

		time.Sleep(time.Duration(randomNum) * time.Second)

		go hs.Scan()
	}
}

func (hs *hashScaner) RunManual() {
	go hs.Scan()
}

func (hs *hashScaner) Scan() {
	hs.hsMutex.Lock()
	defer hs.hsMutex.Unlock()
	resultMap := []FolderMap{}
	for _, folder := range hs.hsConfig.WatchingFolders {
		info, err := os.Stat(folder)
		curFolderMap := FolderMap{
			Path: folder,
		}
		if err != nil {
			curFolderMap.Result = false
			curFolderMap.Error = err.Error()
			resultMap = append(resultMap, curFolderMap)
			continue
		}
		if info.Mode().IsDir() {
			curFolderMap.Files = hs.ScanFolder(folder)
			resultMap = append(resultMap, curFolderMap)
			continue
		}
		if info.Mode().IsRegular() {
			curFolderMap.Files = []FolderFile{hs.ScanFile(folder)}
			resultMap = append(resultMap, curFolderMap)
			continue
		}
	}
}

func (hs *hashScaner) ScanFolder(folder string) []FolderFile {
	hashesMap := make(map[string]interface{})
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	for object := range files {

	}
}

func (hs *hashScaner) ScanFile(file string) FolderFile {

}
