package hashscaner

import (
	"crypto/sha256"
	"fmt"
	"io"
	"light-defender-client/pkg/config"
	folderfile "light-defender-client/pkg/folder_file"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

type HashScanerI interface {
	RunScheduler()
	RunManual()
	Scan()
	ScanFolder(string) ([]folderfile.FolderFile, error)
	ScanFile(string) folderfile.FolderFile
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
	resultMap := []folderfile.FolderMap{}
	for _, folder := range hs.hsConfig.WatchingFolders {
		info, err := os.Stat(folder)
		curFolderMap := folderfile.FolderMap{
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
			curFolderMap.Files = []folderfile.FolderFile{file}
			resultMap = append(resultMap, curFolderMap)
			continue
		}
	}
	responseMap := map[string]interface{}{}

}

func (hs *hashScaner) ScanFolder(folder string) ([]folderfile.FolderFile, error) {
	files := []folderfile.FolderFile{}
	folderFiles, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	for _, file := range folderFiles {
		if slices.Contains(hs.hsConfig.Exceptions, filepath.Join(folder, file.Name())) {
			continue
		}
		if file.Type().IsDir() {
			folderResult, err := hs.ScanFolder(file.Name())
			if err != nil {
				files = append(files, folderfile.FolderFile{
					File:   file.Name(),
					Error:  err.Error(),
					Result: false,
				})
				continue
			}
			files = append(files, folderResult...)
		}
		if file.Type().IsRegular() {
			files = append(files, hs.ScanFile(filepath.Join(folder, file.Name())))
		}
	}
	return files, nil
}

func (hs *hashScaner) ScanFile(file string) folderfile.FolderFile {
	data, err := os.Open(file)
	if err != nil {
		return folderfile.FolderFile{
			File:   file,
			Error:  err.Error(),
			Result: false,
		}
	}
	defer data.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, data); err != nil {
		return folderfile.FolderFile{
			File:   file,
			Error:  err.Error(),
			Result: false,
		}
	}

	hashBytes := hasher.Sum(nil)

	return folderfile.FolderFile{
		File:   file,
		Hash:   fmt.Sprintf("%x", hashBytes),
		Result: true,
		Error:  "",
	}
}
