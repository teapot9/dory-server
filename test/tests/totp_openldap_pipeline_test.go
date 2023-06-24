package tests

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
	"github.com/pquerna/otp/totp"
	"github.com/be-ys-cloud/dory-server/test/connectors"
	"github.com/be-ys-cloud/dory-server/test/helpers"
)

// TestTOTPOpenLDAP generates a TOTP token, verifies it and revoke it
func TestTOTPOpenLDAP(t *testing.T) {
	config.TOTP.Kind = "openldap"
	helpers.ReloadServerConfig(t, &config, &containers)

	defer func() {
		config.TOTP.Kind = "db"
		helpers.ReloadServerConfig(t, &config, &containers)
	}()

	// Trying to create TOTP
	url := baseUrl + "totp/create"

	data := user{
		Username: "otpuser-disabled",
		Authentication: authentication{
			Password: "test",
		},
	}

	marshaled, _ := json.Marshal(data)
	reader := strings.NewReader(string(marshaled))

	code, response, _, err := connectors.WSProvider("POST", url, reader, nil)
	helpers.AssertHTTPResponse(t, "could not create TOTP", data, response, err, 200, code)

	var TOTP totpStruct
	err = json.Unmarshal(response, &TOTP)
	if err != nil {
		t.Fatalf("failed to parse received TOTP struct: %v", err)
	}

	TOTP.TOTP = strings.Split(TOTP.TOTP, "secret=")[1]
	TOTP.TOTP = strings.Split(TOTP.TOTP, "&")[0]
	t.Log(TOTP.TOTP)

	// Verifying TOTP
	totpcode, err := totp.GenerateCode(TOTP.TOTP, time.Now())
	if err != nil {
		t.Fatalf("failed to generate TOTP code (TOTP: %v): %v", TOTP, err)
	}

	url = baseUrl + "totp/verify"

	data = user{
		Username: "otpuser-disabled",
		TOTP:     totpcode,
	}

	marshaled, _ = json.Marshal(data)
	reader = strings.NewReader(string(marshaled))

	code, response, _, err = connectors.WSProvider("POST", url, reader, nil)
	helpers.AssertHTTPResponse(t, "failed to verify TOTP", data, response, err, 200, code)

	// Revoking TOTP
	url = baseUrl + "totp/revoke"

	data = user{
		Username: "otpuser-disabled",
		Authentication: authentication{
			Password: "test" + totpcode,
		},
	}

	marshaled, _ = json.Marshal(data)
	reader = strings.NewReader(string(marshaled))

	code, response, _, err = connectors.WSProvider("POST", url, reader, nil)
	helpers.AssertHTTPResponse(t, "failed to revoke TOTP", data, response, err, 200, code)
}
