package ldap

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap/helpers"
	"github.com/be-ys-cloud/dory-server/internal/structures"
)

func GetToken(userDN string) (string, error) {
	log := logrus.WithField("service", "GetToken")
	log.Debugf("getting openldap totp token for %s", userDN)

	// Connect to LDAP server
	l, err := helpers.GetSession(
		configuration.Configuration.LDAPServer.Address,
		configuration.Configuration.LDAPServer.Port,
		configuration.Configuration.LDAPServer.SkipTLSVerify,
	)
	if err != nil {
		msg := fmt.Sprintf("could not connect to server: %v", err)
		log.Warn(msg)
		return "", &structures.CustomError{Text: msg, HttpCode: 504}
	}
	defer l.Close()

	// Bind as admin
	err = helpers.BindUser(
		l, configuration.Configuration.LDAPServer.Admin.Username,
		configuration.Configuration.LDAPServer.Admin.Password,
	)
	if err != nil { // Probably misconfiguration
		msg := fmt.Sprintf("could not login to LDAP as admin: %v", err)
		log.Warn(msg)
		return "", &structures.CustomError{Text: msg, HttpCode: 500}
	}

	// Get user info
	user, err := helpers.GetUserByDN(l, userDN)
	if err != nil {
		log.Warnf("could not find user %s: %v", userDN, err)
		return "", err
	}

	// Get oathTOTPToken
	values := helpers.FindAttribute(user, "oathSecret")
	if len(values) == 0 {
		msg := fmt.Sprintf("oathSecret attribute not found for %s", userDN)
		log.Error(msg)
		return "", &structures.CustomError{Text: msg, HttpCode: 404}
	}
	return values[0], nil
}
