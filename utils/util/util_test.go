package util

import (
	"reflect"
	"strconv"
	"testing"
)

func TestHashAndCompareString(t *testing.T) {
	hash, err := HashString("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashString: %v", err)
	}
	if hash == "correct horse battery staple" {
		t.Fatal("hash must not equal the plain text")
	}
	if err := CompareString(hash, "correct horse battery staple"); err != nil {
		t.Errorf("CompareString with the right password: %v", err)
	}
	if err := CompareString(hash, "wrong password"); err == nil {
		t.Error("CompareString with a wrong password returned nil")
	}
}

func TestRemoveFromSlice(t *testing.T) {
	tests := []struct {
		name  string
		in    []string
		value string
		want  []string
	}{
		{name: "removes first match", in: []string{"red", "green", "blue"}, value: "green", want: []string{"red", "blue"}},
		{name: "missing value", in: []string{"red"}, value: "blue", want: []string{"red"}},
		{name: "empty slice", in: []string{}, value: "red", want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveFromSlice(tt.in, tt.value)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveFromSlice = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	got := Map([]int{1, 2, 3}, strconv.Itoa)
	if want := []string{"1", "2", "3"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Map = %v, want %v", got, want)
	}
}

func TestWrappers(t *testing.T) {
	if got := WrapData(42); got.Data != 42 {
		t.Errorf("WrapData = %+v", got)
	}
	if got := WrapListData([]string{"a"}); !reflect.DeepEqual(got.Data.List, []string{"a"}) {
		t.Errorf("WrapListData = %+v", got)
	}
	if got := WrapMessage("ok"); got.Message != "ok" {
		t.Errorf("WrapMessage = %+v", got)
	}
}
