package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap/helpers"
	"github.com/be-ys-cloud/dory-server/internal/structures"
)

func CreateToken(userDN string, key string) error {
	log := logrus.WithField("service", "CreateToken")
	log.Debugf("creating openldap totp token for %s", userDN)

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

	// Build LDAP request
	req := ldap.NewModifyRequest(userDN, nil)
	req.Add(
		"objectClass",
		[]string{"oathToken", "oathTOTPToken", "oathTOTPUser", "oathUser"},
	)
	req.Replace("oathTOTPToken", []string{userDN})
	req.Replace("oathTOTPParams", []string{configuration.Configuration.TOTP.OpenLDAPParamsDN})
	req.Replace("oathSecret", []string{key})

	err = l.Modify(req)
	if err != nil && ldap.IsErrorWithCode(err, ldap.LDAPResultAttributeOrValueExists) {
		msg := "otp already enabled"
		log.Info(msg)
		return &structures.CustomError{Text: msg, HttpCode: 409}
	} else if err != nil {
		msg := fmt.Sprintf("failed to send token to LDAP: %v", err)
		log.Error(msg)
		return &structures.CustomError{Text: msg, HttpCode: 500}
	}

	return nil
}
