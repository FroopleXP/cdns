package main

import (
    "fmt"
    "strings"
)

type rrecord struct {
    qnames []string
    qtype qtype
    qclass qclass
    ttl uint32
    rdlen uint16
    rdata []byte

    // srdata is the string (printable) representation of 'rdata'
    srdata string
}

func (r *rrecord) readFrom(b []byte, n int) int {
    s, ls := readLabels(b, n); n+=s
    r.qnames = ls

    r.qtype  = qtype(bytesToUint16(b[n], b[n+1])); n+=2
    r.qclass = qclass(bytesToUint16(b[n], b[n+1])); n+=2
    r.ttl    = bytesToUint32(b[n], b[n+1], b[n+2], b[n+3]); n+=4
    r.rdlen  = bytesToUint16(b[n], b[n+1]); n+=2
    r.rdata  = b[n:n+int(r.rdlen)]; n+=int(r.rdlen)

    if r.qtype == QT_A && r.rdlen == 4 {
        r.srdata = fmt.Sprintf("%d.%d.%d.%d", r.rdata[0], r.rdata[1], r.rdata[2], r.rdata[3])
    } else if r.qtype == QT_NS || r.qtype == QT_CNAME {
        _, labels := readLabels(b, n - int(r.rdlen))
        r.srdata = strings.Join(labels, ".")
    }

    return n
}

func (r *rrecord) String() string {
    return r.srdata
}
