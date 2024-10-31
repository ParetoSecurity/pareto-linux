package claims

import "paretosecurity.com/auditor/check"

type Claim struct {
	Title  string
	Checks []check.Check
}
