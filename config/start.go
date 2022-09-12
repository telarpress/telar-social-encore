package config

import (
	"bytes"
	_ "embed"
	b64 "encoding/base64"
	"io/ioutil"
	"os"

	"github.com/Jeffail/gabs/v2"
	"github.com/joho/godotenv"
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
var envParsed map[string]string

type AllSecrets struct {
	AdminUsername  string
	AdminPassword  string
	MongoHost      string
	MongoDatabase  string
	PhoneAuthId    string
	PhoneAuthToken string
	Key            string
	KeyPub         string
	RefEmailPass   string
	PayloadSecret  string
	ServiceAccount string
	TSClientSecret string
	RecaptchaKey   string
}

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
		var envFile []byte
		var secretSource = jsonParsed.Path("environment.secret_source").Data().(string)
		if secretSource == "env" {
			envFile, err = ioutil.ReadFile("config/.env")
			reader := bytes.NewReader(envFile)
			if err != nil {
				panic(err)
			}
			envParsed, err = godotenv.Parse(reader)
		}
	}

	if err != nil {
		panic(err)
	}
}

// Initiailze core configurations
func InitCoreConfig(microName string, cfg *coreSetting.Configuration, secrets *AllSecrets) {

	appName, ok := jsonParsed.Path("environment.app_name").Data().(string)
	if ok {
		cfg.AppName = &appName
		log.Info("[%s] app_name information loaded from config file |%s|", microName, appName)
	}

	baseRoute, ok := jsonParsed.Path("environment.base_route_domain").Data().(string)
	if ok {
		cfg.BaseRoute = &baseRoute
		log.Info("[%s] base_route_domain information loaded from config file |%s|", microName, baseRoute)
	}

	dbType, ok := jsonParsed.Path("environment.db_type").Data().(string)
	if ok {
		cfg.DBType = &dbType
		log.Info("[%s] db_type information loaded from config file |%s|", microName, dbType)
	}

	headerCookieName, ok := jsonParsed.Path("environment.header_cookie_name").Data().(string)
	if ok {
		cfg.HeaderCookieName = &headerCookieName
		log.Info("[%s] header_cookie_name information loaded from config file |%s|", microName, headerCookieName)
	}

	orgAvatar, ok := jsonParsed.Path("environment.org_avatar").Data().(string)
	if ok {
		cfg.OrgAvatar = &orgAvatar
		log.Info("[%s] org_avatar information loaded from config file |%s|", microName, orgAvatar)
	}

	orgName, ok := jsonParsed.Path("environment.org_name").Data().(string)
	if ok {
		cfg.OrgName = &orgName
		log.Info("[%s] org_name information loaded from config file |%s|", microName, orgName)
	}

	payloadCookieName, ok := jsonParsed.Path("environment.payload_cookie_name").Data().(string)
	if ok {
		cfg.PayloadCookieName = &payloadCookieName
		log.Info("[%s] payload_cookie_name information loaded from config file |%s|", microName, payloadCookieName)
	}

	phoneSourceNumber, ok := jsonParsed.Path("environment.phone_source_number").Data().(string)
	if ok {
		cfg.PhoneSourceNumber = &phoneSourceNumber
		log.Info("[%s] phone_source_number information loaded from config file |%s|", microName, phoneSourceNumber)
	}

	recaptchaSiteKey, ok := jsonParsed.Path("environment.recaptcha_site_key").Data().(string)
	if ok {
		cfg.RecaptchaSiteKey = &recaptchaSiteKey
		log.Info("[%s] recaptcha_site_key information loaded from config file |%s|", microName, recaptchaSiteKey)
	}

	refEmail, ok := jsonParsed.Path("environment.ref_email").Data().(string)
	if ok {
		cfg.RefEmail = &refEmail
		log.Info("[%s] ref_email information loaded from config file |%s|", microName, refEmail)
	}

	signatureCookieName, ok := jsonParsed.Path("environment.signature_cookie_name").Data().(string)
	if ok {
		cfg.SignatureCookieName = &signatureCookieName
		log.Info("[%s] signature_cookie_name information loaded from config file |%s|", microName, signatureCookieName)
	}

	smtpEmail, ok := jsonParsed.Path("environment.smtp_email").Data().(string)
	if ok {
		cfg.SmtpEmail = &smtpEmail
		log.Info("[%s] smtp_email information loaded from config file |%s|", microName, smtpEmail)
	}

	debug, ok := jsonParsed.Path("environment.debug").Data().(bool)
	if ok {
		cfg.Debug = &debug
		log.Info("[%s] debug information loaded from config file |%t|", microName, debug)
	}

	gateway, ok := jsonParsed.Path("environment.gateway").Data().(string)
	if ok {
		cfg.Gateway = &gateway
		log.Info("[%s] gateway information loaded from config file |%s|", microName, gateway)
	}

	origin, ok := jsonParsed.Path("environment.origin").Data().(string)
	if ok {
		cfg.Origin = &origin
		log.Info("[%s] origin information loaded from config file |%s|", microName, origin)
	}

	cfg.InternalGateway = &internalBaseURL

	var secretSource = jsonParsed.Path("environment.secret_source").Data().(string)
	if secretSource == "env" {
		// Set secrets from env file
		payloadSecret, ok := envParsed["PayloadSecret"]
		if ok {
			cfg.PayloadSecret = &payloadSecret
			log.Info("[%s] payloadSecret information loaded from env file", microName)
		}

		publicKey, ok := envParsed["KeyPub"]
		if ok {
			decodedInBytes, _ := b64.StdEncoding.DecodeString(publicKey)
			decodedInString := string(decodedInBytes)
			cfg.PublicKey = &decodedInString
			log.Info("[%s] publicKey information loaded from env file", microName)
		}

		privateKey, ok := envParsed["Key"]
		if ok {
			decodedInBytes, _ := b64.StdEncoding.DecodeString(privateKey)
			decodedInString := string(decodedInBytes)
			cfg.PrivateKey = &decodedInString
			log.Info("[%s] privateKey information loaded from env file", microName)
		}

		refEmailPass, ok := envParsed["RefEmailPass"]
		if ok {
			cfg.RefEmailPass = &refEmailPass
			log.Info("[%s] refEmailPass information loaded from env file", microName)
		}

		recaptchaKey, ok := envParsed["RecaptchaKey"]
		if ok {
			cfg.RecaptchaKey = &recaptchaKey
			log.Info("[%s] recaptchaKey information loaded from env file", microName)
		}

		mongoHost, ok := envParsed["MongoHost"]
		if ok {
			cfg.MongoDBHost = &mongoHost
			log.Info("[%s] mongoHost information loaded from env file", microName)
		}

		mongoDatabase, ok := envParsed["MongoDatabase"]
		if ok {
			cfg.Database = &mongoDatabase
			log.Info("[%s] mongoDatabase information loaded from env file", microName)
		}

	} else {
		// Set secrets from encore
		cfg.PayloadSecret = &secrets.PayloadSecret
		cfg.PublicKey = &secrets.KeyPub
		cfg.PrivateKey = &secrets.Key
		cfg.RefEmailPass = &secrets.RefEmailPass
		cfg.RecaptchaKey = &secrets.RecaptchaKey
		cfg.MongoDBHost = &secrets.MongoHost
		cfg.Database = &secrets.MongoDatabase
	}

}

