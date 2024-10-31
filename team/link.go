package team

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
)

const baseURL = "https://dash.paretosecurity.com/api/v1/team"

func LinkWithDevice(device ReportingDevice, teamAuth string, teamID string) error {
	url := fmt.Sprintf("%s/%s/device", baseURL, teamID)
	ctx := context.Background()
	var res interface{}

	err := requests.
		URL(url).
		BodyJSON(&device).
		Header("X-Device-Auth", teamAuth).
		ToJSON(&res).
		Fetch(ctx)

	if err != nil {
		return err
	}

	return nil
}
