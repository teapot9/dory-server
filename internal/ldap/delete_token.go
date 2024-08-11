package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap/helpers"
	"github.com/be-ys-cloud/dory-server/internal/structures"
)

func DeleteToken(userDN string) error {
	log := logrus.WithField("service", "DeleteToken")
	log.Debugf("removing openldap totp token for %s", userDN)

	// Connect to LDAP server
	l, err := helpers.GetSession(
		configuration.Configuration.LDAPServer.Address,
		configuration.Configuration.LDAPServer.Port,
		configuration.Configuration.LDAPServer.SkipTLSVerify,
	)
	if err != nil {
		msg := fmt.Sprintf("could not connect to server: %v", err)
		log.Warn(msg)
		return &structures.CustomError{Text: msg, HttpCode: 504}
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
		return &structures.CustomError{Text: msg, HttpCode: 500}
	}

	// Check if oathSecret present
	user, err := helpers.GetUserByDN(l, userDN)
	if err != nil {
		log.Warnf("could not find user %s: %v", userDN, err)
		return err
	}
	if len(helpers.FindAttribute(user, "oathTOTPToken")) == 0 {
		msg := fmt.Sprintf("TOTP is not configured for %s", userDN)
		log.Info(msg)
		return &structures.CustomError{Text: msg, HttpCode: 404}
	} else {
		log.Debug("found the oathSecret attribute")
	}

	// Build LDAP request
	req := ldap.NewModifyRequest(userDN, nil)
	req.Delete(
		"objectClass",
		[]string{"oathToken", "oathTOTPToken", "oathTOTPUser", "oathUser"},
	)
	req.Replace("oathTOTPToken", nil)
	req.Replace("oathTOTPParams", nil)
	req.Replace("oathSecret", nil)
	req.Replace("oathTOTPLastTimeStep", nil)
	req.Replace("oathTOTPTimeStepDrift", nil)

	err = l.Modify(req)
	if err != nil {
		msg := fmt.Sprintf("failed to delete token from LDAP: %v", err)
		log.Error(msg)
		return &structures.CustomError{Text: msg, HttpCode: 500}
	}

	return nil
}
