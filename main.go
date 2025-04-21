package main

import (
    "log"
)

const True  uint8 = 0x01
const False uint8 = 0x00

func id(p []byte, id uint16) []byte {
    p[0] = uint8((id << 8) & 0xff)
    p[1] = uint8(id & 0xff)
    return p
}

func qr(p []byte, qr uint8) []byte {
    if qr > 0x01 {
        qr = 0x01
    }
    p[2] |= (qr << 7)
    return p
}

func opcode(p []byte, op uint8) []byte {
    if op > 0x0f {
        op = 0x0f
    }
    p[2] |= (op << 3)
    return p
}

func aa(p []byte, aa uint8) []byte {
    if aa > True {
        aa = True
    }
    p[2] |= (aa << 2)
    return p
}

func tc(p []byte, tc uint8) []byte {
    if tc > True {
        tc = True
    }
    p[2] |= (tc << 1)
    return p
}

func rd(p []byte, rd uint8) []byte {
    if rd > True {
        rd = True
    }
    p[2] |= rd
    return p
}

func ra(p []byte, ra uint8) []byte {
    if ra > True {
        ra = True
    }
    p[3] |= (ra << 7)
    return p
}

func z(p []byte, z uint8) []byte {
    if z > 0x07 {
        z = 0x07
    }
    p[3] |= (z << 4)
    return p
}

func rcode(p []byte, rcode uint8) []byte {
    if rcode > 0x0f {
        rcode = 0x0f
    }
    p[3] |= rcode
    return p
}

func qdcount(p []byte, qdcount uint16) []byte {
    p[4] = uint8((qdcount << 8) & 0xff)
    p[5] = uint8(qdcount & 0xff)
    return p
}

func ancount(p []byte, ancount uint16) []byte {
    p[6] = uint8((ancount << 8) & 0xff)
    p[7] = uint8(ancount & 0xff)
    return p
}

func nscount(p []byte, nscount uint16) []byte {
    p[8] = uint8((nscount << 8) & 0xff)
    p[9] = uint8(nscount & 0xff)
    return p
}

func arcount(p []byte, arcount uint16) []byte {
    p[10] = uint8((arcount << 8) & 0xff)
    p[11] = uint8(arcount & 0xff)
    return p
}

func printHeader(header []byte) {
    log.Printf("ID:          %08b, %08b\n", header[0],  header[1])
    log.Printf("QR:          %08b\n",       header[2] & 0b10000000)
    log.Printf("OPCODE:      %08b\n",      (header[2] & 0b01111000) << 1)
    log.Printf("AA:          %08b\n",      (header[2] & 0b00000100) << 5)
    log.Printf("TC:          %08b\n",      (header[2] & 0b00000010) << 6)
    log.Printf("RD:          %08b\n",      (header[2] & 0b00000001) << 7)
    log.Printf("RA:          %08b\n",       header[3] & 0b10000000)
    log.Printf("Z:           %08b\n",      (header[3] & 0b01110000) << 1)
    log.Printf("RCODE:       %08b\n",      (header[3] & 0b00001111) << 4)
    log.Printf("QDCOUNT:     %08b, %08b\n", header[4],  header[5])
    log.Printf("ANCOUNT:     %08b, %08b\n", header[6],  header[7])
    log.Printf("NSCOUNT:     %08b, %08b\n", header[8],  header[9])
    log.Printf("ARCOUNT:     %08b, %08b\n", header[10], header[11])
}

func main() {
    var header []byte = make([]byte, 96)
    header = id(header, 0x0001)
    header = qr(header, True)
    header = opcode(header, 0x000a)
    header = aa(header, True)
    header = tc(header, True)
    header = rd(header, True)
    header = ra(header, True)
    header = z(header, 0x0a)
    header = rcode(header, 0x0a)
    header = qdcount(header, 0x0a0a)
    header = ancount(header, 0xf1f1)
    header = nscount(header, 0xdada)
    header = arcount(header, 0xfbf1)

    printHeader(header)
}
