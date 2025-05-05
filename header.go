package main

import (
    "fmt"
)

// RCODE
type rcode uint16

const (
	RC_NOERR  rcode = 0b0000000000000000
	RC_FMTERR rcode = 0b0000000000000001
	RC_SRVERR rcode = 0b0000000000000010
	RC_NAMERR rcode = 0b0000000000000011
	RC_NOTIMP rcode = 0b0000000000000100
	RC_REFUSD rcode = 0b0000000000000101
)

func (r rcode) String() string {
    switch r {
    case RC_NOERR:
        return "No error"
    case RC_FMTERR:
        return "Format error"
    case RC_SRVERR:
        return "Server failure"
    case RC_NAMERR:
        return "Name error"
    case RC_NOTIMP:
        return "Not implemented"
    case RC_REFUSD:
        return "Refused"
    }
    return ""
}

// Flags
type flags uint16

const (
	R  flags = 0b1000000000000000
	AA flags = 0b0000010000000000
	TC flags = 0b0000001000000000
	RD flags = 0b0000000100000000
	RA flags = 0b0000000010000000
)

func (f flags) Error() error {
    if rcode(f & 0x000f) != RC_NOERR {
        return fmt.Errorf("server responded with error '%s'", rcode(f & 0x000f).String())
    }
    return nil
}

func (f flags) OpCode() opcode {
    return opcode(f & 0x7800)
}

func (f flags) IsRequest() bool {
    return (f & R) == R
}

func (f flags) IsAuthoritative() bool {
    return (f & AA) == AA
}

func (f flags) IsTruncated() bool {
    return (f & TC) == TC
}

func (f flags) IsRecursionDesired() bool {
    return (f & RD) == RD
}

func (f flags) IsRecursionAvailable() bool {
    return (f & RA) == RA
}

// OpCode
type opcode uint16

const (
	OP_QUERY  opcode = 0b0000000000000000
	OP_IQUERY opcode = 0b0000100000000000
	OP_STATUS opcode = 0b0001000000000000
)

func (o opcode) String() string {
    switch o {
    case OP_QUERY:
        return "Standard query"
    case OP_IQUERY:
        return "Inverse query"
    case OP_STATUS:
        return "Server status request"
    }
    return ""
}

type header struct {
    id uint16
    flags flags
    qdcount uint16
    ancount uint16
    nscount uint16
    arcount uint16
}

func (h *header) writeTo(b []byte) []byte {
	b = writeSplitUint16(b, h.id)
	b = writeSplitUint16(b, uint16(h.flags))
	b = writeSplitUint16(b, h.qdcount)
	b = writeSplitUint16(b, h.ancount)
	b = writeSplitUint16(b, h.nscount)
	b = writeSplitUint16(b, h.arcount)
	return b
}

func (h *header) readFrom(b []byte, n int) int {
    h.id = bytesToUint16(b[0], b[1]); n+=2
    h.flags = flags(bytesToUint16(b[2], b[3])); n+=2
    h.qdcount = bytesToUint16(b[4], b[5]); n+=2
    h.ancount = bytesToUint16(b[6], b[7]); n+=2
    h.nscount = bytesToUint16(b[8], b[9]); n+=2
    h.arcount = bytesToUint16(b[10], b[11]); n+=2
    return n
}

