package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	githubClient "github.com/google/go-github/v39/github"
	"github.com/rookie-ninja/rk-common/common"
	"github.com/rookie-ninja/rk-entry/entry"
	"github.com/rookie-ninja/rk-query"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"strings"
)

const (
	// EntryName name of entry
	EntryName = "ws-oauth"
	// EntryType type of entry
	EntryType = "ws-oauth"
	// EntryDescription description of entry
	EntryDescription = "Entry for oauth management entry."
	// GithubCallbackHost describes default callback address
	GithubCallbackHost = "http://localhost:8080"
	// Github type of oauth destination
	Github = "github"

	GithubAppPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAyAwt0oJ0zxzw8BoBaPx0iAoDoA94DpRtqHExsdbR5YbF7WFk
jGe14PTfrOWDU313ow7434+u9UEDtoXbCZlQ4PPbxHIKK59ztnct06gWjIxKi7Gm
XYShgYNBVc4GrANPz10TUqGCkHvaoHwdPkHyxGGGLHi7NClgKP103yRTilDiwt1q
b/oTrtTvFHXJOjWdy0vTbFFnnhEnzPWSx5L/g9fEfS4F5gGqjzzvxLfQZ75IDKis
77yNRPyvC9rWG3KNYSgQe4qi+YGZyzlq8XNp1YbruBxhALr4q33bkabC4CnhARRK
6lLFJ3Iz3+crXo2JR6mtmKj3+RYCgtuXxE4JdwIDAQABAoIBAQDFfJB38uXR2SZa
QbIGrMN10T0G9H53FjyzPxvqDsKjrssSr0UN/wxkihmOm/1rnL9Qr+Us/rGf2JEL
zVURQtO/X74jTtmPexhGl50cCCJKmey1v6kXRWoTazpCtOONSK+vbJk4x80XSc3p
BN3xqjU5fXVxzbYxpm8RyDVqXHhEYBqHhApa8MPSYv/qhCmYwggaRwgzFxcquB8X
JYD0QBG8CB0WYIIhcYAhns6qDqwcoH7qzHoF/ANIR9Jex/hHmIicMAbR2Y3QWhoy
hK+nE8OEdYa7kw/BHMGve6Spg028DSc50pjIry0Eo4ZD6ooO07KxImH5qCvg4JKl
0nIFVY3xAoGBAOpDBwPd1R62o2+DMjbwXeIYsCjpbqHuU9Rz00x7Ks4k2qDkJ9xr
xugp0wTw8t4PeXYTiEP9HgtPwCC0p8cqxvwSh116zcMPB71B7Vwg3P08neRIdWcg
bi75gAK+Ym8ziXiPASZX30XsElxwmM35yw6luxiZQ84X2n1CnbwLfYHNAoGBANqc
YUNBkA7LlKPMWOX86ELxAs2lx7TUwTQ8BF5NOdmv3AQNRRb1qSxAzyfw9Sf8kASV
gOgADnUnBm6EbMLDv7UG1QtYCkeiFPe3puoOf6uEp65U34WlqFCBCkdf/QD0fscU
4kP38cTfhqOH9kULsSpwmAB5K6QMF11fTnQ0LMRTAoGBAKjx7f073pdn4DZrx6sX
bp3AcEsRDlh6KLrvTVO7AAPrUED4SkcM80Y745OMsZq0TkR5kax2v1QpD8aGgvmA
QEFKm5UvG3WxQUOcaDIpATcgoD4ig4j8OnpmNYvFAfhwkpP/jjS46qzis9s22Pyz
SV4m5+e1oNDhIxFzGY6kOr+BAoGBAMyEgHnrXFp7GxQimQiRErmNwJGkBrGmWRoF
DBEtLnH5lFw2Dezs5tf/yc4UH0bJgfLH61EgvGXdnKbIPPf5KeCyA54ZP4TEndki
d4WBCu/rqvPtczAVSuIF1xfvNUMveWvGnef4jrcgZ1WWXU87IQQTUiEfOzS+Gx2/
jCYqqbcBAoGBAKgFgKEB1X7vpTR/ln96+PWgKQ8DJ9qFwkeKSuaDtROVDyK3Jrt1
Fl7nfF79WSy6aasnKiBIdJadNUhmeqokPhUJaBoBKRn7VouVtfD/t75hj3WffLH1
3gDn7GPBxLBM1AOxpGjmY2IvHGSCj3kFVdtBdWsv7yKPbbYzW3e8FBYU
-----END RSA PRIVATE KEY-----
`
)

// BootConfig is a struct which is for unmarshalled YAML
type BootConfig struct {
	Oauth struct {
		Enabled bool `yaml:"enabled" json:"enabled"`
		Github  struct {
			Enabled      bool     `yaml:"enabled" json:"enabled"`
			CallbackHost string   `yaml:"callbackHost" json:"callbackHost"`
			ClientId     string   `yaml:"clientId" json:"clientId"`
			ClientSecret string   `yaml:"clientSecret" json:"clientSecret"`
			Scopes       []string `yaml:"scopes" json:"scopes"`
		} `yaml:"github" json:"github"`
		Logger struct {
			ZapLogger struct {
				Ref string `yaml:"ref" json:"ref"`
			} `yaml:"zapLogger" json:"zapLogger"`
			EventLogger struct {
				Ref string `yaml:"ref" json:"ref"`
			} `yaml:"eventLogger" json:"eventLogger"`
		} `yaml:"logger" json:"logger"`
	} `yaml:"oauth" json:"oauth"`
}

// RegisterEntryFromConfig is an implementation of:
// type EntryRegFunc func(string) map[string]rkentry.Entry
func RegisterEntryFromConfig(configFilePath string) map[string]rkentry.Entry {
	res := make(map[string]rkentry.Entry)

	// 1: decode config map into boot config struct
	config := &BootConfig{}
	rkcommon.UnmarshalBootConfig(configFilePath, config)

	// 3: construct entry
	if config.Oauth.Enabled {
		opts := make([]EntryOption, 0)
		// github enabled
		if config.Oauth.Github.Enabled {
			githubConfig := &oauth2.Config{
				RedirectURL:  GithubCallbackHost + CallbackPathGithub,
				ClientID:     config.Oauth.Github.ClientId,
				ClientSecret: config.Oauth.Github.ClientSecret,
				Scopes:       config.Oauth.Github.Scopes,
				Endpoint:     github.Endpoint,
			}
			opts = append(opts, WithOauthConfig(Github, githubConfig))
		}

		entry := RegisterEntry(opts...)
		res[entry.GetName()] = entry
	}

	return res
}

// RegisterController will register Entry into GlobalAppCtx
func RegisterEntry(opts ...EntryOption) *Entry {
	entry := &Entry{
		EntryName:        EntryName,
		EntryType:        EntryType,
		EntryDescription: EntryDescription,
		ZapLoggerEntry:   rkentry.GlobalAppCtx.GetZapLoggerEntryDefault(),
		EventLoggerEntry: rkentry.GlobalAppCtx.GetEventLoggerEntryDefault(),
		oauthDest:        make(map[string]*oauth2.Config, 0),
	}

	for i := range opts {
		opts[i](entry)
	}

	fmt.Println(entry.oauthDest[Github].RedirectURL)

	rkentry.GlobalAppCtx.AddEntry(entry)

	return entry
}

// EntryOption will be extended in future.
type EntryOption func(*Entry)

// WithOauthConfig provide user
func WithOauthConfig(dest string, config *oauth2.Config) EntryOption {
	return func(entry *Entry) {
		entry.oauthDest[dest] = config
	}
}

// EntryImpl performs as manager of project and organizations
type Entry struct {
	EntryName        string                    `json:"entryName" yaml:"entryName"`
	EntryType        string                    `json:"entryType" yaml:"entryType"`
	EntryDescription string                    `json:"entryDescription" yaml:"entryDescription"`
	ZapLoggerEntry   *rkentry.ZapLoggerEntry   `json:"zapLoggerEntry" yaml:"zapLoggerEntry"`
	EventLoggerEntry *rkentry.EventLoggerEntry `json:"eventLoggerEntry" yaml:"eventLoggerEntry"`
	CallbackAddr     string                    `json:"callbackAddr" yaml:"callbackAddr"`
	oauthDest        map[string]*oauth2.Config `json:"-" yaml:"-"`
}

// Bootstrap entry
func (entry *Entry) Bootstrap(context.Context) {
	event := entry.EventLoggerEntry.GetEventHelper().Start(
		"bootstrap",
		rkquery.WithEntryName(entry.EntryName),
		rkquery.WithEntryType(entry.EntryType))

	logger := entry.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	initApi()

	entry.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Bootstrapping oauth entry.", event.ListPayloads()...)
}

// Interrupt entry
func (entry *Entry) Interrupt(context.Context) {
	event := entry.EventLoggerEntry.GetEventHelper().Start(
		"interrupt",
		rkquery.WithEntryName(entry.EntryName),
		rkquery.WithEntryType(entry.EntryType))
	logger := entry.ZapLoggerEntry.GetLogger().With(zap.String("eventId", event.GetEventId()))

	// TODO: Interrupting anything related.

	entry.EventLoggerEntry.GetEventHelper().Finish(event)
	logger.Info("Interrupting oauth entry.", event.ListPayloads()...)
}

// GetName returns entry name
func (entry *Entry) GetName() string {
	return entry.EntryName
}

// GetDescription returns entry description
func (entry *Entry) GetDescription() string {
	return entry.EntryDescription
}

// GetType returns entry type as project
func (entry *Entry) GetType() string {
	return entry.EntryType
}

// String returns entry as string
func (entry *Entry) String() string {
	bytes, _ := json.Marshal(entry)
	return string(bytes)
}

func (entry *Entry) IsValidOauthDest(src string) bool {
	_, ok := entry.oauthDest[strings.ToLower(src)]

	return ok
}

func (entry *Entry) GetOauthConfig(dest string) (*oauth2.Config, error) {
	dest = strings.ToLower(dest)

	if !entry.IsValidOauthDest(dest) {
		return nil, fmt.Errorf("unsupported oauth destination:%s", dest)
	}

	return entry.oauthDest[dest], nil
}

func (entry *Entry) GetGithubUser(accessToken string) (*githubClient.User, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := githubClient.NewClient(tc)

	user, _, err := client.Users.Get(context.Background(), "")
	return user, err
}

// GetEntry returns ProjectEntry.
func GetEntry() *Entry {
	if raw := rkentry.GlobalAppCtx.GetEntry(EntryName); raw != nil {
		return raw.(*Entry)
	}

	return nil
}
