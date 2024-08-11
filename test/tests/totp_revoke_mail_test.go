package tests

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/test/connectors"
	"strings"
	"testing"
)

// TestTOTPMail generates a TOTP token and revokes it with mail.
func TestTOTPMail(t *testing.T) {

	// Trying to create TOTP
	url := baseUrl + "totp/create"

	data := user{
		Username: "testuser",
		Authentication: authentication{
			Password: "test",
		},
	}

	marshaled, _ := json.Marshal(data)
	reader := strings.NewReader(string(marshaled))

	code, response, _, err := connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	var TOTP totpStruct
	err = json.Unmarshal(response, &TOTP)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	TOTP.TOTP = strings.Split(TOTP.TOTP, "secret=")[1]
	TOTP.TOTP = strings.Split(TOTP.TOTP, "&")[0]
	t.Log(TOTP.TOTP)

	// Send request for valid user
	url = baseUrl + "request/revoke"

	data = user{
		Username: "testuser",
	}

	marshaled, _ = json.Marshal(data)
	reader = strings.NewReader(string(marshaled))

	code, _, _, err = connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	// Check mail have been sent and retrieve link
	code, response, _, err = connectors.WSProvider("GET", mailUrl, nil, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	var mails []email
	_ = json.Unmarshal(response, &mails)

	resetToken := strings.Split(strings.Split(mails[0].TextAsHtml, "https://localhost:8001/revoke?user=testuser&amp;type=mail&amp;token=")[1], "\"")[0]

	// Revoking TOTP
	url = baseUrl + "totp/revoke"

	data = user{
		Username: "testuser",
		Authentication: authentication{
			Token: resetToken,
		},
	}

	marshaled, _ = json.Marshal(data)
	reader = strings.NewReader(string(marshaled))

	code, _, _, err = connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}
}
