package config

import (
	_ "embed"
	"os"

	"github.com/Jeffail/gabs/v2"
	coreSetting "github.com/red-gold/telar-core/config"
	log "github.com/red-gold/telar-core/pkg/log"
	actionSetting "github.com/red-gold/telar-web/micros/actions/config"
	authSetting "github.com/red-gold/telar-web/micros/auth/config"
	notifySetting "github.com/red-gold/telar-web/micros/notifications/config"
	storageSetting "github.com/red-gold/telar-web/micros/storage/config"
)

//go:embed config.development.json
var devConfig []byte

//go:embed config.production.json
var prodConfig []byte

//go:embed config.test.json
var testConfig []byte

var internalBaseURL string

var jsonParsed *gabs.Container

func init() {
	internalBaseURL = baseURL()
	env := getEnvironment()
	log.Info("App is running on %s environment.", env)
	var err error
	if env == "production" {
		jsonParsed, err = gabs.ParseJSON(prodConfig)
	} else if env == "test" {
		jsonParsed, err = gabs.ParseJSON(testConfig)
	} else {
		jsonParsed, err = gabs.ParseJSON(devConfig)
	}

	if err != nil {
		panic(err)
	}
}

// Initiailze core configurations
func InitCoreConfig(cfg *coreSetting.Configuration, microName string) {

	appName, ok := jsonParsed.Path("environment.app_name").Data().(string)
	if ok {
		cfg.AppName = &appName
		log.Info("[%s] app_name information loaded from env |%s|", microName, appName)
	}

	baseRoute, ok := jsonParsed.Path("environment.base_route_domain").Data().(string)
	if ok {
		cfg.BaseRoute = &baseRoute
		log.Info("[%s] base_route_domain information loaded from env |%s|", microName, baseRoute)
	}

	dbType, ok := jsonParsed.Path("environment.db_type").Data().(string)
	if ok {
		cfg.DBType = &dbType
		log.Info("[%s] db_type information loaded from env |%s|", microName, dbType)
	}

	headerCookieName, ok := jsonParsed.Path("environment.header_cookie_name").Data().(string)
	if ok {
		cfg.HeaderCookieName = &headerCookieName
		log.Info("[%s] header_cookie_name information loaded from env |%s|", microName, headerCookieName)
	}

	orgAvatar, ok := jsonParsed.Path("environment.org_avatar").Data().(string)
	if ok {
		cfg.OrgAvatar = &orgAvatar
		log.Info("[%s] org_avatar information loaded from env |%s|", microName, orgAvatar)
	}

	orgName, ok := jsonParsed.Path("environment.org_name").Data().(string)
	if ok {
		cfg.OrgName = &orgName
		log.Info("[%s] org_name information loaded from env |%s|", microName, orgName)
	}

	payloadCookieName, ok := jsonParsed.Path("environment.payload_cookie_name").Data().(string)
	if ok {
		cfg.PayloadCookieName = &payloadCookieName
		log.Info("[%s] payload_cookie_name information loaded from env |%s|", microName, payloadCookieName)
	}

	phoneSourceNumber, ok := jsonParsed.Path("environment.phone_source_number").Data().(string)
	if ok {
		cfg.PhoneSourceNumber = &phoneSourceNumber
		log.Info("[%s] phone_source_number information loaded from env |%s|", microName, phoneSourceNumber)
	}

	recaptchaSiteKey, ok := jsonParsed.Path("environment.recaptcha_site_key").Data().(string)
	if ok {
		cfg.RecaptchaSiteKey = &recaptchaSiteKey
		log.Info("[%s] recaptcha_site_key information loaded from env |%s|", microName, recaptchaSiteKey)
	}

	refEmail, ok := jsonParsed.Path("environment.ref_email").Data().(string)
	if ok {
		cfg.RefEmail = &refEmail
		log.Info("[%s] ref_email information loaded from env |%s|", microName, refEmail)
	}

	signatureCookieName, ok := jsonParsed.Path("environment.signature_cookie_name").Data().(string)
	if ok {
		cfg.SignatureCookieName = &signatureCookieName
		log.Info("[%s] signature_cookie_name information loaded from env |%s|", microName, signatureCookieName)
	}

	smtpEmail, ok := jsonParsed.Path("environment.smtp_email").Data().(string)
	if ok {
		cfg.SmtpEmail = &smtpEmail
		log.Info("[%s] smtp_email information loaded from env |%s|", microName, smtpEmail)
	}

	debug, ok := jsonParsed.Path("environment.debug").Data().(bool)
	if ok {
		cfg.Debug = &debug
		log.Info("[%s] debug information loaded from env |%t|", microName, debug)
	}

	gateway, ok := jsonParsed.Path("environment.gateway").Data().(string)
	if ok {
		cfg.Gateway = &gateway
		log.Info("[%s] gateway information loaded from env |%s|", microName, gateway)
	}

	origin, ok := jsonParsed.Path("environment.origin").Data().(string)
	if ok {
		cfg.Origin = &origin
		log.Info("[%s] origin information loaded from env |%s|", microName, origin)
	}

	cfg.InternalGateway = &internalBaseURL

}

