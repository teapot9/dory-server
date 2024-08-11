package totp

import (
	"encoding/base32"
	"errors"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/database"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/sirupsen/logrus"
	"github.com/pquerna/otp/totp"
)

func VerifyTOTP(userDN string, token string) (bool, error) {
	var err error
	if configuration.Configuration.Features.DisableTOTP {
		logrus.Warnf("User %s tried to verify TOTP, but this function is disabled.", userDN)
		return false, errors.New("totp is disabled on this server")
	}

	// Get token from backend
	var tokenUser string
	switch configuration.Configuration.TOTP.Kind {
	case "db":
		tokenUser, err = database.GetToken(encodeUser(userDN))
	case "openldap":
		tokenUser, err = ldap.GetToken(userDN)
	default:
		logrus.Fatal("unreachable! bad TOTP kind")
	}
	if err != nil {
		logrus.Warnf("Failed to generate TOTP for user %s. Error was: %s.", userDN, err.Error())
		return false, err
	}

	// Return validation state
	return totp.Validate(token, base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(tokenUser))), nil
}
