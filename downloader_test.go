package gostore

import (
	"os"
	"path/filepath"
	"testing"
)

const tmpPath = "/tmp"

type testGetApk struct {
	Data   AppInfo
	Result error
}

func TestGetApk(t *testing.T) {
	tests := []testGetApk{
		{
			AppInfo{AppName: "Tik Tok Live Photo", PkgName: "com.ss.android.ugc.tiktok.livewallpaper"},
			nil,
		},
	}
	for _, test := range tests {
		err := GetApk(tmpPath, test.Data)
		if test.Result != err {
			t.Error("Expact : ", test.Result, " Result : ", err)
		} else {
			os.Remove(filepath.Join(tmpPath, test.Data.PkgName+".apk"))
		}
	}
}