// Initiailze auth micro configurations
func InitAuthConfig(cfg *authSetting.Configuration) {

	microName := "auth"

	baseRoute, ok := jsonParsed.Path("micros.auth.environment.base_route").Data().(string)
	if ok {
		cfg.BaseRoute = baseRoute
		log.Info("[%s] base_route information loaded from env |%s|", microName, baseRoute)
	}

	externalRedirectDomain, ok := jsonParsed.Path("micros.auth.environment.external_redirect_domain").Data().(string)
	if ok {
		cfg.ExternalRedirectDomain = externalRedirectDomain
		log.Info("[%s] external_redirect_domain information loaded from env |%s|", microName, externalRedirectDomain)
	}

	webURL, ok := jsonParsed.Path("micros.auth.environment.web_url").Data().(string)
	if ok {
		cfg.WebURL = webURL
		log.Info("[%s] web_url information loaded from env |%s|", microName, webURL)
	}

	authWebURI, ok := jsonParsed.Path("micros.auth.environment.auth_web_uri").Data().(string)
	if ok {
		cfg.AuthWebURI = authWebURI
		log.Info("[%s] auth_web_uri information loaded from env |%s|", microName, authWebURI)
	}

	clientID, ok := jsonParsed.Path("micros.auth.environment.client_id").Data().(string)
	if ok {
		cfg.ClientID = clientID
		log.Info("[%s] client_id information loaded from env |%s|", microName, clientID)
	}

	oAuthProvider, ok := jsonParsed.Path("micros.auth.environment.oauth_provider").Data().(string)
	if ok {
		cfg.OAuthProvider = oAuthProvider
		log.Info("[%s] oauth_provider information loaded from env |%s|", microName, oAuthProvider)
	}

	oAuthProviderBaseURL, ok := jsonParsed.Path("micros.auth.environment.oauth_provider_base_url").Data().(string)
	if ok {
		cfg.OAuthProviderBaseURL = oAuthProviderBaseURL
		log.Info("[%s] oauth_provider_base_url information loaded from env |%s|", microName, oAuthProviderBaseURL)
	}

	verifyType, ok := jsonParsed.Path("micros.auth.environment.verify_type").Data().(string)
	if ok {
		cfg.VerifyType = verifyType
		log.Info("[%s] verify_type information loaded from env |%s|", microName, verifyType)
	}

	cookieRootDomain, ok := jsonParsed.Path("environment.cookie_root_domain").Data().(string)
	if ok {
		cfg.CookieRootDomain = cookieRootDomain
		log.Info("[%s] cookie_root_domain information loaded from env |%s|", microName, cookieRootDomain)
	}

	cfg.QueryPrettyURL = true

}

// Initialize notificatin micro configurations
func InitNotifyConfig(cfg *notifySetting.Configuration) {
	microName := "notifications"

	webURL, ok := jsonParsed.Path("micros.notification.environment.web_url").Data().(string)
	if ok {
		cfg.WebURL = webURL
		log.Info("[%s] web_url information loaded from env |%s|", microName, webURL)
	}

	baseRoute, ok := jsonParsed.Path("micros.notification.environment.base_route").Data().(string)
	if ok {
		cfg.BaseRoute = baseRoute
		log.Info("[%s] base_route information loaded from env |%s|", microName, baseRoute)
	}

	cfg.QueryPrettyURL = true
}

// Initialize action micro configurations
func InitActionConfig(cfg *actionSetting.Configuration) {
	microName := "actions"

	baseRoute, ok := jsonParsed.Path("micros.actions.environment.base_route").Data().(string)
	if ok {
		cfg.BaseRoute = baseRoute
		log.Info("[%s] base_route information loaded from env |%s|", microName, baseRoute)
	}

	cfg.QueryPrettyURL = true

	websocketServerURL, ok := jsonParsed.Path("environment.websocket_server_url").Data().(string)
	if ok {
		cfg.WebsocketServerURL = websocketServerURL
		log.Info("[%s] websocket_server_url information loaded from env |%s|", microName, websocketServerURL)
	}

}

// Initialize storage micro configurations
func InitStorageConfig(cfg *storageSetting.Configuration) {
	microName := "storage"

	baseRoute, ok := jsonParsed.Path("micros.storage.environment.base_route").Data().(string)
	if ok {
		cfg.BaseRoute = baseRoute
		log.Info("[%s] base_route information loaded from env |%s|", microName, baseRoute)
	}

	BucketName, ok := jsonParsed.Path("micros.storage.environment.bucket_name").Data().(string)
	if ok {
		cfg.BucketName = BucketName
		log.Info("[%s] bucket_name information loaded from env |%s|", microName, BucketName)
	}

	cfg.QueryPrettyURL = true
}

// Get base URL
func baseURL() string {
	p := os.Getenv("PORT")
	if p == "" {
		p = "4000" // local dev
	}
	return "http://localhost:" + p
}

func getEnvironment() string {

	env := os.Getenv("TELAR_ENV")
	if env == "" {
		return "production"
	} else {
		return env
	}
}
