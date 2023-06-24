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

func TestRevokeFailOpenLDAP(t *testing.T) {
	msg := "should have failed to revoke TOTP"
	url := baseUrl + "totp/revoke"

	config.TOTP.Kind = "openldap"
	helpers.ReloadServerConfig(t, &config, &containers)

	defer func() {
		config.TOTP.Kind = "db"
		helpers.ReloadServerConfig(t, &config, &containers)
	}()

	t.Run("revoke with bad password", func(t *testing.T) {
		data := user{
			Username: "otpuser-enabled",
			Authentication: authentication{
				Password: "badpass",
			},
		}

		marshaled, _ := json.Marshal(data)
		reader := strings.NewReader(string(marshaled))

		code, resp, _, err := connectors.WSProvider("POST", url, reader, nil)
		helpers.AssertHTTPResponse(t, msg, data, resp, err, 401, code)
	})

	t.Run("revoke when already disabled", func(t *testing.T) {
		data := user{
			Username: "otpuser-disabled",
			Authentication: authentication{
				Password: "test",
			},
		}

		marshaled, _ := json.Marshal(data)
		reader := strings.NewReader(string(marshaled))

		code, resp, _, err := connectors.WSProvider("POST", url, reader, nil)
		helpers.AssertHTTPResponse(t, msg, data, resp, err, 404, code)
	})

	t.Run("revoke with OTP authentication", func(t *testing.T) {
		totpcode, err := totp.GenerateCode(helpers.EncodeTOTP("totpsecret"), time.Now())
		if err != nil {
			t.Fatalf("failed to get TOTP code: %v", err)
		}

		data := user{
			Username: "otpuser-enabled",
			Authentication: authentication{
				TOTP: totpcode,
			},
		}

		marshaled, _ := json.Marshal(data)
		reader := strings.NewReader(string(marshaled))

		code, resp, _, err := connectors.WSProvider("POST", url, reader, nil)
		helpers.AssertHTTPResponse(t, msg, data, resp, err, 401, code)
	})
}
