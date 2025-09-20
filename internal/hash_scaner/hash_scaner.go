package hashscaner

import (
	"crypto/sha256"
	"fmt"
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
	ScanFolder(string) ([]FolderFile, error)
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
			files, err := hs.ScanFolder(folder)
			if err != nil {
				curFolderMap.Result = false
				curFolderMap.Error = err.Error()
				continue
			}
			curFolderMap.Files = files
			resultMap = append(resultMap, curFolderMap)
			continue
		}
		if info.Mode().IsRegular() {
			file := hs.ScanFile(folder)
			curFolderMap.Files = []FolderFile{file}
			resultMap = append(resultMap, curFolderMap)
			continue
		}
	}
}

func (hs *hashScaner) ScanFolder(folder string) ([]FolderFile, error) {
	files := []FolderFile{}
	folderFiles, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	for _, file := range folderFiles {
		if file.Type().IsDir() {
			folderResult, err := hs.ScanFolder(file.Name())
			if err != nil {
				files = append(files, FolderFile{
					File:   file.Name(),
					Error:  err.Error(),
					Result: false,
				})
				continue
			}
			files = append(files, folderResult...)
		}
		if file.Type().IsRegular() {
			files = append(files, hs.ScanFile(file.Name()))
		}
	}
	return files, nil
}

func (hs *hashScaner) ScanFile(file string) FolderFile {
	data, err := os.ReadFile(file)
	if err != nil {
		return FolderFile{
			File:   file,
			Error:  err.Error(),
			Result: false,
		}
	}

	hash := sha256.Sum256(data)

	return FolderFile{
		File:   file,
		Hash:   fmt.Sprintf("%x", hash),
		Result: true,
		Error:  "",
	}
}
