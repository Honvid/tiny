package tiny

import (
	"net/http"
	"path"
	"regexp"
	"strings"
)

type RouterGroup struct {
	prefix      string
	middlewares []Handler
	parent      *RouterGroup
	engine      *Engine
}

type router struct {
	roots    map[string]*node
	handlers map[string]Handler
}

func (g *RouterGroup) Use(middlewares ...Handler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) GET(pattern string, handler Handler) {
	g.addRoute(http.MethodGet, pattern, handler)
}

func (g *RouterGroup) HEAD(pattern string, handler Handler) {
	g.addRoute(http.MethodHead, pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler Handler) {
	g.addRoute(http.MethodPost, pattern, handler)
}

func (g *RouterGroup) PUT(pattern string, handler Handler) {
	g.addRoute(http.MethodPut, pattern, handler)
}

func (g *RouterGroup) PATCH(pattern string, handler Handler) {
	g.addRoute(http.MethodPatch, pattern, handler)
}

func (g *RouterGroup) OPTIONS(pattern string, handler Handler) {
	g.addRoute(http.MethodOptions, pattern, handler)
}

func (g *RouterGroup) DELETE(pattern string, handler Handler) {
	g.addRoute(http.MethodDelete, pattern, handler)
}

func (g *RouterGroup) CONNECT(pattern string, handler Handler) {
	g.addRoute(http.MethodConnect, pattern, handler)
}

func (g *RouterGroup) TRACE(pattern string, handler Handler) {
	g.addRoute(http.MethodTrace, pattern, handler)
}

func (g *RouterGroup) Static(relativePath string, root string) {
	handler := g.staticHandler(relativePath, http.Dir(root))
	url := path.Join(relativePath, "/*filepath")
	g.GET(url, handler)
}

func (g *RouterGroup) Any(pattern string, handler Handler) {
	g.GET(pattern, handler)
	g.POST(pattern, handler)
	g.PUT(pattern, handler)
	g.DELETE(pattern, handler)
	g.PATCH(pattern, handler)
	g.HEAD(pattern, handler)
	g.OPTIONS(pattern, handler)
	g.TRACE(pattern, handler)
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]Handler),
	}
}

func (g *RouterGroup) addRoute(method string, comp string, handler Handler) {
	pattern := g.prefix + comp
	g.engine.router.addRouter(method, pattern, handler)
}

func (r *router) addRouter(method string, pattern string, handler Handler) {
	segments := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, segments, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchSegments := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchSegments, 0)
	if n != nil {
		segments := parsePattern(n.pattern)
		for index, segment := range segments {
			if segment[0] == ':' {
				params[segment[1:]] = searchSegments[index]
			}
			if segment[0] == '*' && len(segment) > 1 {
				params[segment[1:]] = strings.Join(searchSegments[index:], "/")
				break
			}
			// 正则匹配
			if segment[0] == '{' && segment[len(segment)-1] == '}' {
				splitPart := strings.Split(segment[1:len(segment)-1], ":")
				rePattern := splitPart[1]
				if rePattern[0] != '^' {
					rePattern = "^" + rePattern
				}
				if rePattern[len(rePattern)-1] != '$' {
					rePattern = rePattern + "$"
				}
				re := regexp.MustCompile(rePattern)
				if re.MatchString(searchSegments[index]) {
					params[splitPart[0]] = searchSegments[index]
				} else {
					return nil, nil
				}
			}
		}
		return n, params
	}
	return nil, nil
}

func (g *RouterGroup) staticHandler(relativePath string, fs http.FileSystem) Handler {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Segment("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (r *router) handle(c *Context) {
	n, segments := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Segments = segments
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}

func parsePattern(pattern string) []string {
	patterns := strings.Split(pattern, "/")
	segments := make([]string, 0)
	for _, segment := range patterns {
		if segment != "" {
			segments = append(segments, segment)
			if segment[0] == '*' {
				break
			}
		}
	}
	return segments
}
