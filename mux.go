package sailgo

import (
	//"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

var RecoverPanic bool

type Router struct {
	regex          *regexp.Regexp
	params         map[int]string
	controllerType reflect.Type
}

type Mux struct {
	routers   []Router
	staticDir map[string]string
}

var method_func_map = map[string]string{
	"GET":     "Get",
	"PUT":     "Put",
	"POST":    "Post",
	"DELETE":  "Delete",
	"HEAD":    "Head",
	"PATCH":   "Patch",
	"OPTIONS": "Options",
}

func NewMux() *Mux {
	return &Mux{
		routers:   make([]Router, 0),
		staticDir: make(map[string]string, 10),
	}
}

func (m *Mux) SetStaticPath(prefix string, static_dir string) {
	m.staticDir[prefix] = static_dir
}

func (m *Mux) RegisterRouter(path string, controller interface{}) error {
	path_split := strings.Split(path, "/")
	//var params map[int]string
	params := make(map[int]string)
	j := 1
	for i, p := range path_split {
		expr := "([^/]*)" // '/' is not allowed, zero is not allowed
		if strings.HasPrefix(p, ":") {
			//'/usr/:id([0-9]+)', path_split: [^/]+, param: :id
			if index := strings.Index(p, "("); index != -1 {
				path_split[i] = p[index:]
				params[j] = p[:index]
			}
			path_split[i] = expr
			params[j] = p[1:]
			j++
		}
	}
	path = strings.Join(path_split, "/")
	//fmt.Println(path)
	ct := reflect.Indirect(reflect.ValueOf(controller)).Type() //get the real value not pointer
	regex, err := regexp.Compile(path)
	if err != nil {
		Critical("regexp compile failed:", err)
		os.Exit(1)
		return err
	}
	//fmt.Println(params)
	r := Router{
		regex:          regex,
		params:         params,
		controllerType: ct,
	}
	m.routers = append(m.routers, r)
	return nil
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			if !RecoverPanic {
				panic(err)
			} else {
				for i := 1; ; i += 1 {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					//fmt.Println(file, line)
					Critical(file, line)
				}
			}
		}
	}()
	//get the path from request, find the params and controller
	rurl := r.URL.Path
	//var params map[string]string
	params := make(map[string]string)
	var started bool
	//handle static file
	for prefix, static_dir := range m.staticDir {
		if strings.HasPrefix(rurl, prefix) {
			file := static_dir + rurl[len(prefix):]
			http.ServeFile(w, r, file)
			started = true
			return
		}
	}

	for _, router := range m.routers {
		if !router.regex.MatchString(rurl) { //url not match
			continue
		}
		matches := router.regex.FindStringSubmatch(rurl)
		//fmt.Println(matches)
		if len(router.params) > 0 {
			query := r.URL.Query()
			for i, param := range matches[1:] {
				query.Add(router.params[i], param)
				params[router.params[i+1]] = matches[i+1]
			}
			r.URL.RawQuery = url.Values(query).Encode() + "&" + r.URL.RawQuery
		}
		//fmt.Println(params)
		v := reflect.New(router.controllerType)
		in := make([]reflect.Value, 2)
		ct := &Context{ResponseWriter: w, Request: r, Params: params}
		in[0] = reflect.ValueOf(ct)
		in[1] = reflect.ValueOf(router.controllerType.Name())
		f := v.MethodByName("Init")
		f.Call(in)
		in = make([]reflect.Value, 0)

		fn := v.MethodByName("Addfunc")
		fn.Call(in)
		fn = v.MethodByName("Prepare")
		fn.Call(in)
		fn = v.MethodByName(method_func_map[r.Method])
		fn.Call(in)
		fn = v.MethodByName("Finish")
		fn.Call(in)
		started = true
		break
	}

	if started == false {
		http.NotFound(w, r)
	}
}
