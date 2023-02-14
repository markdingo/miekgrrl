package miekgrrl

import (
	"testing"

	"github.com/markdingo/rrl"
	"github.com/miekg/dns"
)

// addr implements a mock net.Addr
type addr struct {
	n, s string
}

func (a *addr) Network() string { return a.n }
func (a *addr) String() string  { return a.s }

func newAddr(n, s string) *addr {
	return &addr{n: n, s: s}
}

func TestDerive(t *testing.T) {
	type testCase struct {
		question       string
		answer         string
		ns             string
		rcode          int
		wildOriginName string
		ac             rrl.AllowanceCategory
		sn             string
	}

	testCases := []testCase{
		{"", "", "", 0, "", rrl.AllowanceNoData, ""},
		{"example.net.", "example.net. IN A 10.0.0.1", "", 0, "", rrl.AllowanceAnswer, "example.net."},

		// Referrals should return the authoritative domain rather than the qName
		{"example.net.", "", "subexample.net. IN NS ns1.example.net.", 0, "", rrl.AllowanceReferral,
			"subexample.net."},
		{"example.net.", "", "subexample.net. IN NS ns1.example.net.", 3, "", rrl.AllowanceNXDomain,
			"subexample.net."},

		{"example.net.", "example.net. IN A 10.0.0.1", "", 0, "net.", rrl.AllowanceAnswer, "net."},
	}

	for ix, tc := range testCases {
		var msg dns.Msg
		if len(tc.question) > 0 {
			msg.SetQuestion(tc.question, dns.TypeAAAA)
		}
		if len(tc.answer) > 0 {
			a, err := dns.NewRR(tc.answer)
			if err != nil {
				t.Fatal(ix, "Setup of test failed", err)
			}
			msg.Answer = append(msg.Answer, a)
		}
		if len(tc.ns) > 0 {
			ns, err := dns.NewRR(tc.ns)
			if err != nil {
				t.Fatal(ix, "Setup of test failed", err)
			}
			msg.Ns = append(msg.Ns, ns)
		}
		msg.MsgHdr.Rcode = tc.rcode

		tuple := Derive(&msg, tc.wildOriginName)
		if tuple.AllowanceCategory != tc.ac {
			t.Error(ix, "Category mismatch. Expected", tc.ac, "Got", tuple.AllowanceCategory)
		}
		if tuple.SalientName != tc.sn {
			t.Error(ix, "SalientName mismatch. Expected", tc.sn, "Got", tuple.SalientName)
		}
	}
}
