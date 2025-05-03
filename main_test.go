package main

import (
	"testing"
)

func bytesEqual(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range len(a) {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestBytesEqual(t *testing.T) {
	var a []byte = []byte{0xff, 0xfa}
	var b []byte = []byte{0xff, 0xfa}

	if !bytesEqual(a, b) {
		t.Errorf("%x should equal %x", a, b)
	}
}

func TestWriteQNames(t *testing.T) {
	expected := []byte{
		0x07, 0x66, 0x72, 0x6f, 0x6f, 0x70, 0x6c, 0x65,
		0x03, 0x6e, 0x65, 0x74, 0x00,
	}

	var cases []string = []string{
		"froople.net",
	}

	buf := make([]byte, 0)
	buf = writeQNames(buf, cases)

	if !bytesEqual(buf, expected) {
		t.Errorf("expected '%x', got '%x'", expected, buf)
	}
}

func TestSplitUint16(t *testing.T) {
	expected := [][]byte{
		[]byte{0xff, 0xfa},
		[]byte{0xaa, 0xa1},
	}

	cases := []uint16{
		0xfffa, 
		0xaaa1,
	}

	for i, c := range cases {
		a, b := splitUint16(c)
		e1, e2 := expected[i][0], expected[i][1]
		if a != e1 || b != e2 {
			t.Errorf("expected '%02x', '%02x' for case '%04x' but got '%02x', '%02x'", e1, e2, c, a, b)
		}
	}

}

func TestWriteQuestion(t *testing.T) {
	expected := []byte {
    	0x07, 0x66, 0x72, 0x6f, 0x6f, 0x70, 0x6c, 
		0x65, 0x03, 0x6e, 0x65, 0x74, 0x00, 0x00,
		0x01, 0x00, 0x01,
	}

	buf := make([]byte, 0)
	buf = writeQuestion(buf, []string{"froople.net"}, QTypeA, QClassIn)

	if !bytesEqual(buf, expected) {
		t.Errorf("expected '%x', got '%x'", expected, buf)	
	}
}
