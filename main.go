package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	QT_A = uint16(iota + 1)
)

const (
	QC_IN = uint16(iota + 1)
	QC_CS
	QC_CH
	QC_HS
)

const (
	R         uint16 = 0b1000000000000000
	OP_IQUERY uint16 = 0b0000100000000000
	OP_STATUS uint16 = 0b0001000000000000
	AA        uint16 = 0b0000010000000000
	TC        uint16 = 0b0000001000000000
	RD        uint16 = 0b0000000100000000
	RA        uint16 = 0b0000000010000000
	RC_FMTERR uint16 = 0b0000000000000001
	RC_SRVERR uint16 = 0b0000000000000010
	RC_NAMERR uint16 = 0b0000000000000011
	RC_NOTIMP uint16 = 0b0000000000000100
	RC_REFUSD uint16 = 0b0000000000000101
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

// BUG: This doesn't work when there are multiple domains
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

func writeQuestion(b []byte, qnames []string, qtype uint16, qclass uint16) []byte {
	b = writeQNames(b, qnames)
	b = writeSplitUint16(b, qtype)
	b = writeSplitUint16(b, qclass)
	return b
}

func writeSplitUint16(b []uint8, v uint16) []byte {
	u, l := splitUint16(v)
	b = append(b, u)
	b = append(b, l)
	return b
}

func writeHeader(b []byte, id, flags, qdcount, ancount, nscount, arcount uint16) []byte {
	b = writeSplitUint16(b, id)
	b = writeSplitUint16(b, flags)
	b = writeSplitUint16(b, qdcount)
	b = writeSplitUint16(b, ancount)
	b = writeSplitUint16(b, nscount)
	b = writeSplitUint16(b, arcount)
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

func main() {
	var packet []byte = make([]byte, 0)
	packet = writeHeader(packet, 0x0001, RD, 0x0001, 0x0000, 0x0000, 0x0000)
	packet = writeQuestion(packet, []string{"froople.net"}, QT_A, QC_IN)

	addr, err := net.ResolveUDPAddr("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatalf("failed to resolve udp address: %v\n", err)
	}

	log.Print("dailing udp server\n")
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("failed to open connection: %v\n", err)
	}
	defer conn.Close()

	n, err := conn.Write(packet)
	if err != nil {
		log.Fatalf("failed to write packet: %v\n", err)
	}

	log.Printf("wrote %d byte(s) to server\n", n)

	buf := make([]byte, 1024)
	for {
		log.Println("waiting for data from server")
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("failed to read from server: %v\n", err)
			break
		}

		log.Printf("read %d byte(s) from server\n", n)

		// Header
		var id uint16 = bytesToUint16(buf[0], buf[1])
		var flags uint16 = bytesToUint16(buf[2], buf[3])
		var qdcount uint16 = bytesToUint16(buf[4], buf[5])
		var ancount uint16 = bytesToUint16(buf[6], buf[7])
		var nscount uint16 = bytesToUint16(buf[8], buf[9])
		var arcount uint16 = bytesToUint16(buf[10], buf[11])

		log.Printf("header id=%d, flags=%x, qdc=%d, anc=%d, nsc=%d, arc=%d\n", id, flags, qdcount, ancount, nscount, arcount)

		var ptr int = 12

		// Read Question Section
		for range qdcount {
			var qtype uint16 = 0x0000
			var qclass uint16 = 0x0000

			// Read QNAMEs
            s, ls := readLabels(buf, ptr) 
            for _, l := range ls {
                log.Printf("label = '%s'\n", string(l))
            }
            ptr+=s

			// Read QTYPE
			qtype = bytesToUint16(buf[ptr], buf[ptr+1])
			ptr += 2

			// Read QCLASS
			qclass = bytesToUint16(buf[ptr], buf[ptr+1])
			ptr += 2

			log.Printf("qtype=%d, qclass=%d\n", qtype, qclass)
		}

		// Answer section
		for range ancount {
			var typ uint16 = 0x0000
			var class uint16 = 0x0000
			var ttl uint32 = 0x00000000
			var rdlen uint16 = 0x0000

			// Read QNAMEs
            s, _ := readLabels(buf, ptr) 
            ptr+=s

			// Type
			typ = bytesToUint16(buf[ptr], buf[ptr+1])
			ptr += 2

			// Class
			class = bytesToUint16(buf[ptr], buf[ptr+1])
			ptr += 2

			// TTL
			ttl = bytesToUint32(buf[ptr], buf[ptr+1], buf[ptr+2], buf[ptr+3])
			ptr += 4

			// RDLength
			rdlen = bytesToUint16(buf[ptr], buf[ptr+1])
			ptr += 2

            // The actual IP address!
            ip := bytesToIPString([4]byte{ buf[ptr], buf[ptr+1], buf[ptr+2], buf[ptr+3]})
            ptr += 4

			log.Printf("typ=%d, class=%d, ttl=%d, rdlen=%d, ip=%s\n", typ, class, ttl, rdlen, ip)
		}
	}
}

func readLabelsPtr(b []byte, offset int) (int, [][]byte) {
    labels := make([][]byte, 0)
    ptr := offset

    l := b[ptr]
    ptr++

    l = (l << 2) >> 2 // Shave off the 2 MSBs
    o := bytesToUint16(l, b[ptr])
    ptr++

    _, labels = readLabels(b, int(o))

    return ptr - offset, labels
}

func readLabel(b []byte, offset int) (int, []byte) {
    label := make([]byte, 0)
    ptr := offset

    l := int(b[ptr])
    ptr++

    label = b[ptr:ptr+l]
    ptr += l

    return ptr - offset, label
}

func readLabels(b []byte, offset int) (int, [][]byte) {
    labels := make([][]byte, 0)
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
                labels = append(labels, label)
            }
            ptr+=shifted
            break

        } else {
            shifted, label := readLabel(b, ptr)
            labels = append(labels, label)
            ptr+=shifted
        }
    }

    return ptr - offset, labels
}
