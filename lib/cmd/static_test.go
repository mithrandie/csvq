package cmd

import (
	"os"
	"reflect"
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
	t1 := Now()
	time.Sleep(time.Millisecond)
	t2 := Now()

	if !t1.Equal(t2) {
		t.Errorf("function Now() returns different time")
	}
}

func TestGetReader(t *testing.T) {
	fp := os.Stdout

	r := GetReader(fp, UTF8)
	if reflect.TypeOf(r).String() != "*bufio.Reader" {
		t.Errorf("reader = %q, want %q", reflect.TypeOf(r).String(), "*bufio.Reader")
	}

	r = GetReader(fp, SJIS)
	if reflect.TypeOf(r).String() != "*transform.Reader" {
		t.Errorf("reader = %q, want %q", reflect.TypeOf(r).String(), "*transform.Reader")
	}
}
