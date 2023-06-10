package helpers

import (
	"github.com/ory/dockertest/v3"
)

type ContainersEnvironment struct {
	Pool *dockertest.Pool
	Network *dockertest.Network
	LDAP *dockertest.Resource
	Mail *dockertest.Resource
	Server *dockertest.Resource
}
