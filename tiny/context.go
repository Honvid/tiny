package tiny

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type H map[string]interface{}

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	Path     string
	Method   string
	Segments map[string]string

	StatusCode int
	handlers   []Handler
	index      int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Path:    r.URL.Path,
		Method:  r.Method,
		index:   -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Segment(key string) string {
	value, _ := c.Segments[key]
	return value
}

func (c *Context) Post(param string) string {
	return c.Request.FormValue(param)
}

func (c *Context) Query(param string) string {
	return c.Request.URL.Query().Get(param)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *Context) Fail(code int, message string)  {
	c.StatusCode = code
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.Write([]byte(fmt.Sprint(message)))
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/html")
	c.Writer.Write([]byte(html))
}
