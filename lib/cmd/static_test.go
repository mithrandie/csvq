package cmd

import (
	"testing"
	"time"
)

func TestGetRand(t *testing.T) {
	p1 := GetRand()
	p2 := GetRand()

	if p1 != p2 {
		t.Errorf("function GetRand() returns different pointer")
	}
}

func TestGetLocation(t *testing.T) {
	p1 := GetLocation()
	p2 := GetLocation()

	if p1 != p2 {
		t.Errorf("function GetLocation() returns different pointer")
	}
}

func TestNow(t *testing.T) {
	TestTime, _ = time.ParseInLocation("2006-01-02 15:04:05.999999999", "2012-02-01 12:03:23", GetLocation())

	expect := time.Date(2012, 2, 1, 12, 3, 23, 0, GetLocation())

	if !Now().Equal(expect) {
		t.Errorf("function Now() returns current time")
		t.Log(Now())
	}

	TestTime = time.Time{}
	if Now().Equal(expect) {
		t.Errorf("function Now() returns fixed time")
		t.Log(Now())
	}

}
