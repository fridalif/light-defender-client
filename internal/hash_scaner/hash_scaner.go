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
	ScanFolder(string) (map[string]interface{}, error)
	ScanFile(string) (string, error)
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
	resultMap := make(map[string]interface{})
	for _, folder := range hs.hsConfig.WatchingFolders {
		info, err := os.Stat(folder)
		if err != nil {
			resultMap[folder] = map[string]interface{}{
				"error":  err.Error(),
				"hashes": map[string]interface{}{},
				"result": false,
			}
			continue
		}
		if info.Mode().IsDir() {
			folderHashes, err := hs.ScanFolder(folder)
			if err != nil {
				resultMap[folder] = map[string]interface{}{
					"error":  err.Error(),
					"hashes": map[string]interface{}{},
					"result": false,
				}
				continue
			}
			resultMap[folder] = map[string]interface{}{
				"error":  "",
				"hashes": folderHashes,
				"result": true,
			}
			continue
		}
		if info.Mode().IsRegular() {
			hash, err := hs.ScanFile(folder)
			if err != nil {
				resultMap[folder] = map[string]interface{}{
					"error":  err.Error(),
					"hashes": map[string]interface{}{},
					"result": false,
				}
				continue
			}
			resultMap[folder] = map[string]interface{}{
				"error": "",
				"hashes": map[string]interface{}{
					folder:   hash,
					"error":  "",
					"result": true,
				},
				"result": true,
			}
		}
	}
}

func (hs *hashScaner) ScanFolder(folder string) (map[string]interface{}, error) {
	hashesMap := make(map[string]interface{})
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	for object := range files {

	}
}

func (hs *hashScaner) ScanFile(file string) (string, error) {

}
