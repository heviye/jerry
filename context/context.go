/**
 * Creator: hevi
 * Time: 2018/12/3 09:25
 * Description: http 请求上下文
 */

package context

import "net/http"

type Context struct {
	r *http.Request
	w http.ResponseWriter
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w: w, r: r}
}

func (c *Context) Write(b []byte) (int, error) {
	return c.w.Write(b)
}
