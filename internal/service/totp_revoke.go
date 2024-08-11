package service

import (
	"github.com/be-ys-cloud/dory-server/internal/authentication/totp"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/be-ys-cloud/dory-server/internal/mailer"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/sirupsen/logrus"
)

func RevokeTOTP(user structures.UserCreateTOTP) error {

	// Check authentication
	valid, err := checkAuth(user.Username, user.Authentication)
	if err != nil {
		return &structures.CustomError{HttpCode: 401, Text: err.Error()}
	}
	if !valid {
		return &structures.CustomError{HttpCode: 401, Text: "invalid authentication"}
	}

	userDN, err := ldap.GetUserDN(user.Username)
	if err != nil {
		return err
	}

	err = totp.DeleteTOTP(userDN)
	if err != nil {
		return err
	}

	// Send email
	email, err := ldap.GetUserMail(user.Username)

	if err != nil {
		logrus.Warnf("Could not send totp deleted mail to user %s because there is no email associated to it on Active Directory.", user.Username)
	} else {
		_ = mailer.SendMail("mail_totp_deleted", email, struct {
			Name string
			URL  string
			LDAP string
		}{Name: user.Username, URL: configuration.Configuration.FrontAddress, LDAP: configuration.Configuration.LDAPServer.Address})
	}

	return err
}