// Initiailze auth micro configurations
func InitAuthConfig(cfg *authSetting.Configuration, secrets *AllSecrets) {

	microName := "auth"

	baseRoute, ok := jsonParsed.Path("micros.auth.environment.base_route").Data().(string)
	if ok {
		cfg.BaseRoute = baseRoute
		log.Info("[%s] base_route information loaded from config file |%s|", microName, baseRoute)
	}

	externalRedirectDomain, ok := jsonParsed.Path("micros.auth.environment.external_redirect_domain").Data().(string)
	if ok {
		cfg.ExternalRedirectDomain = externalRedirectDomain
		log.Info("[%s] external_redirect_domain information loaded from config file |%s|", microName, externalRedirectDomain)
	}

	webURL, ok := jsonParsed.Path("micros.auth.environment.web_url").Data().(string)
	if ok {
		cfg.WebURL = webURL
		log.Info("[%s] web_url information loaded from config file |%s|", microName, webURL)
	}

	authWebURI, ok := jsonParsed.Path("micros.auth.environment.auth_web_uri").Data().(string)
	if ok {
		cfg.AuthWebURI = authWebURI
		log.Info("[%s] auth_web_uri information loaded from config file |%s|", microName, authWebURI)
	}

	clientID, ok := jsonParsed.Path("micros.auth.environment.client_id").Data().(string)
	if ok {
		cfg.ClientID = clientID
		log.Info("[%s] client_id information loaded from config file |%s|", microName, clientID)
	}

	oAuthProvider, ok := jsonParsed.Path("micros.auth.environment.oauth_provider").Data().(string)
	if ok {
		cfg.OAuthProvider = oAuthProvider
		log.Info("[%s] oauth_provider information loaded from config file |%s|", microName, oAuthProvider)
	}

	oAuthProviderBaseURL, ok := jsonParsed.Path("micros.auth.environment.oauth_provider_base_url").Data().(string)
	if ok {
		cfg.OAuthProviderBaseURL = oAuthProviderBaseURL
		log.Info("[%s] oauth_provider_base_url information loaded from config file |%s|", microName, oAuthProviderBaseURL)
	}

	verifyType, ok := jsonParsed.Path("micros.auth.environment.verify_type").Data().(string)
	if ok {
		cfg.VerifyType = verifyType
		log.Info("[%s] verify_type information loaded from config file |%s|", microName, verifyType)
	}

	cookieRootDomain, ok := jsonParsed.Path("environment.cookie_root_domain").Data().(string)
	if ok {
		cfg.CookieRootDomain = cookieRootDomain
		log.Info("[%s] cookie_root_domain information loaded from config file |%s|", microName, cookieRootDomain)
	}

	cfg.QueryPrettyURL = true

	var secretSource = jsonParsed.Path("environment.secret_source").Data().(string)
	if secretSource == "env" {
		// Set secrets from env file
		cfg.OAuthClientSecret = envParsed["TSClientSecret"]
		cfg.AdminUsername = envParsed["AdminUsername"]
		cfg.AdminPassword = envParsed["AdminPassword"]
	} else {
		// Set secrets from encore
		cfg.OAuthClientSecret = secrets.TSClientSecret
		cfg.AdminUsername = secrets.AdminUsername
		cfg.AdminPassword = secrets.AdminPassword
	}

}

