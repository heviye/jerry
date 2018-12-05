/**
 * Creator: hevi
 * Time: 2018/12/3 09:25
 * Description: http 请求上下文
 */

package context

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"net/http"
	"sync"
)

type Context struct {
	sync.RWMutex

	r *http.Request
	w http.ResponseWriter

	items map[string]interface{}

	printLog func(string)
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:     w,
		r:     r,
		items: make(map[string]interface{}),
	}
}

func (c *Context) Write(b []byte) (int, error) {
	return c.w.Write(b)
}

func (c *Context) WriteHeader(statusCode int) {
	c.w.WriteHeader(statusCode)
}

func (c *Context) Header() http.Header {
	return c.w.Header()
}

func (c *Context) SetPrintLogFunc(fn func(string)) {
	c.printLog = fn
}

func (c *Context) Request() *http.Request {
	return c.r
}

func (c *Context) RemoteAddr() string {
	return c.r.RemoteAddr
}

func (c *Context) GetCookie(name string) (*http.Cookie, error) {
	return c.r.Cookie(name)
}

func (c *Context) SetCookie(cooke *http.Cookie) {
	http.SetCookie(c.w, cooke)
}

// 获取参与的值
// GET、POST、PUT、DELETE都可以通过该函数获取
func (c *Context) Query(param string) string {
	return c.r.Form.Get(param)
}

// 增加用户自定义的值
func (c *Context) Add(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()

	c.items[key] = value
}

// 用户自定义的值是否存在
func (c *Context) IsExists(key string) bool {
	c.RLock()
	defer c.RUnlock()

	_, ok := c.items[key]

	return ok
}

// 获取用户自定义的值
func (c *Context) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()

	return c.items[key]
}

// 删除用户自定义的值
func (c *Context) Del(key string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.items[key]; ok {
		delete(c.items, key)
	}
}

func (c *Context) log(content string) {
	if c.printLog == nil {
		return
	}

	c.printLog(content)
}

// 写入PROTOBUF格式的数据
func (c *Context) WritePB(pb proto.Message) (int, error) {
	buf, err := proto.Marshal(pb)
	if err != nil {
		return 0, err
	}

	c.log(string(buf))

	return c.w.Write(buf)
}

// 写入JSON格式的数据
func (c *Context) WriteJson(data interface{}) (int, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	c.log(string(buf))

	return c.w.Write(buf)
}
