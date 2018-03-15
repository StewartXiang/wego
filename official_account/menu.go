package official_account

import (
	"github.com/godcong/wego/core"
	"github.com/godcong/wego/core/menu"
)

type Menu struct {
	config  core.Config
	account *OfficialAccount
	client  *core.Client
	token   *core.AccessToken
	buttons core.Map
	menuid  int
}

func newMenu(account *OfficialAccount) *Menu {
	return &Menu{
		config:  defaultConfig,
		account: account,
		client:  account.client,
		token:   account.token,
		buttons: make(core.Map),
	}
}

func NewMenu() *Menu {
	return newMenu(account)
}

func (m *Menu) SetButtons(b []*menu.Button) *Menu {
	m.buttons["button"] = b
	return m
}

func (m *Menu) GetButtons() []*menu.Button {
	if v, b := m.buttons["button"]; b {
		if v0, b := v.([]*menu.Button); b {
			return v0
		}
	}
	return nil
}

func (m *Menu) AddButton(b *menu.Button) *Menu {
	if v := m.GetButtons(); v != nil {
		m.buttons["button"] = append(v, b)
	} else {
		m.buttons["button"] = []*menu.Button{b}
	}
	return m
}

func (m *Menu) SetMatchRule(rule *menu.MatchRule) *Menu {
	m.buttons["matchrule"] = rule
	return m
}

func (m *Menu) SetMenuId(id int) *Menu {
	m.menuid = id
	return m
}

//个性化创建
//https://api.weixin.qq.com/cgi-bin/menu/addconditional?access_token=ACCESS_TOKEN
//成功:
//{"errcode":0,"errmsg":"ok"}

//自定义菜单
//https://api.weixin.qq.com/cgi-bin/menu/create?access_token=ACCESS_TOKEN
//成功:
// {"menuid":429680901}]
func (m *Menu) Create() *core.Response {
	token := m.token.GetToken()
	if _, b := m.buttons["matchrule"]; !b {
		resp := m.client.HttpPostJson(
			m.client.Link(MENU_CREATE_URL_SUFFIX),
			m.buttons,
			core.Map{core.REQUEST_TYPE_QUERY.String(): token.KeyMap()})
		return resp
	}
	resp := m.client.HttpPostJson(
		m.client.Link(MENU_ADDCONDITIONAL_URL_SUFFIX),
		m.buttons,
		core.Map{core.REQUEST_TYPE_QUERY.String(): token.KeyMap()})
	return resp
}

func (m *Menu) List() *core.Response {
	token := m.token.GetToken()
	resp := m.client.HttpGet(m.client.Link(MENU_GET_URL_SUFFIX), core.Map{
		core.REQUEST_TYPE_QUERY.String(): token.KeyMap(),
	})
	return resp

}

func (m *Menu) Current() *core.Response {
	token := m.token.GetToken()
	resp := m.client.HttpGet(m.client.Link(GET_CURRENT_SELFMENU_INFO_URL_SUFFIX), core.Map{
		core.REQUEST_TYPE_QUERY.String(): token.KeyMap(),
	})
	return resp
}

func (m *Menu) TryMatch(userId string) *core.Response {
	token := m.token.GetToken()
	resp := m.client.HttpPostJson(m.client.Link(MENU_TRYMATCH_URL_SUFFIX),
		core.Map{"user_id": userId},
		core.Map{core.REQUEST_TYPE_QUERY.String(): token.KeyMap()})
	return resp
}

func (m *Menu) Delete() *core.Response {
	token := m.token.GetToken()
	if m.menuid == 0 {
		resp := m.client.HttpGet(m.client.Link(MENU_DELETE_URL_SUFFIX), core.Map{
			core.REQUEST_TYPE_QUERY.String(): token.KeyMap(),
		})
		return resp
	}

	resp := m.client.HttpPostJson(m.client.Link(MENU_DELETECONDITIONAL_URL_SUFFIX),
		core.Map{"menuid": m.menuid},
		core.Map{core.REQUEST_TYPE_QUERY.String(): token.KeyMap()})
	return resp
}