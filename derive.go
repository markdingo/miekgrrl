/*
Package miekgrrl is a tiny helper package which creates a
[github.com/markdingo/rrl.ResponseTuple] from a [github.com/miekg/dns.Msg].

A “ResponseTuple” is passed to the [github.com/markdingo/rrl.Debit] function for
“Response Rate Limiting” analysis.
*/
package miekgrrl

import (
	"github.com/markdingo/rrl"
	"github.com/miekg/dns"
)

// Derive creates a [github.com/markdingo/rrl.ResponseTuple] from a [github.com/miekg/dns]
// response Msg for the purpose of passing to the [github.com/markdingo/rrl.Debit]
// “Response Rate Limiting” function.
//
// If the response has been formulated from a wildcard the caller *must* supply the
// origin name of the owning zone in the wildcardOriginName argument.
// Normally that will be qName with the first label removed but the following is also a
// valid zone:
//
//	$ORIGIN example.net.
//	*.a.b.c IN TXT "my origin name is example.net."
//
// thus one cannot blindly remove the first label and hope for the best.
func Derive(response *dns.Msg, wildcardOriginName string) (tuple *rrl.ResponseTuple) {

	tuple = &rrl.ResponseTuple{}

	// As of rfc7873#5.4 (circa 2016) it's legal to have a request (and thus response)
	// with an empty Question section which implies that there is no Class or Type.
	// On that basis, we leave the tuple fields as their default of zero which both
	// serendipitously were never allocated by rfc1035 or anything since.
	//
	// The net result of an invalid Class and Type is normally to categorize the
	// response as an error, but it probably makes more sense to categorize it as a
	// NoData response which is what this function does.
	// Mind you, in over a year of continuous monitoring, I've not seen a single
	// genuine instance of such a query so it's hardly likely to be a big issue.

	if len(response.Question) > 0 { // Can only set if there is a question
		tuple.Class = response.Question[0].Qclass
		tuple.Type = response.Question[0].Qtype
	}

	// Each category has a separate, configurable per-second allowance

	tuple.AllowanceCategory = rrl.NewAllowanceCategory(response.Rcode, len(response.Answer), len(response.Ns))

	// AllowanceCategory influences SalientName

	if tuple.AllowanceCategory == rrl.AllowanceNXDomain || tuple.AllowanceCategory == rrl.AllowanceReferral {
		if len(response.Ns) > 0 {
			tuple.SalientName = response.Ns[0].Header().Name
		}
	} else if len(response.Question) > 0 { // At least get something in there
		if len(wildcardOriginName) > 0 {
			tuple.SalientName = wildcardOriginName
		} else {
			tuple.SalientName = response.Question[0].Name
		}
	}

	return
}
