package models

import (
	"gopkg.in/macaron.v1"
	"time"
)

/*
type logItem struct {
	Name       string
	CreateTime int64
	Text       string
}

// LogErrors logs error bytes to tmp/log directory.
func AddLog(bytes []byte) {
	dir := "data/log"
	file := path.Join(dir, pkg.DateInt64(pkg.Now(), "MMDDHHmmss.log"))
	ioutil.WriteFile(file, bytes, 0755)
}

func Logs() []*logItem {
	logs := make([]*logItem, 0)
	dir := filepath.Join("tmp","log")
	filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		if err == nil {
			if info.IsDir() {
				return nil
			}
			ext := filepath.Ext(info.Name())
			if ext != ".log" {
				return nil
			}
			bytes, e := ioutil.ReadFile(filepath.Join(dir, info.Name()))
			if e != nil {
				return nil
			}
			l := new(logItem)
			l.Name = info.Name()
			l.CreateTime = info.ModTime().Unix()
			l.Text = string(bytes)
			logs = append([]*logItem{l}, logs...)
		}
		return nil
	})
	return logs
}

func RemoveLog(file string) {
	f := filepath.Join(filepath.Join("data", "log"), file)
	os.Remove(f)
}

func RemoveAllLog() {
	f := filepath.Join("data", "log")
	os.Remove(f)
	os.MkdirAll(f, 0755)
}
*/

var (
	loginErrorCount map[string]int
	loginErrorLog []loginErrLog
)

type loginErrLog struct {
	User string
	Password string
	Message string
	Date int64
	Ip string
	UserAgent string
}

func LogInit() {
	loginErrorCount = make(map[string]int)
	loginErrorLog = make([]loginErrLog, 0)
}

func AddLoginErrLog(msg string, context *macaron.Context) {
	ip := context.RemoteAddr()
	loginErrorCount[ip]++
	loginErrorLog = append(loginErrorLog, loginErrLog{
		User: context.Query("user"),
		Password: context.Query("password"),
		Message: msg,
		Date: time.Now().Unix(),
		Ip: ip,
		UserAgent: context.Req.Header.Get("User-Agent"),
	})
}

func GetLoginErrCount(ip string) int {
	return loginErrorCount[ip]
}
