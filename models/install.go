package models

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/Unknwon/cae/zip"
	"github.com/joyde68/blog/pkg"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

var (
	tmpZipFile      = "tmp.zip"
	installLockFile = "install.lock"
	createdCommentTemplate = `
<table style="width: 99.8%;height:99.8% "><tbody><tr><td style="background:#FAFAFA">
    <div style="background-color:white;border-top:2px solid #0079BC;box-shadow:0 1px 3px #AAAAAA;line-height:180%;padding:0 15px 12px;width:500px;margin:50px auto;color:#555555;font-family:'Century Gothic','Trebuchet MS','Hiragino Sans GB',微软雅黑,'Microsoft Yahei',Tahoma,Helvetica,Arial,'SimSun',sans-serif;font-size:12px;">
        <h2 style="border-bottom:1px solid #DDD;font-size:14px;font-weight:normal;padding:13px 0 10px 8px;"><span style="color: #0079bc;font-weight: bold;">&gt;</span>{{.author}}&nbsp;在文章《{{.title}}》发表了评论！</h2>
        <div style="padding:0 12px 0 12px;margin-top:18px">
            <p><strong>{{.author}}</strong>&nbsp;同学，在文章《{{.title}}》上发表评论:</p>
            <p style="background-color: #f5f5f5;border: 0px solid #DDD;padding: 10px 15px;margin:18px 0">{{.text}}</p>
            <p>您可以点击 <a style="text-decoration:none; color:#e64346" href="{{.permalink}}">查看完整內容 </a>，欢迎再次光临 <a style="text-decoration:none; color:#0079bc" href="{{.link}}">{{.site}}</a>。</p>
        </div>
    </div>
</td></tr></tbody></table>
`
	replyCommentTemplate = `
<table style="width: 99.8%;height:99.8% "><tbody><tr><td style="background:#FAFAFA">
    <div style="background-color:white;border-top:2px solid #0079BC;box-shadow:0 1px 3px #AAAAAA;line-height:180%;padding:0 15px 12px;width:500px;margin:50px auto;color:#555555;font-family:'Century Gothic','Trebuchet MS','Hiragino Sans GB',微软雅黑,'Microsoft Yahei',Tahoma,Helvetica,Arial,'SimSun',sans-serif;font-size:12px;">
        <h2 style="border-bottom:1px solid #DDD;font-size:14px;font-weight:normal;padding:13px 0 10px 8px;"><span style="color: #0079bc;font-weight: bold;">&gt; </span>您的留言有回复啦！</h2>
        <div style="padding:0 12px 0 12px;margin-top:18px">
            <p><strong>{{.author_p}}</strong>&nbsp;同学，您曾在文章《{{.title}}》上发表评论:</p>
            <p style="background-color: #f5f5f5;border: 0px solid #DDD;padding: 10px 15px;margin:18px 0">{{.text_p}}</p>
            <p><strong>{{.author}}</strong>&nbsp;给您的回复如下:</p>
            <p style="background-color: #f5f5f5;border: 0px solid #DDD;padding: 10px 15px;margin:18px 0">{{.text}}</p>
            <p>您可以点击 <a style="text-decoration:none; color:#e64346" href="{{.permalink}}">查看回复的完整內容 </a>，欢迎再次光临 <a style="text-decoration:none; color:#0079bc" href="{{.link}}">{{.site}}</a>。</p>
        </div>
    </div>
</td></tr></tbody></table>
`
)

func CheckInstall() bool {
	return pkg.IsFile(installLockFile)
}

