package main

import (
	"log"
	"net"
    "flag"
    "fmt"
    "strings"
)

func printPacket(pk *packet) {
    fmt.Printf("QUESTIONS\n") 
    for _, q := range pk.questions {
        domain := strings.Join(q.qnames, ".")
        
        fmt.Printf("%s\t", domain)
        if q.qtype == QT_A {
            fmt.Printf("A\t")
        } else if q.qtype == QT_CNAME {
            fmt.Printf("CNAME\t")
        }
        fmt.Printf("(%s)\n", q.qtype)
    }
    fmt.Printf("\n")

    fmt.Printf("ANSWERS\n") 
    for _, a := range pk.answers {
        domain := strings.Join(a.qnames, ".")
        
        fmt.Printf("%s\t", domain)
        if a.qtype == QT_A {
            fmt.Printf("A\t")
        } else if a.qtype == QT_CNAME {
            fmt.Printf("CNAME\t")
        }
        fmt.Printf("(%s)\t", a.qtype)
        fmt.Printf("%s\n", a)
    }
}

func main() {
    fDNSServerAddr := flag.String("s", "8.8.8.8:53", "dns server address to query")
    flag.Parse()

    if flag.NArg() < 2 {
        log.Fatal("insufficient arguments\n")
    }

    var host string = flag.Arg(0)
    var typ string = flag.Arg(1)

	var pk []byte = make([]byte, 0)

	var h *header = &header{0x0001, RD, 0x0001, 0x0000, 0x0000, 0x0000}
    pk = h.writeTo(pk)

    var q *question = &question{[]string{host}, QT_A, QC_IN}

    switch typ {
    case "NS":
        q.qtype = QT_NS
    case "CNAME":
        q.qtype = QT_CNAME
    default:
        q.qtype = QT_A   
    }

    pk = q.writeTo(pk)

	addr, err := net.ResolveUDPAddr("udp", *fDNSServerAddr)
	if err != nil {
		log.Fatalf("failed to resolve udp address: %v\n", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("failed to open connection: %v\n", err)
	}
	defer conn.Close()

	_, err = conn.Write(pk)
	if err != nil {
		log.Fatalf("failed to write packet: %v\n", err)
	}

	buf := make([]byte, 512)
    _, err = conn.Read(buf)
    if err != nil {
        log.Fatalf("failed to read from server: %v\n", err)
    }

    var p *packet = &packet{}
    p.read(buf)

    var e error = p.header.flags.Error()
    if e != nil {
        log.Fatalf("request failed: %v\n", e)
    }

    printPacket(p)
}

