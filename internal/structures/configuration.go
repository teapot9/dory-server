package structures

type TLSMode string
const (
	TLSModeNone = "none"
	TLSModeSTARTTLS = "starttls"
	TLSModeTLS = "tls"
)

type Configuration struct {
	LDAPServer   LDAPServerConfig `json:"ldap_server"`
	Server       ServerConfig     `json:"server"`
	TOTP         TOTPConfig       `json:"totp"`
	MailServer   MailServerConfig `json:"mail_server"`
	FrontAddress string           `json:"front_address"`
	Features     Features         `json:"features"`
}

type LDAPServerConfig struct {
	Admin         LDAPServerAdmin `json:"admin"`
	BaseDN        string          `json:"base_dn"`
	FilterOn      string          `json:"filter_on"`
	Address       string          `json:"address"`
	Port          int             `json:"port"`
	Kind          string          `json:"kind"`
	SkipTLSVerify bool            `json:"skip_tls_verify"`
	EmailField    string          `json:"email_field"`
}

type LDAPServerAdmin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Port         int    `json:"port"`
	BasePath     string `json:"base_path"`
	DatabasePath string `json:"database_path"`
}

type TOTPConfig struct {
	Kind              string `json:"kind"`
	CustomServiceName string `json:"custom_service_name"`
	Secret            string `json:"secret"`
	OpenLDAPParamsDN  string `json:"openldap_params_dn"`
}

type MailServerConfig struct {
	Address       string  `json:"address"`
	Port          int     `json:"port"`
	Password      string  `json:"password"`
	SenderAddress string  `json:"sender_address"`
	SenderName    string  `json:"sender_name"`
	Subject       string  `json:"subject"`
	SkipTLSVerify bool    `json:"skip_tls_verify"`
	TLSMode       TLSMode `json:"tls_mode"`
}

type Features struct {
	DisableUnlock                   bool `json:"disable_unlock"`
	DisablePasswordUpdate           bool `json:"disable_password_update"`
	DisablePasswordReinitialization bool `json:"disable_password_reinitialization"`
	DisableTOTP                     bool `json:"disable_totp"`
}
