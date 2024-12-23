package claims

import "github.com/ParetoSecurity/pareto-linux/check"

type Claim struct {
	Title  string
	Checks []check.Check
}
