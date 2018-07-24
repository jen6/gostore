package gostore

import "testing"

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
		err := GetApk("/tmp", test.Data)
		if test.Result != err {
			t.Error("Expact : ", test.Result, " Result : ", err)
		}
	}
}