// Initialize notificatin micro configurations
func InitNotifyConfig(cfg *notifySetting.Configuration) {
	microName := "notifications"

	webURL, ok := jsonParsed.Path("micros.notification.environment.web_url").Data().(string)
	if ok {
		cfg.WebURL = webURL
		log.Info("[%s] web_url information loaded from config file |%s|", microName, webURL)
	}

	baseRoute, ok := jsonParsed.Path("micros.notification.environment.base_route").Data().(string)
	if ok {
		cfg.BaseRoute = baseRoute
		log.Info("[%s] base_route information loaded from config file |%s|", microName, baseRoute)
	}

	cfg.QueryPrettyURL = true
}

// Initialize action micro configurations
func InitActionConfig(cfg *actionSetting.Configuration) {
	microName := "actions"

	baseRoute, ok := jsonParsed.Path("micros.actions.environment.base_route").Data().(string)
	if ok {
		cfg.BaseRoute = baseRoute
		log.Info("[%s] base_route information loaded from config file |%s|", microName, baseRoute)
	}

	cfg.QueryPrettyURL = true

	websocketServerURL, ok := jsonParsed.Path("environment.websocket_server_url").Data().(string)
	if ok {
		cfg.WebsocketServerURL = websocketServerURL
		log.Info("[%s] websocket_server_url information loaded from config file |%s|", microName, websocketServerURL)
	}

}

// Initialize storage micro configurations
func InitStorageConfig(cfg *storageSetting.Configuration, secrets *AllSecrets) {
	microName := "storage"

	baseRoute, ok := jsonParsed.Path("micros.storage.environment.base_route").Data().(string)
	if ok {
		cfg.BaseRoute = baseRoute
		log.Info("[%s] base_route information loaded from config file |%s|", microName, baseRoute)
	}

	BucketName, ok := jsonParsed.Path("micros.storage.environment.bucket_name").Data().(string)
	if ok {
		cfg.BucketName = BucketName
		log.Info("[%s] bucket_name information loaded from config file |%s|", microName, BucketName)
	}

	cfg.QueryPrettyURL = true

	var secretSource = jsonParsed.Path("environment.secret_source").Data().(string)
	if secretSource == "env" {
		// Set secrets from env file
		decodedInBytes, _ := b64.StdEncoding.DecodeString(envParsed["ServiceAccount"])
		decodedInString := string(decodedInBytes)
		cfg.StorageSecret = decodedInString
	} else {
		cfg.StorageSecret = secrets.ServiceAccount
	}
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
