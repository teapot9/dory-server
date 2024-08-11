package tests

import (
	"encoding/json"
	"strings"
	"testing"
	"github.com/be-ys-cloud/dory-server/test/connectors"
	"github.com/be-ys-cloud/dory-server/test/helpers"
)

func TestVerifyFailOpenLDAP(t *testing.T) {
	msg := "should have failed to verify TOTP"
	url := baseUrl + "totp/verify"

	config.TOTP.Kind = "openldap"
	helpers.ReloadServerConfig(t, &config, &containers)

	defer func() {
		config.TOTP.Kind = "db"
		helpers.ReloadServerConfig(t, &config, &containers)
	}()

	t.Run("verify with bad code", func(t *testing.T) {
		data := user{
			Username: "otpuser-enabled",
			TOTP:     "000000",
		}

		marshaled, _ := json.Marshal(data)
		reader := strings.NewReader(string(marshaled))

		code, resp, _, err := connectors.WSProvider("POST", url, reader, nil)
		helpers.AssertHTTPResponse(t, msg, data, resp, err, 401, code)
	})

	t.Run("verify when disabled", func(t *testing.T) {
		data := user{
			Username: "otpuser-disabled",
			TOTP:     "000000",
		}

		marshaled, _ := json.Marshal(data)
		reader := strings.NewReader(string(marshaled))

		code, resp, _, err := connectors.WSProvider("POST", url, reader, nil)
		helpers.AssertHTTPResponse(t, msg, data, resp, err, 404, code)
	})
}
