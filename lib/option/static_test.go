package option

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
	p1, _ := GetLocation("Local")
	p2, _ := GetLocation("Local")

	if p1 != p2 {
		t.Error("function GetLocation() returns different pointer")
	}

	_, err := GetLocation("Err")
	if err == nil {
		t.Error("function GetLocation() must return an error")
	} else if err.Error() != "timezone \"Err\" does not exist" {
		t.Errorf("function GetLocation() returned an error %q, want %q", err.Error(), "timezone \"Err\" does not exist")
	}

}

func TestNow(t *testing.T) {
	loc, _ := GetLocation("UTC")
	TestTime, _ = time.ParseInLocation("2006-01-02 15:04:05.999999999", "2012-02-01 12:03:23", loc)

	expect := time.Date(2012, 2, 1, 12, 3, 23, 0, loc)

	if !Now(loc).Equal(expect) {
		t.Errorf("function Now() returns current time")
		t.Log(Now(loc))
	}

	TestTime = time.Time{}
	if Now(loc).Equal(expect) {
		t.Errorf("function Now() returns fixed time")
		t.Log(Now(loc))
	}

}
