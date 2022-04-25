package config

import (
	_ "embed"
	"fmt"
	"os"

	coreSetting "github.com/red-gold/telar-core/config"
	authSetting "github.com/red-gold/telar-web/micros/auth/config"
	"gopkg.in/yaml.v2"
)

//go:embed app_config.yml
var appConfigYaml []byte
//go:embed auth_config.yml
var authConfigYaml []byte
//go:embed gateway_config.yml
var gatewayConfigYaml []byte

// Is development mode
var isDev = false

type AppConfig struct {
	Environment struct {
		AppName string `yaml:"app_name"`
		BaseRouteDomain string `yaml:"base_route_domain"`
		DBType	string `yaml:"db_type"`
		HeaderCookieName string `yaml:"header_cookie_name"`
		OrgAvatar string `yaml:"org_avatar"`
		OrgName string `yaml:"org_name"`
		PayloadCookieName string `yaml:"payload_cookie_name"`
		PhoneSourceNumber string `yaml:"phone_source_number"`
		ReadTimeout int `yaml:"read_timeout"`
		WriteTimeout int `yaml:"write_timeout"`
		RecaptchaSiteKey string `yaml:"recaptcha_site_key"`
		RedisAddress string `yaml:"redis_address"`
		RefEmail string `yaml:"ref_email"`
		SignatureCookieName string `yaml:"signature_cookie_name"`
		SmtpEmail string `yaml:"smtp_email"`
		Debug bool `yaml:"debug"`

	} `yaml:"environment"`
}

type AuthConfig struct {
	Environment struct {
		BaseRoute string `yaml:"base_route"`
		ExternalRedirectDomain string `yaml:"external_redirect_domain"`
		WebURL string `yaml:"web_url"`
		AuthWebURI string `yaml:"auth_web_uri"`
		ClientID string `yaml:"client_id"`
		GithubAppID string `yaml:"github_app_id"`
		OAuthProvider string `yaml:"oauth_provider"`
		OAuthProviderBaseURL string `yaml:"oauth_provider_base_url"`
		ReportStatus string `yaml:"report_status"`
		VerifyType string `yaml:"verify_type"`
		WriteDebug bool `yaml:"write_debug"`
		ExecTimeout int `yaml:"exec_timeout"`
		ReadTimeout int `yaml:"read_timeout"`
		WriteTimeout int `yaml:"write_timeout"`
	} `yaml:"environment"`
}

type GatewayConfig struct {
	Environment struct {
		CookieRootDomain string `yaml:"cookie_root_domain"`
		Gateway string `yaml:"gateway"`
		InternalGateway string `yaml:"internal_gateway"`
		Origin string `yaml:"origin"`
		WebSocketServerURL string `yaml:"web_socket_server_url"`
	} `yaml:"environment"`
}

	
// Initiailze core configurations
func InitCoreConfig(cfg *coreSetting.Configuration) {
	fmt.Println("[ "+os.Getenv("PORT")+" ]")
	
	// Parse app config
	var appConfig AppConfig
	yaml.Unmarshal(appConfigYaml, &appConfig)
	cfg.AppName = &appConfig.Environment.AppName
	cfg.BaseRoute = &appConfig.Environment.BaseRouteDomain
	cfg.DBType = &appConfig.Environment.DBType
	cfg.HeaderCookieName = &appConfig.Environment.HeaderCookieName
	cfg.OrgAvatar = &appConfig.Environment.OrgAvatar
	cfg.OrgName = &appConfig.Environment.OrgName
	cfg.PayloadCookieName = &appConfig.Environment.PayloadCookieName
	cfg.PhoneSourceNumber = &appConfig.Environment.PhoneSourceNumber
	cfg.RecaptchaSiteKey = &appConfig.Environment.RecaptchaSiteKey
	cfg.RefEmail = &appConfig.Environment.RefEmail
	cfg.SignatureCookieName = &appConfig.Environment.SignatureCookieName
	cfg.SmtpEmail = &appConfig.Environment.SmtpEmail
	cfg.Debug = &appConfig.Environment.Debug
	
	// Parse gateway config
	var gatewayConfig GatewayConfig
	yaml.Unmarshal(gatewayConfigYaml, &gatewayConfig)
	cfg.Gateway = &gatewayConfig.Environment.Gateway
	cfg.InternalGateway = &gatewayConfig.Environment.InternalGateway
	cfg.Origin = &gatewayConfig.Environment.Origin

	
}

// Initiailze auth micro configurations
func InitAuthConfig(cfg *authSetting.Configuration) {

	var authConfig AuthConfig

	// Parse auth config
	yaml.Unmarshal(authConfigYaml, &authConfig)
	cfg.BaseRoute = authConfig.Environment.BaseRoute
	cfg.ExternalRedirectDomain = authConfig.Environment.ExternalRedirectDomain
	cfg.WebURL = authConfig.Environment.WebURL
	cfg.AuthWebURI = authConfig.Environment.AuthWebURI
	cfg.ClientID = authConfig.Environment.ClientID
	cfg.OAuthProvider = authConfig.Environment.OAuthProvider
	cfg.OAuthProviderBaseURL = authConfig.Environment.OAuthProviderBaseURL
	cfg.VerifyType = authConfig.Environment.VerifyType
	cfg.QueryPrettyURL = true
	
	// Parse gateway config
	var gatewayConfig GatewayConfig
	yaml.Unmarshal(gatewayConfigYaml, &gatewayConfig)
	cfg.CookieRootDomain = gatewayConfig.Environment.CookieRootDomain
	

	
}
