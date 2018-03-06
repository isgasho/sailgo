package sailgo

import (
	"html/template"
	//"log"
	"net/http"
	"path"
	"strings"
)

var ViewsPath = "templates"

type Controller struct {
	Ct        *Context
	FuncMap   map[string]func()
	Tpl       *template.Template
	Data      map[interface{}]interface{}
	ChildName string
	TplNames  string
	Layout    []string
	TplExt    string
}

type ControllerInterface interface {
	Init(ct *Context, cn string)
	Prepare()
	Get()
	Post()
	Put()
	Delete()
	Head()
	Patch()
	Finish()
	Render() error //执行完method对应的方法之后渲染页面
}

func (c *Controller) Init(ct *Context, cn string) {
	c.Data = make(map[interface{}]interface{})
	c.FuncMap = make(map[string]func())
	c.Layout = make([]string, 0)
	c.TplNames = ""
	c.ChildName = cn
	c.Ct = ct
	c.TplExt = "tpl"
}

func (c *Controller) Addfunc() {
}

func (c *Controller) Handle(method string, f func()) {
	c.FuncMap[strings.ToLower(method)] = f
}

func (c *Controller) Prepare() {
}

func (c *Controller) Finish() {
}

func (c *Controller) Get() {
	Info(c.FuncMap)
	if f, ok := c.FuncMap["get"]; ok {
		f()
		return
	}
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Post() {
	if f, ok := c.FuncMap["post"]; ok {
		f()
		return
	}
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Put() {
	if f, ok := c.FuncMap["put"]; ok {
		f()
		return
	}
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Delete() {
	if f, ok := c.FuncMap["delete"]; ok {
		f()
		return
	}
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Head() {
	if f, ok := c.FuncMap["head"]; ok {
		f()
		return
	}
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Patch() {
	if f, ok := c.FuncMap["patch"]; ok {
		f()
		return
	}
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Options() {
	if f, ok := c.FuncMap["options"]; ok {
		f()
		return
	}
	http.Error(c.Ct.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Render() error {
	if len(c.Layout) > 0 {
		var filenames []string
		for _, file := range c.Layout {
			filenames = append(filenames, path.Join(ViewsPath, file))
		}
		t, err := template.ParseFiles(filenames...)
		if err != nil {
			//log.Println("template parse files err:", err)
			Critical("template parse files err:", err)
			return err
		}
		err = t.ExecuteTemplate(c.Ct.ResponseWriter, c.TplNames, c.Data)
		if err != nil {
			//log.Println("template execute err:", err)
			Critical("template execute err:", err)
			return err
		}
	} else {
		if c.TplNames == "" {
			c.TplNames = c.ChildName + "/" + c.Ct.Request.Method + "." + c.TplExt
		}
		t, err := template.ParseFiles(path.Join(ViewsPath, c.TplNames))
		if err != nil {
			Critical("template parse files err:", err)
			//log.Println("template parse files err:", err)
			return err
		}
		err = t.Execute(c.Ct.ResponseWriter, c.Data)
		if err != nil {
			Critical("template execute err:", err)
			//log.Println("template execute err:", err)
			return err
		}
	}
	return nil
}
