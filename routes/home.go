package routes

import (
	"fmt"
	"github.com/joyde68/blog/models"
	"github.com/joyde68/blog/pkg"
	"gopkg.in/macaron.v1"
	"net/url"
	"strconv"
	"strings"
)

func Login(context *macaron.Context) {
	ip := context.RemoteAddr()
	if logginErrorNumber := models.GetLoginErrCount(ip); logginErrorNumber > 3 {
		context.Redirect("/")
	}
	if context.Req.Method == "POST" {
		user := models.GetUserByName(context.Query("user"))
		if user == nil {
			models.AddLoginErrLog("用户名错误", context)
			models.Json(context, false).End()
			return
		}
		if !user.CheckPassword(context.Query("password")) {
			models.AddLoginErrLog("密码错误", context)
			models.Json(context, false).End()
			return
		}
		exp := 3600 * 24 * 3
		expStr := strconv.Itoa(exp)
		s := models.CreateToken(user, context, int64(exp))
		context.SetCookie("token-user", strconv.Itoa(s.UserId), expStr)
		context.SetCookie("token-value", s.Value, expStr)
		models.Json(context, true).End()
		return
	}
	if context.GetCookie("token-value") != "" {
		context.Redirect("/admin/")
		return
	}
	err := models.Theme(true).Tpl("login").Render(context, 200, nil)
	if err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func Auth(context *macaron.Context) {
	tokenValue := context.GetCookie("token-value")
	token := models.GetTokenByValue(tokenValue)
	if token == nil {
		context.Redirect("/logout/")
		return
	}
	if !token.IsValid() {
		context.Redirect("/logout/")
		return
	}
}

func Logout(context *macaron.Context) {
	context.SetCookie("token-user", "", "-3600")
	context.SetCookie("token-value", "", "-3600")
	context.Redirect("/login/")
}

func TagArticles(context *macaron.Context) {
	page := context.ParamsInt("page")
	tag, _ := url.QueryUnescape(context.Params("tag"))
	size := getArticleListSize()
	articles, pager := models.GetTaggedArticleList(tag, page, getArticleListSize())
	// fix dotted tag
	if len(articles) < 1 && strings.Contains(tag, "-") {
		articles, pager = models.GetTaggedArticleList(strings.Replace(tag, "-", ".", -1), page, size)
	}

	data := map[string]interface{}{
		"Articles": articles,
		"Pager":    pager,
		"Tag":      tag,
		"Title":    tag,
	}

	err := models.Theme(false).Layout("layout").Tpl("list").Render(context, 200, data)
	if err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func Home(context *macaron.Context) {
	page := context.ParamsInt("page")
	articles, pager := models.GetPublishArticleList(page, getArticleListSize())
	data := map[string]interface{}{
		"Articles":    articles,
		"Pager":       pager,
		"SidebarHtml": "SidebarHtml",
	}
	if page > 1 {
		data["Title"] = "第 " + strconv.Itoa(page) + " 页"
	}

	err := models.Theme(false).Layout("layout").Tpl("list").Render(context, 200, data)
	if err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func Article(context *macaron.Context) {
	slug := context.Params("slug")
	article := models.GetContentBySlug(slug)
	if article == nil || article.Status != "publish" {
		models.Theme(false).Tpl("404").Render(context, 404, nil)
		return
	}
	if article.Type != "article" {
		models.Theme(false).Tpl("404").Render(context, 404, nil)
		return
	}

	data := map[string]interface{}{
		"Title":   article.Title,
		"Article": article,
		"CommentHtml": models.RenderText("comment", map[string]interface{}{
			"Content":  article,
			"Comments": article.Comments,
		}),
	}

	// 渲染内容
	err := models.Theme(false).Layout("layout").Tpl("article").Render(context, 200, data)
	if err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
		return
	}

	article.Hits++
}

func Page(context *macaron.Context) {
	slug := context.Params("slug")
	article := models.GetContentBySlug(slug)
	if article == nil || article.Status != "publish" {
		models.Theme(false).Tpl("404").Render(context, 404, nil)
		return
	}
	if article.Type != "page" {
		models.Theme(false).Tpl("404").Render(context, 404, nil)
		return
	}

	data := map[string]interface{}{
		"Title": article.Title,
		"Page":  article,
		//"CommentHtml": Comments(context, article),
	}

	err := models.Theme(false).Layout("layout").Tpl("page").Render(context, 200, data)
	if err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
		return
	}

	article.Hits++
}

func TopPage(context *macaron.Context) {
	slug := context.Params("slug")
	page := models.GetContentBySlug(slug)
	if page == nil || page.Status != "publish" {
		models.Theme(false).Tpl("404").Render(context, 404, nil)
		return
	}

	if !page.IsLinked || page.Type != "page" {
		models.Theme(false).Tpl("404").Render(context, 404, nil)
		return
	}

	data := map[string]interface{}{
		"Title": page.Title,
		"Page":  page,
	}

	err := models.Theme(false).Layout("layout").Tpl("page").Render(context, 200, data)
	if err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
		return
	}

	page.Hits++
}

func Comment(context *macaron.Context) {
	cid := context.ParamsInt("id")
	if cid < 1 {
		models.Json(context, false).End()
		return
	}
	if models.GetContentById(cid) == nil {
		models.Json(context, false).End()
		return
	}

	msg := validateComment(context)
	if msg != "" {
		models.Json(context, false).Set("msg", msg).End()
		return
	}
	co := new(models.Comment)
	co.Author = context.Query("user")
	co.Email = context.Query("email")
	co.Url = context.Query("url")
	co.Content = context.Query("content")
	co.Avatar = pkg.Gravatar(co.Email, "50")
	co.Pid = context.QueryInt("pid")
	co.Ip = context.RemoteAddr()
	co.UserAgent = context.Req.Header.Get("User-Agent")
	co.IsAdmin = false
	models.CreateComment(cid, co)
	models.Json(context, true).Set("comment", co.ToJson()).End()
	models.CreateMessage("comment", co)

	go models.SendEmail(co)
}

func validateComment(context *macaron.Context) string {
	if pkg.IsEmptyString(context.Query("user")) || pkg.IsEmptyString(context.Query("content")) {
		return "称呼，邮箱，内容必填"
	}
	if !pkg.IsEmail(context.Query("email")) {
		return "邮箱格式错误"
	}
	if !pkg.IsEmptyString(context.Query("url")) && !pkg.IsURL(context.Query("url")) {
		return "网址格式错误"
	}
	return ""
}

func getArticleListSize() int {
	size, _ := strconv.Atoi(models.GetSetting("article_size"))
	if size < 1 {
		size = 5
	}
	return size
}
