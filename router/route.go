/**
 * Creator: hevi
 * Time: 2018/12/3 09:59
 * Description: 路由接口
 */

package router

import (
	"github.com/Hevi-Ye/jerry/context"
	"net/http"
	"os"
)

type handleFun func(*context.Context)
type middlewareFun func(*context.Context) bool

type handle struct {
	h   handleFun
	mws []middlewareFun
}

var handles = make(map[string]map[string]handle)

type Route struct {
	middlewareList []middlewareFun
}

func (r *Route) any(pattern, httpMethod string, h handleFun) {
	if _, ok := handles[pattern]; !ok {
		handles[pattern] = make(map[string]handle)
	}

	handles[pattern][httpMethod] = handle{h: h, mws: r.middlewareList}
}

func (r *Route) Get(pattern string, h handleFun) {
	r.any(pattern, http.MethodGet, h)
}

func (r *Route) Post(pattern string, h handleFun) {
	r.any(pattern, http.MethodPost, h)
}

type RouteMux struct {
	staticPath string
}

func NewRouteMux() *RouteMux {
	return &RouteMux{}
}

func (rt *RouteMux) Run(addr string) error {
	return http.ListenAndServe(addr, rt)
}

func (rt *RouteMux) Static(path string) {
	rt.staticPath = path
}

func (rt *RouteMux) serveFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if len(path) == 0 || path == "/" {
		path = "index.html"
	}

	file := rt.staticPath + path

	f, err := os.Stat(file)
	if err != nil || f.IsDir() {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, file)
}

func (rt *RouteMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m, ok := handles[r.URL.Path]
	if !ok {
		rt.serveFile(w, r)
		return
	}

	hd, ok := m[r.Method]
	if !ok {
		http.NotFound(w, r)
		return
	}

	r.ParseForm()

	ctx := context.NewContext(w, r)

	for _, mw := range hd.mws {
		if !mw(ctx) {
			return
		}
	}

	hd.h(context.NewContext(w, r))
}

func (rt *RouteMux) Group(ms ...middlewareFun) *Route {
	r := &Route{
		middlewareList: make([]middlewareFun, 0),
	}

	r.middlewareList = append(r.middlewareList, ms...)

	return r
}
