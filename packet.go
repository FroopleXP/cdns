package main

type packet struct {
    header    *header
    questions []*question
    answers   []*rrecord
}

func (p *packet) read(b []byte) {
    var ptr int = 0

    // Header
    p.header = &header{}
    ptr = p.header.readFrom(b, ptr)

    // Read Question Section
    for range p.header.qdcount {
        var q *question = &question{}
        ptr = q.readFrom(b, ptr)
        p.questions = append(p.questions, q)
    }

    // Answer section
    for range p.header.ancount {
        var rr *rrecord = &rrecord{}
        ptr = rr.readFrom(b, ptr)
        p.answers = append(p.answers, rr)
    }
}

