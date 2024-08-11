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

func TestFailCreateTOTPOpenLDAP(t *testing.T) {
	msg := "should have failed to create TOTP"
	url := baseUrl + "totp/create"

	config.TOTP.Kind = "openldap"
	helpers.ReloadServerConfig(t, &config, &containers)

	defer func() {
		config.TOTP.Kind = "db"
		helpers.ReloadServerConfig(t, &config, &containers)
	}()

	t.Run("create with bad password", func(t *testing.T) {
		data := user{
			Username: "otpuser-disabled",
			Authentication: authentication{
				Password: "badpassword",
			},
		}

		marshaled, _ := json.Marshal(data)
		reader := strings.NewReader(string(marshaled))

		code, resp, _, err := connectors.WSProvider("POST", url, reader, nil)
		helpers.AssertHTTPResponse(t, msg, data, resp, err, 401, code)
	})

	t.Run("create when already enabled", func(t *testing.T) {
		totpcode, err := totp.GenerateCode(helpers.EncodeTOTP("totpsecret"), time.Now())
		if err != nil {
			t.Fatalf("failed to get TOTP code: %v", err)
		}

		data := user{
			Username: "otpuser-enabled",
			Authentication: authentication{
				Password: "test" + totpcode,
			},
		}

		marshaled, _ := json.Marshal(data)
		reader := strings.NewReader(string(marshaled))

		code, resp, _, err := connectors.WSProvider("POST", url, reader, nil)
		helpers.AssertHTTPResponse(t, msg, data, resp, err, 409, code)
	})
}
