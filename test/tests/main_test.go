package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/be-ys-cloud/dory-server/test/helpers"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	_ "github.com/go-sql-driver/mysql"
)

var baseUrl string
var mailUrl string
var config structures.Configuration
var containers helpers.ContainersEnvironment

// TestMain contains only the minimum required to start test suite.
func TestMain(m *testing.M) {
	// Generate stack
	err := setupEnv()
	if err != nil {
		log.Fatalln(err)
	}

	// Set variables that will be used by other test files !
	baseUrl = "http://127.0.0.1:" + containers.Server.GetPort("8000/tcp") + "/"
	mailUrl = "http://127.0.0.1:" + containers.Mail.GetPort("1080/tcp") + "/api/emails"

	// Run tests
	_ = m.Run()

	// Destroy stack
	err = destroyEnv(containers.Pool, containers.Network, containers.LDAP, containers.Mail, containers.Server)
	if err != nil {
		log.Fatalln("Could not destroy stack")
	}
}

//--------------- All-in-one methods

func setupEnv() (err error) {

	containers.Pool, err = createPool()
	if err != nil {
		return
	}

	containers.Network, err = containers.Pool.CreateNetwork("dory-tests")
	if err != nil {
		return
	}

	containers.LDAP, err = createOpenLDAPContainer(containers.Pool, containers.Network)
	if err != nil {
		return
	}

	containers.Mail, err = createMailContainer(containers.Pool, containers.Network)
	if err != nil {
		return
	}

	containers.Server, err = createServerContainer(
		containers.Pool, containers.Network,
		containers.LDAP.GetPort("636/tcp"),
		containers.Mail.GetPort("1025/tcp"),
	)
	if err != nil {
		return
	}

	// Wait for servers to be up...
	time.Sleep(10 * time.Second)

	return
}

func destroyEnv(pool *dockertest.Pool, network *dockertest.Network, ldap *dockertest.Resource, mail *dockertest.Resource, server *dockertest.Resource) (err error) {

	if ldap != nil {
		if err = deleteContainer(pool, ldap); err != nil {
			return
		}
	}

	if server != nil {
		if err = deleteContainer(pool, server); err != nil {
			return
		}
	}

	if mail != nil {
		if err = deleteContainer(pool, mail); err != nil {
			return
		}
	}

	_ = pool.RemoveNetwork(network)

	path, err := os.Getwd()
	if err != nil {
		return
	}

	_ = exec.Command("rm", "-f", path+"/configuration.json").Run()
	_ = exec.Command("docker", "image", "rm", "-f", "dory_base_test").Run()

	return
}

//--------------- Unit methods

func createPool() (*dockertest.Pool, error) {
	return dockertest.NewPool("")
}

func createOpenLDAPContainer(pool *dockertest.Pool, network *dockertest.Network) (ressource *dockertest.Resource, err error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path = strings.TrimSuffix(path, "/tests")
	dockerfile := path + "/ldap_data/Dockerfile"

	ressource, err = pool.BuildAndRunWithOptions(dockerfile, &dockertest.RunOptions{
		Name:       "dory_test_openldap",
		Networks:   []*dockertest.Network{network},
		Mounts: []string{
			path + "/ldap_data/bootstrap:/bootstrap:ro",
		},
	})
	return
}

func createMailContainer(pool *dockertest.Pool, network *dockertest.Network) (ressource *dockertest.Resource, err error) {
	ressource, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "dory_test_mailserver",
		Repository: "reachfive/fake-smtp-server",
		Tag:        "latest",
		Networks:   []*dockertest.Network{network},
	})
	return
}

func createServerContainer(pool *dockertest.Pool, network *dockertest.Network, ldapPort string, mailPort string) (*dockertest.Resource, error) {
	ldapPortInt, _ := strconv.Atoi(ldapPort)
	mailPortInt, _ := strconv.Atoi(mailPort)

	config = structures.Configuration{
		LDAPServer: structures.LDAPServerConfig{
			Admin: structures.LDAPServerAdmin{
				Username: "cn=admin,dc=localhost,dc=priv",
				Password: "admin",
			},
			BaseDN:        "dc=localhost,dc=priv",
			FilterOn:      "(&(objectClass=person)(cn=%s))",
			Address:       "host.docker.internal",
			Port:          ldapPortInt,
			Kind:          "openldap",
			SkipTLSVerify: true,
			EmailField:    "email",
		},
		TOTP: structures.TOTPConfig{
			Kind: "db",
			Secret: "AZERTYUIOPQSDFGHJKLMWXCVBN0123456789!",
			OpenLDAPParamsDN: "cn=otp,dc=localhost,dc=priv",
		},
		MailServer: structures.MailServerConfig{
			Address:       "host.docker.internal",
			Port:          mailPortInt,
			Password:      "",
			SenderAddress: "noreply@dory.localhost",
			SenderName:    "Dory",
			Subject:       "LDAP Account Management",
			SkipTLSVerify: true,
			TLSMode:       structures.TLSModeNone,
		},
		FrontAddress: "https://localhost:8001/",
	}

	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(config)
	err = os.WriteFile(path+"/configuration.json", data, 0777)
	if err != nil {
		return nil, err
	}

	port, err := helpers.GetFreePortTCP()
	if err != nil {
		return nil, err
	}

	return pool.BuildAndRunWithOptions(strings.TrimSuffix(path, "/test/tests")+"/Dockerfile", &dockertest.RunOptions{
		Name:     "dory_base_test",
		Tag:      "latest",
		Networks: []*dockertest.Network{network},
		Mounts: []string{
			path + "/configuration.json:/app/configuration.json",
		},
		ExtraHosts: []string{
			"host.docker.internal:host-gateway",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8000/tcp": {{ HostPort: fmt.Sprintf("%d", port) }},
		},
	})
}

func deleteContainer(pool *dockertest.Pool, container *dockertest.Resource) error {

	if err := pool.Purge(container); err != nil {
		return err
	}

	return nil
}

// ---- Structures

type email struct {
	Text       string    `json:"text"`
	TextAsHtml string    `json:"textAsHtml"`
	Subject    string    `json:"subject"`
	Date       time.Time `json:"date"`
}

type user struct {
	Username       string         `json:"username"`
	Password       string         `json:"password"`
	TOTP           string         `json:"totp"`
	NewPassword    string         `json:"new_password"`
	OldPassword    string         `json:"old_password"`
	Authentication authentication `json:"authentication"`
}

type totpStruct struct {
	TOTP string `json:"totp"`
}

type authentication struct {
	Token string `json:"token"`
	TOTP  string `json:"totp"`
}
