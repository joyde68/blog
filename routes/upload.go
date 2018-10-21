package routes

import (
	"github.com/joyde68/blog/models"
	"gopkg.in/macaron.v1"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func AdminFiles(context *macaron.Context) {
	if context.Req.Method == "DELETE" {
		id := context.QueryInt("id")
		models.RemoveFile(id)
		models.Json(context, true).End()
		//context.Do("attach_delete", id)
		return
	}
	files, pager := models.GetFileList(context.ParamsInt("page"), 10)

	data := map[string]interface{}{
		"Title": "媒体文件",
		"Files": files,
		"Pager": pager,
	}


	err := models.Theme(true).Layout("layout").Tpl("files").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func FileUpload(context *macaron.Context) {
	req := context.Req.Request
	req.ParseMultipartForm(32 << 20)
	f, h, e := req.FormFile("file")
	if e != nil {
		models.Json(context, false).Set("msg", e.Error()).End()
		return
	}
	data, _ := ioutil.ReadAll(f)
	maxSize := 10485760
	defer func() {
		f.Close()
		data = nil
		h = nil
	}()
	if len(data) >= maxSize {
		models.Json(context, false).Set("msg", "文件应小于10M").End()
		return
	}
	uploadFileSuffix := ".jpg,.png,.gif,.zip,.txt,.doc,.docx,.xls,.xlsx,.ppt,.pptx"
	if !strings.Contains(uploadFileSuffix, path.Ext(h.Filename)) {
		models.Json(context, false).Set("msg", "文件只支持Office文件，图片和zip存档").End()
		return
	}
	ff := new(models.File)
	ff.Name = h.Filename
	ff.Type = context.Query("type")
	if ff.Type == ""{
		ff.Type = "image"
	}
	ff.Size = int64(len(data))
	ff.ContentType = h.Header["Content-Type"][0]
	ff.Author = context.GetCookieInt("token-user")
	ff.Url = models.CreateFilePath(path.Join("public", "upload"), ff)
	e = ioutil.WriteFile(ff.Url, data, os.ModePerm)
	if e != nil {
		models.Json(context, false).Set("msg", e.Error()).End()
		return
	}
	models.CreateFile(ff)
	models.Json(context, true).Set("file", ff).End()
	//context.Do("attach_created", ff)
}