func writeDefaultData() {
	// write user
	u := new(User)
	u.Id = Storage.TimeInc(10)
	u.Name = "admin"
	u.Password = pkg.Sha1("adminxxxxx")
	u.Nick = "管理员"
	u.Email = "admin@example.com"
	u.Url = "http://example.com/"
	u.CreateTime = pkg.Now()
	u.Bio = "这是站点的管理员，你可以添加一些个人介绍，支持换行不支持markdown"
	u.LastLoginTime = u.CreateTime
	u.Role = "ADMIN"
	Storage.Set("users", []*User{u})

	// write token
	Storage.Set("tokens", map[string]*Token{})

	// write contents
	a := new(Content)
	a.Id = Storage.TimeInc(9)
	a.Title = "欢迎使用 Fxh.Go"
	a.Slug = "welcome-fxh-go"
	a.Text = "如果您看到这篇文章,表示您的 blog 已经安装成功."
	a.Tags = []string{"Fxh.Go"}
	a.CreateTime = pkg.Now()
	a.EditTime = a.CreateTime
	a.UpdateTime = a.CreateTime
	a.IsComment = true
	a.IsLinked = false
	a.AuthorId = u.Id
	a.Type = "article"
	a.Status = "publish"
	a.Format = "markdown"
	a.Template = "blog.html"
	a.Hits = 1
	// write comments
	co := new(Comment)
	co.Author = u.Nick
	co.Email = u.Email
	co.Url = u.Url
	co.Content = "欢迎加入使用 Fxh.Go"
	co.Avatar = pkg.Gravatar(co.Email, "50")
	co.Pid = 0
	co.Ip = "127.0.0.1"
	co.UserAgent = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.17 (KHTML, like Gecko) Chrome/24.0.1312.57 Safari/537.17"
	co.IsAdmin = true
	co.Id = Storage.TimeInc(7)
	co.CreateTime = pkg.Now()
	co.Status = "approved"
	co.Cid = a.Id
	a.Comments = []*Comment{co}
	Storage.Set("content/article-"+strconv.Itoa(a.Id), a)

	// write pages
	p := new(Content)
	p.Id = a.Id + Storage.TimeInc(6)
	p.Title = "关于"
	p.Slug = "about-me"
	p.Text = "本页面由 Fxh.Go 创建, 这只是个测试页面."
	p.Tags = []string{}
	p.CreateTime = pkg.Now()
	p.EditTime = p.CreateTime
	p.UpdateTime = p.CreateTime
	p.IsComment = true
	p.IsLinked = true
	p.AuthorId = u.Id
	p.Type = "page"
	p.Status = "publish"
	p.Format = "markdown"
	p.Comments = make([]*Comment, 0)
	p.Template = "page.html"
	p.Hits = 1
	Storage.Set("content/page-"+strconv.Itoa(p.Id), p)
	p2 := new(Content)
	p2.Id = p.Id + Storage.TimeInc(6)
	p2.Title = "好友"
	p2.Slug = "friends"
	p2.Text = "本页面由 Fxh.Go 创建, 这只是个测试页面."
	p2.Tags = []string{}
	p2.CreateTime = pkg.Now()
	p2.EditTime = p2.CreateTime
	p2.UpdateTime = p2.CreateTime
	p2.IsComment = true
	p2.IsLinked = true
	p2.AuthorId = u.Id
	p2.Type = "page"
	p2.Status = "publish"
	p2.Format = "markdown"
	p2.Comments = make([]*Comment, 0)
	p2.Template = "page.html"
	p2.Hits = 1
	Storage.Set("content/page-"+strconv.Itoa(p2.Id), p2)

	// write new reader
	Storage.Set("readers", map[string]*Reader{})

	// write version
	/*
	v := new(version)
	v.Name = "Fxh.Go"
	v.BuildTime = pkg.pkg.Now()
	v.Version = appVersion
	Storage.Set("version", v)
	*/

	// write settings
	s := map[string]string{
		"site_title":         "Fxh.Go",
		"site_sub_title":     "Go开发的简单博客",
		"site_keywords":      "Fxh.Go,Golang,Blog",
		"site_description":   "Go语言开发的简单博客程序",
		"site_url":           "http://localhost/",
		"article_size":       "4",
		"popular-size": "4",
		"recent-comment-size": "4",
		"site_theme":         "default",
		"enable_go_markdown": "false",
		"c_footer_weibo":     "#",
		"c_footer_github":    "#",
		"c_footer_email":     "#",
		"c_home_avatar":      "/public/img/site.png",
		"c_footer_ga":        "<!-- google analytics or other -->",
		"create-comment-template": createdCommentTemplate,
		"reply-comment-template": replyCommentTemplate,
	}
	Storage.Set("settings", s)

	// write files
	Storage.Set("files", []*File{})

	// write message
	Storage.Set("messages", []*Message{})

	// write navigators
	n := new(NavItem)
	n.Order = 1
	n.Text = "文章"
	n.Title = "文章"
	n.Link = "/"
	n2 := new(NavItem)
	n2.Order = 2
	n2.Text = "关于"
	n2.Title = "关于"
	n2.Link = "/about-me.html"
	n3 := new(NavItem)
	n3.Order = 3
	n3.Text = "好友"
	n3.Title = "好友"
	n3.Link = "/friends.html"
	Storage.Set("navigators", []*NavItem{n, n2, n3})

	// write default tmp data
	writeDefaultTmpData()
}

func writeDefaultTmpData() {
	TmpStorage.Set("contents", make(map[string][]int))
}

func DoInstall() {
	// init some settings
	os.MkdirAll(path.Join("data", "log"), 0755)
	os.MkdirAll(path.Join("tmp", "data"), 0755)
	os.MkdirAll(path.Join("public", "upload"), 0755)

	os.Mkdir(Storage.Dir, 0755)
	os.Mkdir(Storage.Dir + "/content", 0755)
	
	writeDefaultData()

	ExtractBundleBytes()

	ioutil.WriteFile(installLockFile, []byte(fmt.Sprint(pkg.Now())), os.ModePerm)
	println("install success")
}


func ExtractBundleBytes() {
	// origin from https://github.com/wendal/gor/blob/master/gor/gor.go
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(pkg.ZipBytes))
	b, _ := ioutil.ReadAll(decoder)
	ioutil.WriteFile(tmpZipFile, b, 0755)
	z, e := zip.Open(tmpZipFile)
	if e != nil {
		panic(e)
		os.Exit(1)
	}
	defer func() {
		z.Close()
		decoder = nil
		os.Remove(tmpZipFile)
	}()
	z.ExtractTo("")
}

/*
func DoUpdateZipBytes(file string) error {
	// copy from https://github.com/wendal/gor/blob/master/gor/gor.go
	bytes, _ := ioutil.ReadFile(file)
	zipWriter, _ := os.OpenFile("app/cmd/zip.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	header := `package cmd
const zipBytes="`
	zipWriter.Write([]byte(header))
	encoder := base64.NewEncoder(base64.StdEncoding, zipWriter)
	encoder.Write(bytes)
	encoder.Close()
	zipWriter.Write([]byte(`"`))
	zipWriter.Sync()
	zipWriter.Close()
	println("update success")
	return nil
}
*/