package file_nsec5

import (
	"github.com/coredns/coredns/plugin/file_nsec5/tree"

	"github.com/miekg/dns"
)

// ClosestEncloser returns the closest encloser for qname.
func (z *Zone) ClosestEncloser(qname string) (*tree.Elem, bool) {

	offset, end := dns.NextLabel(qname, 0)
	for !end {
		elem, _ := z.Tree.Search(qname)
		if elem != nil {
			return elem, true
		}
		qname = qname[offset:]

		offset, end = dns.NextLabel(qname, 0)
	}

	return z.Tree.Search(z.origin)
}

// NextClosestEncloser returns the next-closest encloser for qname.
func (z *Zone) NextClosestEncloser(qname string) string {
	nex := qname

	offset, end := dns.NextLabel(qname, 0)
	last_offset := 0
	for !end {
		elem, _ := z.Tree.Search(qname)
		if elem != nil {
			return nex
		}
		qname = qname[offset:]
		nex = nex[last_offset:]
		last_offset = offset
		offset, end = dns.NextLabel(qname, offset)
	}

	return nex
}
