package service

import (
	"github.com/be-ys-cloud/dory-server/internal/authentication/token"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/ldap"
	"github.com/be-ys-cloud/dory-server/internal/mailer"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/sirupsen/logrus"
)

func ReinitializePassword(user structures.UserReinitialize) error {
	//Check that token is valid
	valid, err := checkAuth(user.Username, user.Authentication)
	if err != nil {
		return &structures.CustomError{HttpCode: 401, Text: err.Error()}
	}
	if !valid {
		return &structures.CustomError{HttpCode: 401, Text: "invalid authentication"}
	}

	//Token is valid, now modifying password !
	err = ldap.ReinitializePassword(user.Username, user.NewPassword)
	if err != nil {
		logrus.Warn("Error while reinitializing password for user %s. Error was: %s", user.Username, err.Error())
		return &structures.CustomError{Text: "an error occurred while reinitializing password", HttpCode: 500}
	}

	//Removing key
	token.DeleteKey(user.Username)

	//Modification done, sending mail
	email, err := ldap.GetUserMail(user.Username)

	if err != nil {
		logrus.Warnf("Could not send password changed mail to user %s because there is no email associated to it on Active Directory.", user.Username)
	} else {
		_ = mailer.SendMail("mail_info_changed", email, struct {
			Name string
			URL  string
			LDAP string
		}{Name: user.Username, URL: configuration.Configuration.FrontAddress, LDAP: configuration.Configuration.LDAPServer.Address})
	}

	return nil
}
