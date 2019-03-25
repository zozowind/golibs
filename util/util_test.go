package util

import (
	"fmt"
	"testing"
	"time"
)

func TestRound(t *testing.T) {
	over := 20
	total := 22
	p := Round(float64((over+1)*100)/float64(total), 3)
	fmt.Printf("%#v", p)
}

func TestCardBinCheck(t *testing.T) {
	rsp, err := CardBinCheck("6212261001070083099")
	if nil != err {
		t.Errorf(err.Error())
		return
	}
	fmt.Printf("%#v", rsp)
}

func TestTimeDurationToAboutStr(t *testing.T) {
	fmt.Println(TimeDurationToAboutStr(35 * 24 * time.Hour))
	fmt.Println(TimeDurationToAboutStr(10 * 24 * time.Hour))
	fmt.Println(TimeDurationToAboutStr(80 * time.Hour))
	fmt.Println(TimeDurationToAboutStr(6 * time.Hour))
	fmt.Println(TimeDurationToAboutStr(54 * time.Minute))
	fmt.Println(TimeDurationToAboutStr(54 * time.Second))
}
