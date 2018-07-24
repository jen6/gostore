package gostore

import (
	"strings"
	"testing"
)

func TestGetNewApps(t *testing.T) {
	_, err := GetNewAppsReader(10)
	if err != nil {
		t.Fatal(err)
	}
}

type InfoTest struct {
	Data     string
	AppNames []string
	PkgNames []string
}

func TestGetAppList(t *testing.T) {
	tests := []InfoTest{
		{
			Data:     `<div class="details"> <a class="card-click-target" href="/store/apps/details?id=com.papegames.evol.kr.release" aria-hidden="true" tabindex="-1"> </a> <a class="title" href="/store/apps/details?id=com.papegames.evol.kr.release" title="러브앤프로듀서" aria-hidden="true" tabindex="-1"> 1. 러브앤프로듀서 <span class="paragraph-end"></span> </a> <div class="subtitle-container"> <a class="subtitle" href="/store/apps/developer?id=Paper+Games" title="Paper Games">Paper Games</a> <span class="price-container"> <span class="paragraph-end"></span> </span> </div><div class="description"> 너의 마음을 두드릴 초능력! <span class="paragraph-end"></span> <a class="card-click-target" href="/store/apps/details?id=com.papegames.evol.kr.release" aria-hidden="true" tabindex="-1"> </a> </div></div>`,
			AppNames: []string{"러브앤프로듀서"},
			PkgNames: []string{"com.papegames.evol.kr.release"},
		},
	}

	for _, testData := range tests {
		reader := strings.NewReader(testData.Data)
		infos := GetNewAppList(reader)
		for i, info := range infos {
			if testData.AppNames[i] != info.AppName {
				t.Error("Expact : ", testData.AppNames[i], " Value : ", info.AppName)
			}
			if testData.PkgNames[i] != info.PkgName {
				t.Error("Expact : ", testData.PkgNames[i], " Value : ", info.PkgName)
			}
		}
	}
}
