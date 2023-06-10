package totp

import (
	"errors"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/database"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/sirupsen/logrus"
	"github.com/thanhpk/randstr"
)
import "github.com/pquerna/otp/totp"

func CreateTOTP(userDN string, username string) (string, error) {
	var err error

	if configuration.Configuration.Features.DisableTOTP {
		logrus.Warnf("User %s tried to generate TOTP, but this function is disabled.", userDN)
		return "", errors.New("totp is disabled on this server")
	}

	// Create random security key and store it into backend
	secret := randstr.String(256)
	switch configuration.Configuration.TOTP.Kind {
	case "db":
		err = database.CreateToken(encodeUser(userDN), secret)
	case "openldap":
		err = ldap.CreateToken(userDN, secret)
	default:
		logrus.Fatal("unreachable! bad TOTP kind")
	}
	if err != nil {
		logrus.Warnf("Failed to generate TOTP for user %s. Error was: %s.", userDN, err.Error())
		return "", err
	}

	// Generate TOTP
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      configuration.Configuration.TOTP.CustomServiceName,
		AccountName: username,
		Secret:      []byte(secret),
		Period:      30,
	})
	if err != nil {
		logrus.Warnf("Failed to generate TOTP for user %s. Error was: %s.", username, err.Error())
		return "", err
	}

	return key.String(), nil

}
