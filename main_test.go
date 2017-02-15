package main

import (
	"bytes"
	"io/ioutil"
	//"os"
	"strings"
	"testing"
)

var inputFile = []byte(`
struct A {
	int        a;	// field 1    
	int     b;   // field 2   
	char*      c;   	  // field 3
};
    
struct X {
    struct D     x;  // field 1
    const char*    y;   // field 2
    char*    z;            // field 3
`)

var expectedFile = []byte(`
struct A {
	int    a;  // field 1
	int    b;  // field 2
	char*  c;  // field 3
};

struct X {
    struct D     x;  // field 1
    const char*  y;  // field 2
    char*        z;  // field 3
`)

func listSpaces(s string) string {
	return strings.Replace(strings.Replace(s, "\t", "▸▸▸▸", -1), " ", "·", -1)
}

func BenchmarkFixTabstopsWritten(b *testing.B) {
	for i := 0; i < b.N; i++ {
		target, err := os.Create(os.DevNull)
		if err != nil {
			panic(err)
		}
		FixTabstops(bytes.NewReader(inputFile), target)
		target.Close()
	}
}

func BenchmarkFixTabstopsUnchanged(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FixTabstops(bytes.NewReader(expectedFile), ioutil.Discard)
	}
}

func TestFixTabstops(t *testing.T) {
	var out bytes.Buffer

	if err := FixTabstops(bytes.NewReader(inputFile), &out); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(out.Bytes(), expectedFile) {
		t.Errorf(
			"\n--- input ---\n%s\n--- expected ---\n%s\n--- actual ---\n %s",
			listSpaces(string(inputFile)),
			listSpaces(string(expectedFile)),
			listSpaces(string(out.Bytes())),
		)
	}
}

func TestFixTabstopsUnchanged(t *testing.T) {
	err := FixTabstops(bytes.NewReader(expectedFile), ioutil.Discard)
	if err != Unchanged {
		t.Errorf("Expected Unchanged, got %v", err)
	}
}
