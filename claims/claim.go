package claims

import "github.com/ParetoSecurity/pareto-core/check"

type Claim struct {
	Title  string
	Checks []check.Check
}
