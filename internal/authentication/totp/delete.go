package totp

import (
	"errors"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/database"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/sirupsen/logrus"
)

func DeleteTOTP(userDN string) error {
	var err error

	if configuration.Configuration.Features.DisableTOTP {
		logrus.Warnf("User %s tried to revoke TOTP, but this function is disabled.", userDN)
		return errors.New("totp is disabled on this server")
	}

	// Delete token from backend
	switch configuration.Configuration.TOTP.Kind {
	case "db":
		err = database.DeleteToken(encodeUser(userDN))
	case "openldap":
		err = ldap.DeleteToken(userDN)
	default:
		logrus.Fatal("unreachable! bad TOTP kind")
	}
	if err != nil {
		logrus.Warnf("Failed to revoke TOTP for user %s. Error was: %s.", userDN, err.Error())
		return err
	}

	return nil

}
