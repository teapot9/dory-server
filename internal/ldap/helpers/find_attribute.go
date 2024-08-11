package helpers

import (
	"github.com/go-ldap/ldap"
)

func FindAttribute(entry *ldap.Entry, attribute string) []string {
	for _, attr := range entry.Attributes {
		if attr.Name == attribute {
			return attr.Values
		}
	}
	return nil
}
