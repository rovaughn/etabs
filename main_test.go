package main

import (
	"bytes"
	"testing"
)

var inputFile = `
struct A {
	int        a;	// field 1    
	int     b;   // field 2   
	char*      c;   	  // field 3
};
    
struct X {
    struct D     x;  // field 1
    const char*    y;   // field 2
    char*    z;            // field 3
`

var expectedFile = `
struct A {
	int    a;  // field 1
	int    b;  // field 2
	char*  c;  // field 3
};

struct X {
    struct D     x;  // field 1
    const char*  y;  // field 2
    char*        z;  // field 3
`

func TestFixTabstops(t *testing.T) {
	var out bytes.Buffer

	if err := FixTabstops(bytes.NewReader([]byte(inputFile)), &out); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(out.Bytes(), []byte(expectedFile)) {
		t.Errorf("expected\n%s\ngot\n %s", expectedFile, string(out.Bytes()))
	}
}
