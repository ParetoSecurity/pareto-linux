package shared

func IsValidParetoURL(url string) bool {
	return len(url) > 18 && url[:18] == "paretosecurity://"
}
