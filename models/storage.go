package models

import (
	"encoding/json"
	"github.com/joyde68/blog/pkg"
	"io/ioutil"
	"os"
	"path"
)

var (
	//appVersion int
	// global data storage instance
	Storage *JsonStorage
	// global tmp data storage instance. Temp data are generated for special usages, will not backup.
	TmpStorage *JsonStorage
)

type JsonStorage struct {
	Dir string
}

func (jss *JsonStorage) Init(dir string) {
	jss.Dir = dir
}

func (jss *JsonStorage) Has(key string) bool {
	/*
	file := path.Join(jss.Dir, key+".json")
	_, e := os.Stat(file)
	*/
	return pkg.IsFile(jss.Dir + key + ".json")
}

func (jss *JsonStorage) Get(key string, v interface{}) {
	file := path.Join(jss.Dir, key+".json")
	bytes, e := ioutil.ReadFile(file)
	if e != nil {
		println("read storage '" + key + "' error")
		return
	}
	e = json.Unmarshal(bytes, v)
	if e != nil {
		println("json decode '" + key + "' error")
	}
}

func (jss *JsonStorage) Set(key string, v interface{}) {
	locker.Lock()
	defer locker.Unlock()

	bytes, e := json.Marshal(v)
	if e != nil {
		println("json encode '" + key + "' error")
		return
	}
	file := path.Join(jss.Dir, key+".json")
	e = ioutil.WriteFile(file, bytes, 0777)
	if e != nil {
		println("write storage '" + key + "' error")
	}
}

func (jss *JsonStorage) GetDir(name string) {
	os.MkdirAll(path.Join(jss.Dir, name), os.ModePerm)
}

func loadAllData() {
	//loadVersion()
	LoadSettings()
	LoadNavigators()
	LoadUsers()
	LoadTokens()
	LoadContents()
	LoadMessages()
	LoadReaders()
	LoadComments()
	LoadFiles()
}

// TimeInc returns time step value devided by d int with time unix stamp.
func (jss *JsonStorage) TimeInc(d int) int {
	return int(pkg.Now())%d + 1
}

// Init does model initialization.
// If first run, write default data.
// v means app.Version number. It's needed for version data.
func Init() {
	//appVersion = v
	Storage = new(JsonStorage)
	Storage.Init("data")
	TmpStorage = new(JsonStorage)
	TmpStorage.Init("tmp/data")

	if !CheckInstall() {
		DoInstall()
	}
}

// All loads all data from storage to memory.
// Start timers for content, comment and message.
func All() {
	loadAllData()
	// generate indexes
	SyncIndexes()
	// start model timer, do all timer stuffs
	StartModelTimer()
}

func SyncIndexes() {
	// generate indexes
	generatePublishArticleIndex()
	generateContentTmpIndexes()
}

// SyncAll writes all current memory data to storage files.
func SyncAll() {
	SyncContents()
	SyncMessages()
	SyncFiles()
	SyncReaders()
	SyncSettings()
	SyncNavigators()
	SyncTokens()
	SyncUsers()
	//SyncVersion()
}
