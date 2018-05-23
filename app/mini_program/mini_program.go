package mini_program

import (
	"github.com/godcong/wego"
	"github.com/godcong/wego/config"
	"github.com/godcong/wego/core"
	"github.com/godcong/wego/log"
)

type MiniProgram struct {
	config.Config
	client   *core.Client
	app      *core.Application
	token    *core.AccessToken
	auth     *Auth
	appCode  *AppCode
	dataCube *DataCube
	template *Template
}

var defaultConfig config.Config
var program *MiniProgram

func init() {
	defaultConfig = config.GetConfig("mini_program.default")
	app := core.App()
	program = newMiniProgram(defaultConfig, app)
	app.Register("mini_program", program)
	//app.Register(newMiniProgram())
}

func newMiniProgram(config config.Config, application *core.Application) *MiniProgram {
	mini0 := &MiniProgram{
		Config: config,
		app:    application,
		client: core.NewClient(config),
	}
	mini0.token = core.NewAccessToken(config, mini0.client)
	mini0.client.SetDataType(core.DATA_TYPE_JSON)
	mini0.client.SetDomain(core.NewDomain("mini_program"))
	return mini0
}

func (m *MiniProgram) SetClient(c *core.Client) *MiniProgram {
	m.client = c
	return m
}

func (m *MiniProgram) GetClient() *core.Client {
	return m.client
}
func (m *MiniProgram) Auth() wego.Auth {
	if m.auth == nil {
		m.auth = &Auth{
			Config:      m.Config,
			MiniProgram: m,
		}
	}
	return m.auth
}

func (m *MiniProgram) AppCode() wego.AppCode {
	if m.appCode == nil {
		m.appCode = &AppCode{
			Config:      m.Config,
			MiniProgram: m,
		}
	}
	return m.appCode
}

func (m *MiniProgram) DataCube() wego.DataCube {
	if m.dataCube == nil {
		m.dataCube = &DataCube{
			Config:      m.Config,
			MiniProgram: m,
		}
	}
	return m.dataCube
}

func (m *MiniProgram) Template() *Template {
	if m.template == nil {
		m.template = &Template{
			Config:      m.Config,
			MiniProgram: m,
		}
	}
	return m.template
}

func (m *MiniProgram) AccessToken() *core.AccessToken {
	log.Debug("MiniProgram|AccessToken")
	if m.token == nil {
		m.token = core.NewAccessToken(m.Config, m.client)
	}
	return m.token
}

//func (m *MiniProgram) accessToken() token.AccessTokenInterface {
//	if m.acc == nil {
//		m.acc = NewMiniProgramAccessToken(m.app, m.Config)
//	}
//	return m.acc
//}
//
//func (m *MiniProgram) Client() Client {
//	if m.client == nil {
//		m.client = app.Client(m.Config)
//	}
//	return m.client
//}
//
//func NewMiniProgram(application Application) MiniProgram {
//	config := application.GetConfig("mini_program.default")
//	return &MiniProgram{
//		Config: config,
//		app:    application,
//		client: app.Client(config),
//	}
//}
//
//type MiniProgramAccessToken struct {
//	token.accessToken
//	config.Config
//	app core.Application
//}
//
//func NewMiniProgramAccessToken(application Application, config Config) *MiniProgramAccessToken {
//	return &MiniProgramAccessToken{
//		Config: config,
//		app:    application,
//	}
//}
//
//func (m *MiniProgramAccessToken) getCredentials() util.Map {
//	return util.Map{
//		"grant_type": "client_credential",
//		"appid":      m.Get("app_id"),
//		"secret":     m.Get("secret"),
//	}
//}