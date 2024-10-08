package helpers

import (
	"fmt"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/go-ldap/ldap"
	"github.com/sirupsen/logrus"
)

func searchUser(l *ldap.Conn, req *ldap.SearchRequest) (*ldap.Entry, error) {
	sr, err := l.Search(req)
	if err != nil {
		logrus.Warnln("Unable to search into LDAP. Detailed error : " + err.Error())
		return nil, &structures.CustomError{Text: "could not search into LDAP", HttpCode: 503}
	}

	if len(sr.Entries) == 0 {
		logrus.Warnln("No user matched.")
		return nil, &structures.CustomError{Text: "user not found in LDAP", HttpCode: 404}
	}

	if len(sr.Entries) > 1 {
		logrus.Warnln("Too many user matched.")
		return nil, &structures.CustomError{Text: "too many user matched LDAP. could not process to avoid modifying one undesired account", HttpCode: 404}
	}

	return sr.Entries[0], nil
}

func GetUser(l *ldap.Conn, baseDN string, filterOn string, username string) (*ldap.Entry, error) {

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(filterOn, ldap.EscapeFilter(username)),
		[]string{},
		nil,
	)

	return searchUser(l, searchRequest)
}

func GetUserByDN(l *ldap.Conn, dn string) (*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(
		dn,
		ldap.ScopeBaseObject, ldap.NeverDerefAliases, 1, 0, false,
		"(objectClass=*)", nil, nil,
	)

	return searchUser(l, searchRequest)
}
