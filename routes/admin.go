package routes

import (
	"fmt"
	"github.com/joyde68/blog/models"
	"github.com/joyde68/blog/pkg"
	"gopkg.in/macaron.v1"
	"strconv"
	"strings"
)

func Admin(context *macaron.Context) {
	uid := context.GetCookieInt("token-user")
	user := models.GetUserById(uid)

	data := map[string]interface{}{
		"Title":    "控制台",
		"Statis":   models.NewStatis(),
		"User":     user,
		"Messages": models.GetUnreadMessages(),
	}

	err := models.Theme(true).Layout("layout").Tpl("home").Render(context, 200, data)
	if err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func AdminProfile(context *macaron.Context) {
	uid := context.GetCookieInt("token-user")
	user := models.GetUserById(uid)
	if context.Req.Method == "POST" {
		if !user.ChangeEmail(context.Query("email")) {
			models.Json(context, false).Set("msg", "邮箱与别的用户重复").End()
			return
		}
		user.Name = context.Query("user")
		user.Email = context.Query("email")
		user.Avatar = pkg.Gravatar(user.Email, "180")
		user.Url = context.Query("url")
		user.Nick = context.Query("nick")
		user.Bio = context.Query("bio")
		//Json(context, true).End()
		models.Json(context, true).End()
		go models.SyncUsers()
		go models.UpdateCommentAdmin(user)
		// 更新最近登录时间
		//context.Do("profile_update", user)
		return
	}

	data := map[string]interface{}{
		"Title": "个性资料",
		"User":  user,
	}


	err := models.Theme(true).Layout("layout").Tpl("profile").Render(context, 200, data)
	if err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func AdminPassword(context *macaron.Context) {
	if context.Req.Method == "POST" {
		uid := context.GetCookieInt("token-user")
		user := models.GetUserById(uid)
		if !user.CheckPassword(context.Query("old")) {
			models.Json(context, false).Set("msg", "旧密码错误").End()
			return
		}
		user.ChangePassword(context.Query("new"))
		go models.SyncUsers()
		models.Json(context, true).End()
		// 待完善
		//context.Do("password_update", user)
		return
	}

	data := map[string]interface{}{
		"Title": "修改密码",
		//"User":user,
	}

	err := models.Theme(true).Layout("layout").Tpl("password").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func AdminArticle(context *macaron.Context) {
	articles, pager := models.GetArticleList(context.ParamsInt("page"), 10)

	data := map[string]interface{}{
		"Title":    "文章",
		"Articles": articles,
		"Pager":    pager,
	}


	err := models.Theme(true).Layout("layout").Tpl("articles").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func ArticleWrite(context *macaron.Context) {
	if context.Req.Method == "POST" {
		c := new(models.Content)
		c.Id = 0
		if !c.ChangeSlug(context.Query("slug")) {
			//Json(context, false).Set("msg", "固定链接重复").End()
			models.Json(context, false).End()
			return
		}
		c.Title = context.Query("title")
		c.Text = context.Query("content")
		c.Tags = strings.Split(strings.Replace(context.Query("tags"), "，", ",", -1), ",")
		c.IsComment = context.QueryBool("comment")
		c.IsLinked = false
		c.AuthorId = context.GetCookieInt("token-user")
		c.Template = "blog.html"
		c.Status = context.Query("status")
		c.Format = "markdown"
		c.Hits = 1
		var e error
		c, e = models.CreateContent(c, "article")
		if e != nil {
			//Json(context, false).Set("msg", e.Error()).End()
			models.Json(context, false).Set("msg", e.Error()).End()
			return
		}
		//Json(context, true).Set("content", c).End()
		models.Json(context, true).Set("content", c).End()
		// 待完善
		//context.Do("article_created", c)
		//c.Type = "article"
		return
	}

	data := map[string]interface{}{
		"Title": "撰写文章",
	}

	err := models.Theme(true).Layout("layout").Tpl("write_article").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func ArticleEdit(context *macaron.Context) {
	id := context.ParamsInt("id")
	c := models.GetContentById(id)
	if c == nil {
		models.Theme(false).Tpl("404").Render(context, 404, nil)
		return
	}
	if context.Req.Method == "DELETE" {
		models.RemoveContent(c)
		//Json(context, true).End()
		models.Json(context, true).End()
		return
	}
	if context.Req.Method == "POST" {
		if !c.ChangeSlug(context.Query("slug")) {
			models.Json(context, false).Set("msg", "固定链接重复").End()
			return
		}
		c.Title = context.Query("title")
		c.Text = context.Query("content")
		c.Tags = strings.Split(strings.Replace(context.Query("tags"), "，", ",", -1), ",")
		c.IsComment = context.QueryBool("comment")
		//c.IsLinked = false
		//c.AuthorId, _ = strconv.Atoi(context.Cookie("token-user"))
		//c.Template = "blog.html"
		c.Status = context.Query("status")
		//c.Format = "markdown"
		models.SaveContent(c)
		//Json(context, true).Set("content", c).End()
		models.Json(context, true).Set("content", c).End()
		//待完善
		//context.Do("article_modified", c)
		//c.Type = "article"
		return
	}

	data := map[string]interface{}{
		"Title":   "编辑文章",
		"Article": c,
	}


	err := models.Theme(true).Layout("layout").Tpl("edit_article").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func AdminPage(context *macaron.Context) {
	pages, pager := models.GetPageList(context.ParamsInt("page"), 10)
	data := map[string]interface{}{
		"Title": "页面",
		"Pages": pages,
		"Pager": pager,
	}

	err := models.Theme(true).Layout("layout").Tpl("pages").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func PageWrite(context *macaron.Context) {
	if context.Req.Method == "POST" {
		c := new(models.Content)
		c.Id = 0
		if !c.ChangeSlug(context.Query("slug")) {
			models.Json(context, false).Set("msg", "固定链接重复").End()
			return
		}
		c.Title = context.Query("title")
		c.Text = context.Query("content")
		c.Tags = make([]string, 0)
		c.IsComment = context.QueryBool("comment")
		c.IsLinked = context.QueryBool("link")
		c.AuthorId = context.GetCookieInt("token-user")
		c.Template = "page.html"
		c.Status = context.Query("status")
		c.Format = "markdown"
		c.Hits = 1
		var e error
		c, e = models.CreateContent(c, "page")
		if e != nil {
			//Json(context, false).Set("msg", e.Error()).End()
			models.Json(context, false).Set("msg", e.Error()).End()
			return
		}
		models.Json(context, true).Set("content", c).End()
		//c.Type = "article"
		// 待完善
		//context.Do("page_created", c)
		return
	}

	data := map[string]interface{}{
		"Title": "撰写页面",
	}


	err := models.Theme(true).Layout("layout").Tpl("write_page").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func PageEdit(context *macaron.Context) {
	id := context.ParamsInt("id")
	c := models.GetContentById(id)
	if c == nil {
		context.Redirect("/admin/pages/")
		return
	}
	if context.Req.Method == "DELETE" {
		models.RemoveContent(c)
		models.Json(context, true).End()
		return
	}
	if context.Req.Method == "POST" {
		if !c.ChangeSlug(context.Query("slug")) {
			models.Json(context, false).Set("msg", "固定链接重复").End()
			return
		}
		c.Title = context.Query("title")
		c.Text = context.Query("content")
		//c.Tags = strings.Split(strings.Replace(data["tag"], "，", ",", -1), ",")
		c.IsComment = context.QueryBool("comment")
		c.IsLinked = context.QueryBool("link")
		//c.AuthorId, _ = strconv.Atoi(context.Cookie("token-user"))
		//c.Template = "blog.html"
		c.Status = context.Query("status")
		//c.Format = "markdown"
		models.SaveContent(c)
		models.Json(context, true).Set("content", c).End()
		// 待修改
		//context.Do("page_modified", c)
		//c.Type = "article"
		return
	}

	data := map[string]interface{}{
		"Title": "编辑文章",
		"Page":  c,
	}

	err := models.Theme(true).Layout("layout").Tpl("edit_page").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func AdminSetting(context *macaron.Context) {
	if context.Req.Method == "POST" {
		context.Req.ParseForm()
		for k, v := range context.Req.Form {
			if v[0] == "" {
				if context.Req.Form.Get(k+"_def") != "" {
					v[0] = context.Req.Form.Get(k+"_def")
				}
			}
			models.SetSetting(k, v[0])
		}
		models.SyncSettings()
		//Json(context, true).End()
		models.Json(context, true).End()
		// 待完善
		//context.Do("setting_saved")
		return
	}

	data := map[string]interface{}{
		"Title":      "配置",
		"Custom":     models.GetCustomSettings(),
		"Navigators": models.GetNavigators(),
	}

	err := models.Theme(true).Layout("layout").Tpl("setting").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func CustomSetting(context *macaron.Context) {
	keys := context.QueryStrings("key")
	values := context.QueryStrings("value")

	for i, k := range keys {
		if len(k) < 1 {
			continue
		}
		models.SetSetting("c_"+k, values[i])
	}
	models.SyncSettings()
	models.Json(context, true).End()
	// 待完善
	//context.Do("setting_saved")
}

func NavigatorSetting(context *macaron.Context) {
	order := context.QueryStrings("order")
	text := context.QueryStrings("text")
	title := context.QueryStrings("title")
	link := context.QueryStrings("link")
	models.SetNavigators(order, text, title, link)
	models.Json(context, true).End()
	//context.Do("setting_saved")
}

func AdminComments(context *macaron.Context) {
	if context.Req.Method == "DELETE" {
		id := context.QueryInt("id")
		cmt := models.GetCommentById(id)
		models.RemoveComment(cmt.Cid, id)
		//Json(context, true).End()
		models.Json(context, true).End()
		//context.Do("comment_delete", id)
		return
	}
	if context.Req.Method == "PUT" {
		id := context.QueryInt("id")
		cmt2 := models.GetCommentById(id)
		cmt2.Status = "approved"
		cmt2.GetReader().Active = true
		models.SaveComment(cmt2)
		models.Json(context, true).End()
		//context.Do("comment_change_status", cmt2)
		return
	}
	if context.Req.Method == "POST" {
		// get required data
		pid := context.QueryInt("pid")
		cid := models.GetCommentById(pid).Cid
		uid, _ := strconv.Atoi(context.GetCookie("token-user"))
		user := models.GetUserById(uid)

		co := new(models.Comment)
		co.Author = user.Nick
		co.Email = user.Email
		co.Url = user.Url
		co.Content = context.Query("content")
		co.Avatar = pkg.Gravatar(co.Email, "50")
		co.Pid = pid
		co.Ip = context.RemoteAddr()
		co.UserAgent = context.Req.Header.Get("User-Agent")
		co.IsAdmin = true
		models.CreateComment(cid, co)
		models.Json(context, true).Set("comment", co.ToJson()).End()
		models.CreateMessage("comment", co)
		go models.SendEmail(co)
		return
	}
	page := context.QueryInt("page")
	if page == 0 {
		page = 1
	}
	comments, pager := models.GetCommentList(page, 10)

	data := map[string]interface{}{
		"Title":    "评论",
		"Comments": comments,
		"Pager":    pager,
	}

	err := models.Theme(true).Layout("layout").Tpl("comments").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func AdminMessageRead(context *macaron.Context) {
	id := context.QueryInt("id")
	if id < 0 {
		//Json(context, false).End()
		models.Json(context, false).End()
		return
	}
	m := models.GetMessage(id)
	if m == nil {
		//Json(context, false).End()
		models.Json(context, false).End()
		return
	}
	models.SaveMessageRead(m)
	models.Json(context, true).End()
}

func AdminReader(context *macaron.Context) {
	if context.Req.Method == "POST" {
		email :=context.Query("email")
		models.RemoveReader(email)
		//Json(context, true).End()
		models.Json(context, true).End()
		return
	}
	/*
	context.Layout("admin/cmd")
	context.Render("admin/cmd/reader", map[string]interface{}{
		"Title":   "读者",
		"Readers": models.GetReaders(),
	})
	*/
	data := map[string]interface{}{
		"Title":   "读者",
		"Readers": models.GetReaders(),
	}

	err := models.Theme(true).Layout("layout").Tpl("reader").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func AdminMessage(context *macaron.Context) {
	/*
	context.Layout("admin/cmd")
	context.Render("admin/cmd/message", map[string]interface{}{
		"Title":    "消息",
		"Messages": models.GetMessages(),
	})
	*/
	data := map[string]interface{}{
		"Title":    "消息",
		"Messages": models.GetMessages(),
	}

	err := models.Theme(true).Layout("layout").Tpl("message").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func AdminTemplates(context *macaron.Context) {
	if context.Req.Method == "POST" {
		change := context.Query("cache")
		if change != "" {
			// 设置主题缓存
			//models.SetThemeCache(context, change == "true")
			models.Json(context,true).End()
			return
		}
		theme := context.Query("theme")
		if theme != "" {
			models.SetSetting("site_theme", theme)
			models.SyncSettings()
			//Json(context, true).End()
			models.Json(context,true).End()
			return
		}
		return
	}

	data := map[string]interface{}{
		"Title":        "主题",
		"Themes":       models.GetThemes("templates"),
		"CurrentTheme": models.GetSetting("site_theme"),
	}

	err := models.Theme(true).Layout("layout").Tpl("templates").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

/*
func AdminLogs(context *macaron.Context) {
	if context.Req.Method == "DELETE" {
		models.RemoveLog(context.Query("file"))
		models.Json(context, true).End()
		return
	}
	data := map[string]interface{}{
		"Title": "日志",
		"Logs":  models.Logs(),
	}

	err := models.Theme(true).Layout("layout").Tpl("logs").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}
*/

func AdminBackup(context *macaron.Context) {
	if context.Req.Method == "POST" {
		file, e := models.DoBackup()
		if e != nil {
			models.Json(context, false).Set("msg", e.Error()).End()
			return
		}
		models.Json(context, true).Set("file", file).End()
		models.CreateMessage("backup", "[1]"+file)
		return
	}
	if context.Req.Method == "DELETE" {
		file :=context.Query("file")
		if file == "" {
			models.Json(context, false).End()
			return
		}
		models.RemoveBackupFile(file)
		models.Json(context, true).End()
		return
	}
	files, err := models.GetBackupFiles()
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}

	data := map[string]interface{}{
		"Files": files,
		"Title": "备份",
	}

	err = models.Theme(true).Layout("layout").Tpl("backup").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func AdminBackupFile(context *macaron.Context) {
	context.ServeFile("backup/" + context.Params("filename"))
}

func AdminMonitor(context *macaron.Context) {
	data := map[string]interface{}{
		"Title": "系统监控",
		"M":     models.ReadMemStats(),
	}

	err := models.Theme(true).Layout("layout").Tpl("monitor").Render(context, 200, data)
	if err != nil {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}