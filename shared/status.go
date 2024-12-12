package shared

// IsLinked checks if the team is linked by verifying that both the TeamID and AuthToken
// in the shared configuration are not empty strings.
// It returns true if both values are present, indicating that the team is linked;
// otherwise, it returns false.
func IsLinked() bool {
	return Config.TeamID != "" && Config.AuthToken != ""
}
