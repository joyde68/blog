package main

import (
	"github.com/joyde68/blog/models"
	"github.com/joyde68/blog/routes"
	"gopkg.in/macaron.v1"
	"os"
	"os/signal"
	"syscall"
)

var (
	// Global GoInk application
	App *macaron.Macaron
)

func init() {
	// init application
	App = macaron.Classic()

	// init storage
	models.Init()

	// set static handler
	App.Use(macaron.Static("public", macaron.StaticOptions{
		Prefix: "public",
	}))

	// set not found handler
	App.NotFound(func(context *macaron.Context) {
		models.Theme(false).Tpl("404").Render(context, 404, nil)
	})

	// set recover handler
	App.InternalServerError(func(context *macaron.Context) {
		models.Theme(false).Tpl("500").Render(context, 500, nil)
	})

	// load all data
	models.All()

	// catch exit command
	go catchExit()
}

func main() {

	registerAdminRoutes()
	registerHomeRoutes()

	App.Run()
}

func registerHomeRoutes() {
	App.Route("/login/", "GET,POST", routes.Login)
	App.Get("/logout/", routes.Logout)

	App.Get("/article/:slug", routes.Article)
	App.Get("/page/:slug/", routes.Page)
	App.Get("/p/:page/", routes.Home)
	App.Post("/comment/:id/", routes.Comment)
	App.Get("/tag/:tag/", routes.TagArticles)
	App.Get("/tag/:tag/p/:page/", routes.TagArticles)

	App.Get("/feed/", routes.Rss)
	App.Get("/sitemap", routes.SiteMap)

	App.Get("/:slug", routes.TopPage)
	App.Get("/", routes.Home)
}

func registerAdminRoutes() {
	// add admin handlers
	App.Group("/admin", func() {
		App.Get("/", routes.Admin)

		App.Route("/profile/", "GET,POST", routes.AdminProfile)

		App.Route("/password/", "GET,POST", routes.AdminPassword)

		App.Get("/articles/", routes.AdminArticle)
		App.Get("/articles/p/:page/", routes.AdminArticle)
		App.Route("/article/write/", "GET,POST", routes.ArticleWrite)
		App.Route("/article/:id/", "GET,POST,DELETE", routes.ArticleEdit)

		App.Get("/pages/", routes.AdminPage)
		App.Get("/pages/p/:page/", routes.AdminPage)
		App.Route("/page/write/", "GET,POST", routes.PageWrite)
		App.Route("/page/:id/", "GET,POST,DELETE", routes.PageEdit)

		App.Route("/comments/", "GET,POST,PUT,DELETE", routes.AdminComments)

		App.Route("/settings/", "GET,POST", routes.AdminSetting)
		App.Post("/settings/custom/", routes.CustomSetting)
		App.Post("/settings/nav/", routes.NavigatorSetting)

		App.Route("/files/", "GET,DELETE", routes.AdminFiles)
		App.Route("files/p/:page/", "GET,DELETE", routes.AdminFiles)
		App.Post("/files/upload/", routes.FileUpload)

		App.Route("/message/", "GET,POST,DELETE", routes.AdminMessage)
		App.Post("/message/read/", routes.Auth, routes.AdminMessageRead)

		App.Get("/monitor/", routes.AdminMonitor)

		App.Route("/reader/", "GET,POST", routes.AdminReader)

		// CMD Handler
		/*
		App.Get("/cmd/backup/file/", routes.CmdBackupFile)
		*/

		App.Route("/templates/", "GET,POST", routes.AdminTemplates)
		App.Route("/logs/", "GET,DELETE", routes.AdminLogs)
		App.Route("/backup/", "GET,POST,DELETE", routes.AdminBackup)
	}, routes.Auth)
}

// code from https://github.com/Unknwon/gowalker/blob/master/gowalker.go
func catchExit() {
	sigTerm := syscall.Signal(15)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, sigTerm)

	for {
		switch <-sig {
		case os.Interrupt, sigTerm:
			println("\n退出前保存数据")
			models.SyncAll()
			println("准备退出")
			os.Exit(0)
		}
	}
}