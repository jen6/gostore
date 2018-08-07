package gostore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const dbFileName = "apks.db"
const translateLayout = "2006. 1. 2."

func TransformToEucKrStr(in string) string {
	out, _, _ := transform.String(korean.EUCKR.NewEncoder(), in)
	return out
}

func ConvertLastUpdate(timestr string) time.Time {
	t, _ := time.Parse(translateLayout, timestr)
	return t
}

type ApkInfo struct {
	Id         int       `gorm:"AUTO_INCREMENT; not null; primary_key"`
	PkgName    string    `gorm:"not null; unique"`
	LastUpdate time.Time `gorm:"not null"`
}

func SearchApkInfo(proxy, pkgName string, token Token) (ApkInfo, error) {
	eucPkgName := TransformToEucKrStr(pkgName)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		1*time.Minute,
	)
	defer cancel()

	cmd := exec.CommandContext(ctx, "gplaycli",
		"-s", eucPkgName,
		"-ts", token.TokenStr,
		"-g", token.GsfStr,
		"-n", "1",
	)

	env := os.Environ()

	proxy_env := fmt.Sprintf("https_proxy=%s", proxy)
	env = append(env, proxy_env)
	cmd.Env = env

	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(string(out))
		return ApkInfo{}, err
	}

	apk, ok := parseApkInfo(out, eucPkgName)
	if !ok {
		return ApkInfo{}, errors.New("pkg not found!")
	}

	return apk, nil
}

func parseApkInfo(out []byte, pkgName string) (ApkInfo, bool) {
	rows := bytes.Split(out, []byte{'\n'})
	if len(rows) != 3 {
		fmt.Println("row parse error")
		return ApkInfo{}, false
	}

	cols := bytes.Split(rows[1], []byte("  "))
	if len(cols) < 8 {
		fmt.Println("col parse error")
		return ApkInfo{}, false
	}

	var apk ApkInfo

	apk.LastUpdate = ConvertLastUpdate(string(cols[4]))
	apk.PkgName = string(cols[5])

	if apk.PkgName != pkgName {
		fmt.Println("pkg name not match : ", apk.PkgName, " : ", pkgName)
		return apk, false
	}

	return apk, true
}

type ApkManager struct {
	ApkChecker map[string]ApkInfo
	db         *gorm.DB
}

func (am *ApkManager) load() error {
	var apks []ApkInfo
	if err := am.db.Find(&apks).Error; err != nil {
		return err
	}
	for _, apk := range apks {
		am.ApkChecker[apk.PkgName] = apk
	}
	return nil
}

func (am *ApkManager) Close() {
	am.db.Close()
}

func (am *ApkManager) NewApk(apk ApkInfo) error {
	if newr := am.db.NewRecord(apk); !newr {
		return errors.New("already exist record")
	}
	err := am.db.Create(&apk).Error
	if err != nil {
		return err
	}
	am.ApkChecker[apk.PkgName] = apk
	return nil
}

func (am *ApkManager) UpdateApk(apk ApkInfo) error {
	if newr := am.db.NewRecord(apk); newr {
		return errors.New("not exist record")
	}

	err := am.db.Save(&apk).Error
	if err != nil {
		return err
	}
	am.ApkChecker[apk.PkgName] = apk
	return nil
}

func NewApkManager(apkPath string) (*ApkManager, error) {
	apkPath = filepath.Join(apkPath, dbFileName)
	db, err := gorm.Open("sqlite3", apkPath)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&ApkInfo{})

	manager := ApkManager{db: db, ApkChecker: map[string]ApkInfo{}}
	if err = manager.load(); err != nil {
		return &manager, err
	}

	return &manager, nil
}
