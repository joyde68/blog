package models

import (
	"github.com/Unknwon/cae/zip"
	"github.com/joyde68/blog/pkg"
	"os"
	"path"
	"path/filepath"
	"time"
)

var backupDir = "backup"

func init() {
	// close zip terminal output
	zip.Verbose = false
}

// DoBackup backups whole files to zip archive.
// If withData is false, it compresses static files to zip archive without data files, config files and install lock file.
func DoBackup(withData bool) (string, error) {
	os.Mkdir(backupDir, 0755)
	// create zip file name from time unix
	filename := path.Join(backupDir, pkg.DateTime(time.Now(), "YYYYMMDDHHmmss"))
	if withData {
		filename += ".zip"
	} else {
		filename += "_public.zip"
	}
	z, e := zip.Create(filename)
	if e != nil {
		return "", e
	}
	root, _ := os.Getwd()
	if withData {
		// if with data, add install lock file and config file
		lockFile := path.Join(root, "install.lock")
		if pkg.IsFile(lockFile) {
			z.AddFile("install.lock", lockFile)
		}
		configFile := path.Join(root, "config.json")
		if pkg.IsFile(configFile) {
			z.AddFile("config.json", configFile)
		}
	}
	z.AddDir("public/css", path.Join(root, "public", "css"))
	z.AddDir("public/img", path.Join(root, "public", "img"))
	z.AddDir("public/js", path.Join(root, "public", "js"))
	z.AddDir("public/lib", path.Join(root, "public", "lib"))
	z.AddFile("public/favicon.ico", path.Join(root, "public", "favicon.ico"))
	if withData {
		// if with data, backup data files and uploaded files
		z.AddDir("data", path.Join(root, "data"))
		z.AddDir("public/upload", path.Join(root, "public", "upload"))
	}
	z.AddDir("view/default", path.Join(root, "templates", "default"))
	e = z.Flush()
	if e != nil {
		return "", e
	}
	println("backup success in " + filename)
	return filename, nil
}

// RemoveBackupFile removes backup zip file with filename(not filepath).
func RemoveBackupFile(file string) {
	file = path.Join(backupDir, file)
	os.Remove(file)
}

// GetBackupFileAbsPath returns backup zip absolute filepath by filename.
func GetBackupFileAbsPath(name string) string {
	return path.Join(backupDir, name)
}

// GetBackupFile returns fileinfo slice of all backup files.
func GetBackupFiles() ([]os.FileInfo, error) {
	fi := make([]os.FileInfo, 0)
	e := filepath.Walk(backupDir, func(_ string, info os.FileInfo, _ error) error {
		if info == nil {
			return nil
		}
		if !info.IsDir() {
			fi = append([]os.FileInfo{info}, fi...)
		}
		return nil
	})
	return fi, e
}

// StartBackupTimer starts backup operation timer for auto backup stuff.
func StartBackupTimer(t int) {
	SetTimerFunc("backup-data", 144, func() {
		filename, e := DoBackup(true)
		if e != nil {
			CreateMessage("backup", "[0]"+e.Error())
		} else {
			CreateMessage("backup", "[1]"+filename)
		}
		println("backup files in", t, "hours")
	})
}
