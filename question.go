package main

// QTYPE
type qtype uint16

const (
	QT_A = qtype(uint16(iota + 1))
    QT_NS
    QT_MD
    QT_MF
    QT_CNAME
)

func (q qtype) String() string {
    switch q {
    case QT_A:
        return "host address"
    case QT_NS:
        return "unauthoritative name server"
    case QT_MD:
        return "mail destination"
    case QT_MF:
        return "mail forwarder"
    case QT_CNAME:
        return "canonical name for an alias"
    }
    return ""
}

// QCLASS
type qclass uint16

// NOTE: There are many more defined in the RFC but I'm only
// implementing a small subset of them
const (
	QC_IN = qclass(uint16(iota + 1))
)

func (q qclass) String() string {
    switch q {
    case QC_IN:
        return "internet"
    }
    return ""
}

type question struct {
    qnames []string
    qtype  qtype
    qclass qclass
}

func (q *question) writeTo(b []byte) []byte {
	b = writeQNames(b, q.qnames)
	b = writeSplitUint16(b, uint16(q.qtype))
	b = writeSplitUint16(b, uint16(q.qclass))
	return b
}

func (q *question) readFrom(b []byte, n int) int {
    s, ls := readLabels(b, n); n+=s
    q.qnames = ls

    q.qtype = qtype(bytesToUint16(b[n], b[n+1])); n+=2
    q.qclass = qclass(bytesToUint16(b[n], b[n+1])); n+=2
    return n
}
