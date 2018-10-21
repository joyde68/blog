package routes

import (
	"fmt"
	"github.com/joyde68/blog/models"
	"github.com/joyde68/blog/pkg"
	"gopkg.in/macaron.v1"
	"html/template"
	"strings"
	"time"
)

func SiteMap(context *macaron.Context) {
	baseUrl := models.GetSetting("site_url")
	fmt.Println(baseUrl)
	article, _ := models.GetPublishArticleList(1, 50)
	navigators := models.GetNavigators()
	now := time.Unix(pkg.Now(), 0).Format(time.RFC3339)

	articleMap := make([]map[string]string, len(article))
	for i, a := range article {
		m := make(map[string]string)
		m["Link"] = strings.Replace(baseUrl+a.Link(), baseUrl+"/", baseUrl, -1)
		m["Created"] = time.Unix(a.CreateTime, 0).Format(time.RFC3339)
		articleMap[i] = m
	}

	navMap := make([]map[string]string, 0)
	for _, n := range navigators {
		m := make(map[string]string)
		if n.Link == "/" {
			continue
		}
		if strings.HasPrefix(n.Link, "/") {
			m["Link"] = strings.Replace(baseUrl+n.Link, baseUrl+"/", baseUrl, -1)
		} else {
			m["Link"] = n.Link
		}
		m["Created"] = now
		navMap = append(navMap, m)
	}

	data := map[string]interface{}{
		"Title":      models.GetSetting("site_title"),
		"Link":       baseUrl,
		"Created":    now,
		"Articles":   articleMap,
		"Navigators": navMap,
	}

	context.Header().Set("Content-Type", "text/xml")
	t, err := template.New("template").Funcs(template.FuncMap{
		"Html": func(data string) template.HTML {
			return template.HTML(data)
		},
	}).ParseFiles("templates/sitemap.xml")
	if err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
	if err := t.ExecuteTemplate(context.Resp, "sitemap.xml", data); err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}

func Rss(context *macaron.Context) {
	baseUrl := models.GetSetting("site_url")
	article, _ := models.GetPublishArticleList(1, 20)
	author := models.GetUsersByRole("ADMIN")[0]

	articleMap := make([]map[string]string, len(article))
	for i, a := range article {
		m := make(map[string]string)
		m["Title"] = a.Title
		m["Link"] = strings.Replace(baseUrl+a.Link(), baseUrl+"/", baseUrl, -1)
		m["Author"] = author.Nick
		str := pkg.Markdown2Html(a.Content())
		str = strings.Replace(str, `src="/`, `src="`+strings.TrimSuffix(baseUrl, "/")+"/", -1)
		str = strings.Replace(str, `href="/`, `href="`+strings.TrimSuffix(baseUrl, "/")+"/", -1)
		m["Desc"] = str
		m["Created"] = time.Unix(a.CreateTime, 0).Format(time.RFC822)
		articleMap[i] = m
	}

	data := map[string]interface{}{
		"Title":    models.GetSetting("site_title"),
		"Link":     baseUrl,
		"Desc":     models.GetSetting("site_description"),
		"Created":  time.Unix(pkg.Now(), 0).Format(time.RFC822),
		"Articles": articleMap,
	}

	context.Header().Set("Content-Type", "application/rss+xml;charset=UTF-8")
	t, err := template.New("template").Funcs(template.FuncMap{
		"Html": func(data string) template.HTML {
			return template.HTML(data)
		},
	}).ParseFiles("templates/rss.xml")
	if err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
	if err := t.ExecuteTemplate(context.Resp, "rss.xml", data); err != nil {
		fmt.Println(err)
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	}
}
