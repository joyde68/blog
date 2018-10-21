package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joyde68/blog/pkg"
	"gopkg.in/macaron.v1"
	"html/template"
	"path"
)

// 返回模板渲染后的文本
func RenderText(name string, data map[string]interface{}) string {
	t, err := template.New("template").Funcs(template.FuncMap{
		"Html": func(data string) template.HTML {
			return template.HTML(data)
		},
		"DateInt64":  pkg.DateInt64,
		"DateString": pkg.DateString,
		"DateTime":   pkg.DateTime,
		"Now":        pkg.Now,
		"Html2str":   pkg.Html2str,
		"FileSize":   pkg.FileSize,
		"Setting":    GetSetting,
		"Navigator":  GetNavigators,
		"Md2html":    pkg.Markdown2HtmlTemplate,
	}).ParseFiles(path.Join("templates", "default", name+".tmpl"))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var contentHtml bytes.Buffer
	err = t.ExecuteTemplate(&contentHtml,name + ".tmpl", data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return contentHtml.String()
}









type jsonContext struct {
	context *macaron.Context
	data    map[string]interface{}
}

// Json creates a json context response.
func Json(context *macaron.Context, res bool) *jsonContext {
	c := new(jsonContext)
	c.context = context
	c.data = make(map[string]interface{})
	c.data["res"] = res
	return c
}

func (jc *jsonContext) Set(key string, v interface{}) *jsonContext {
	jc.data[key] = v
	return jc
}

func (jc *jsonContext) End() {
	jc.context.Resp.Header().Set("Content-Type","application/json;charset=UTF-8")
	dataJson, err := json.Marshal(jc.data)
	if err != nil {
		fmt.Println(err)
		jc.context.Resp.Write([]byte(`{"res":false"}`))
	}
	jc.context.Write(dataJson)
}

type themeContext struct {
	template   string
	layout string
	tpl string
}

// Theme creates themed context response.
func Theme(isAdmin bool) *themeContext {
	t := new(themeContext)
	t.template = GetSetting("site_theme")
	if isAdmin {
		t.template = "admin"
	}
	return t
}

func (tc *themeContext) Layout(layout string) *themeContext  {
	/*
		if layout == "" {
			context.Layout("")
			return tc
		}
		context.Layout(path.Join(tc.template, layout))
	*/
	tc.layout = layout
	return tc
}

func (tc *themeContext) Tpl(tpl string) *themeContext {
	//return context.Tpl(path.Join(tc.template, tpl), data)
	tc.tpl = tpl
	return tc
}
/*
func (tc *themeContext) Has(tpl string) bool {
	file := path.Join(tc.theme, tpl)
	return context.App().View().Has(file)
}
*/
func (tc *themeContext) Render(context *macaron.Context, statusCode int, data map[string]interface{}) error {
	context.Resp.WriteHeader(statusCode)
	//context.Render(path.Join(tc.theme, tpl), data)
	t := template.New("template").Funcs(template.FuncMap{
		"Html": func(data string) template.HTML {
			return template.HTML(data)
		},
		"DateInt64":  pkg.DateInt64,
		"DateString": pkg.DateString,
		"DateTime":   pkg.DateTime,
		"Now":        pkg.Now,
		"Html2str":   pkg.Html2str,
		"FileSize":   pkg.FileSize,
		"Setting":    GetSetting,
		"Navigator":  GetNavigators,
		"Md2html":    pkg.Markdown2HtmlTemplate,
	})

	if tc.layout == "" || pkg.IsFile(tc.layout) {
		t, err := t.ParseFiles(path.Join("templates", tc.template, tc.tpl+".tmpl"))
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = t.ExecuteTemplate(context.Resp, tc.tpl+".tmpl",data)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}

	t, err := t.ParseFiles(path.Join("templates", tc.template, tc.layout+".tmpl"), path.Join("templates", tc.template, tc.tpl+".tmpl"))
	if err != nil {
		fmt.Println(err)
		return err
	}

	var contentHtml bytes.Buffer
	err = t.ExecuteTemplate(&contentHtml, tc.tpl+".tmpl", data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	data["LayoutContent"] = contentHtml.String()
	err = t.ExecuteTemplate(context.Resp, tc.layout+".tmpl", data)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}