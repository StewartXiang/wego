package official

import (
	"github.com/godcong/wego/core"
	"github.com/godcong/wego/core/media"
	"github.com/godcong/wego/log"
	"github.com/godcong/wego/net"
	"github.com/godcong/wego/util"
)

/*Material Material */
type Material struct {
	*Media
}

func newMaterial(media *Media) *Material {
	return &Material{
		Media: media,
	}
}

/*NewMaterial NewMaterial */
func NewMaterial() *Material {
	return newMaterial(NewMedia())
}

//AddNews 新增永久素材
// http请求方式: POST，https协议
// https://api.weixin.qq.com/cgi-bin/material/add_news?access_token=ACCESS_TOKEN
func (m *Material) AddNews(articles []*media.Article) *net.Response {
	log.Debug("Material|AddNews", articles)
	key := m.token.GetToken().KeyMap()
	resp := m.client.HTTPPostJSON(
		m.client.Link(materialAddNewsURLSuffix),
		util.Map{"articles": articles},
		util.Map{net.RequestTypeQuery.String(): key})
	return resp
}

//AddMaterial 新增其他类型永久素材
// http请求方式: POST，需使用https
// https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=ACCESS_TOKEN&type=TYPE
// 调用示例（使用curl命令，用FORM表单方式新增一个其他类型的永久素材，curl命令的使用请自行查阅资料）
// 成功:
// {"media_id":"HIWcj9t3AI_b8qCQSu8lrY5DkGL1LMl8_eTrDv4aUo8","url":"http:\/\/mmbiz.qpic.cn\/mmbiz_jpg\/gJHMd2C74XpfUBCTPocUe1Dd8cXnAlDmRqdPoFWq1DvJZjdW5BCaYyu7NfHusicU50nRs8Vb1oiaNrwMbTtNcFtQ\/0?wx_fmt=jpeg"}
func (m *Material) AddMaterial(filePath string, mediaType core.MediaType) *net.Response {
	log.Debug("Material|AddMaterial", filePath, mediaType)
	if mediaType == core.MediaTypeVideo {
		log.Error("please use Material.UploadVideo() function")
	}

	p := m.token.GetToken().KeyMap()
	p.Set("type", mediaType.String())
	resp := m.client.HTTPUpload(
		m.client.Link(materialAddMaterialURLSuffix),
		util.Map{
			"media": filePath,
		},
		util.Map{
			net.RequestTypeQuery.String(): p,
		})
	return resp
}

//UploadVideo 新增其他类型永久素材
// http请求方式: POST，需使用https
// https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=ACCESS_TOKEN&type=TYPE
// 成功:
// {"media_id":"HIWcj9t3AI_b8qCQSu8lrTBEyIAO-uPSQhTiI2uoENk"}
func (m *Material) UploadVideo(filePath string, title, introduction string) *net.Response {
	log.Debug("Media|UploadVideo", filePath, title, introduction)
	p := m.token.GetToken().KeyMap()
	p.Set("type", core.MediaTypeVideo.String())
	resp := m.client.HTTPUpload(
		m.client.Link(materialAddMaterialURLSuffix),
		util.Map{
			"media": filePath,
			"description": util.Map{
				"title":        title,
				"introduction": introduction,
			},
		},
		util.Map{
			net.RequestTypeQuery.String(): p,
		})
	return resp
}

//Get 获取永久素材
// http请求方式: POST,https协议
// https://api.weixin.qq.com/cgi-bin/material/get_material?access_token=ACCESS_TOKEN
// 失败:
// {"errcode":40007,"errmsg":"invalid media_id hint: [YoxHSA07631538]"}
// 成功:
// {"title":"ceshi2","description":"only test","down_url":"http:\/\/203.205.158.71\/vweixinp.tc.qq.com\/1007_ad755ea12b3043e893e174d18de97f24.f10.mp4?vkey=22A7BCCDB429DF3613D50C1CAC510BDDCD12030895D782B3FAE00FB6989E4FFA640EB7EB8B498E560D08C84D808EF352BFFB0B15FA743556DB96BBF0239FC41F6DAFEEBA1024DBCA0954FBE09A66AA5381AB9CA50D1F8AE2&sha=0&save=1"}
func (m *Material) Get(mediaID string) *net.Response {
	log.Debug("Material|Get", mediaID)
	p := m.token.GetToken().KeyMap()
	resp := m.client.HTTPPostJSON(
		m.client.Link(materialGetMaterialURLSuffix),
		util.Map{
			"media_id": mediaID,
		},
		util.Map{
			net.RequestTypeQuery.String(): p,
		})
	return resp
}

//Del 删除永久素材
// http请求方式: POST
// https://api.weixin.qq.com/cgi-bin/material/del_material?access_token=ACCESS_TOKEN
// 成功:
// {"errcode":0,"errmsg":"ok"}
// 失败:
// {"errcode":-1,"errmsg":"system error hint: [NX0zcA05993060]"}
func (m *Material) Del(mediaID string) *net.Response {
	log.Debug("Material|Del", mediaID)
	p := m.token.GetToken().KeyMap()
	resp := m.client.HTTPPostJSON(
		m.client.Link(materialDelMaterialURLSuffix),
		util.Map{
			"media_id": mediaID,
		},
		util.Map{
			net.RequestTypeQuery.String(): p,
		})
	return resp

}

//UpdateNews 修改永久图文素材
// http请求方式: POST
// https://api.weixin.qq.com/cgi-bin/material/update_news?access_token=ACCESS_TOKEN
func (m *Material) UpdateNews(mediaID string, index int, articles []*media.Article) *net.Response {
	log.Debug("Material|UpdateNews", mediaID)
	p := m.token.GetToken().KeyMap()
	resp := m.client.HTTPPostJSON(
		m.client.Link(materialUpdateNewsURLSuffix),
		util.Map{
			"media_id": mediaID,
			"index":    index,
			"articles": articles,
		},
		util.Map{
			net.RequestTypeQuery.String(): p,
		})
	return resp

}

//GetMaterialCount 获取素材总数
// http请求方式: GET
// https://api.weixin.qq.com/cgi-bin/material/get_materialcount?access_token=ACCESS_TOKEN
// 成功:
// {"voice_count":0,"video_count":2,"image_count":0,"news_count":0}
// 失败:
// {"errcode":-1,"errmsg":"system error"}
func (m *Material) GetMaterialCount() *net.Response {
	log.Debug("Material|GetMaterialCount")
	p := m.token.GetToken().KeyMap()
	resp := m.client.HTTPGet(
		m.client.Link(materialGetMaterialcountURLSuffix),
		util.Map{
			net.RequestTypeQuery.String(): p,
		})
	return resp
}

//BatchGet BatchGet
// http请求方式: POST
// https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=ACCESS_TOKEN
// 成功:
// {"item":[{"media_id":"HIWcj9t3AI_b8qCQSu8lrTgTis9nPHNyIkIEWaDdHzY","name":"ceshi2","update_time":1521355652}],"total_count":2,"item_count":1}
// 失败:
// {"errcode":40007,"errmsg":"invalid media_id"}
func (m *Material) BatchGet(mediaType core.MediaType, offset, count int) *net.Response {
	log.Debug("Material|BatchGet", mediaType, offset, count)
	p := m.token.GetToken().KeyMap()
	resp := m.client.HTTPPostJSON(
		m.client.Link(materialBatchgetMaterialURLSuffix),
		util.Map{
			"type":   mediaType.String(),
			"offset": offset,
			"count":  count,
		},
		util.Map{
			net.RequestTypeQuery.String(): p,
		})
	return resp

}