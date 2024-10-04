package main

import (
	"strings"
	"testing"
)

var testStrings []string = []string{
	"test", "hello world", "how much wood could a woodchuck chuck if a woodchuck could chuck wood",
}
var testChars []byte = []byte{'e', 'w', 'd'}

func Test_findChar(t *testing.T) {
	type args struct {
		s string
		c byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"test1", args{testStrings[0], testChars[0]}, strings.IndexByte(testStrings[0], testChars[0])},
		{"hw", args{testStrings[1], testChars[1]}, strings.IndexByte(testStrings[1], testChars[1])},
		{"woodchuck", args{testStrings[2], testChars[2]}, strings.IndexByte(testStrings[2], testChars[2])},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findChar(tt.args.s, tt.args.c); got != tt.want {
				t.Errorf("findChar() = %v, want %v", got, tt.want)
			}
		})
	}
}
