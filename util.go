package main

import (
    "fmt"
    "strings"
)

func splitUint16(v uint16) (uint8, uint8) {
	return uint8((v >> 8) & 0xff), uint8(v & 0xff)
}

func writeStringToAscii(b []byte, s string) []byte {
	for _, c := range s {
		if uint8(c) > 127 {
			return b
		}
		b = append(b, uint8(c))
	}
	return b
}

func writeLabel(b []byte, domain string) []byte {
	for _, part := range strings.Split(domain, ".") {
		if len(part) > 255 {
			return b
		}
		b = append(b, uint8(len(part)))
		b = writeStringToAscii(b, part)
	}
	return b
}

func readLabelsPtr(b []byte, offset int) (int, []string) {
    labels := make([]string, 0)
    ptr := offset

    l := b[ptr]
    ptr++

    l = (l << 2) >> 2 // Shave off the 2 MSBs
    o := bytesToUint16(l, b[ptr])
    ptr++

    _, labels = readLabels(b, int(o))

    return ptr - offset, labels
}

func readLabel(b []byte, offset int) (int, string) {
    label := ""
    ptr := offset

    l := int(b[ptr])
    ptr++

    label = string(b[ptr:ptr+l])
    ptr += l

    return ptr - offset, label
}

func readLabels(b []byte, offset int) (int, []string) {
    labels := make([]string, 0)
    ptr := offset

    for {
        l := b[ptr] 
        if l == 0x00 {
            ptr++
            break
        }

        if l >= 0b11000000 {
            shifted, ptrLabels := readLabelsPtr(b, ptr)
            for _, label := range ptrLabels {
                labels = append(labels, string(label))
            }
            ptr+=shifted
            break

        } else {
            shifted, label := readLabel(b, ptr)
            labels = append(labels, string(label))
            ptr+=shifted
        }
    }

    return ptr - offset, labels
}

func printByteHex(b []byte) {
	for i := range len(b) {
		fmt.Printf("%#02x ", b[i])
	}
	fmt.Printf("\n")
}

func writeQNames(b []byte, qnames []string) []byte {
	for _, d := range qnames {
		b = writeLabel(b, d)
	}
	return append(b, 0x00)
}

func writeSplitUint16(b []uint8, v uint16) []byte {
	u, l := splitUint16(v)
	b = append(b, u)
	b = append(b, l)
	return b
}

func bytesToUint16(a, b byte) uint16 {
	return uint16(a)<<8 | uint16(b)
}

func bytesToUint32(a, b, c, d byte) uint32 {
	return uint32(a)<<24 | uint32(b)<<16 | uint32(c)<<8 | uint32(d)
}

func bytesToIPString(b [4]byte) string {
    return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}
