package mygin

import (
	"io"
	"net/http"
	"os"
	"strings"
)

type RouterGroup struct {
	Handlers HandlerFuncChain
	fullPath string
	engine   *Engine
}

func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) (g *RouterGroup) {
	absolutePath := group.calcAbsolutePath(relativePath)
	if absolutePath == "/" {
		g = &group.engine.RouterGroup
	} else {
		g = &RouterGroup{
			fullPath: absolutePath,
			engine:   group.engine,
		}
	}
	mergedHandlers := append(group.Handlers, handlers...)
	g.Handlers = mergedHandlers
	return g
}

func (group *RouterGroup) Use(handlers ...HandlerFunc) {
	group.Handlers = append(group.Handlers, handlers...)
}

func (group *RouterGroup) Handle(method, relativePath string, handlers ...HandlerFunc) {
	absolutePath := group.calcAbsolutePath(relativePath)
	mergedHandlers := append(group.Handlers, handlers...)
	group.engine.addRoute(method, absolutePath, mergedHandlers)
	debugPrint(method+"\t"+absolutePath+"\t--> (%v handlers)", len(mergedHandlers))
}

func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) {
	group.Handle(http.MethodGet, relativePath, handlers...)
}

func (group *RouterGroup) POST(path string, handlers ...HandlerFunc) {
	group.Handle(http.MethodPost, path, handlers...)
}

func (group *RouterGroup) PUT(path string, handlers ...HandlerFunc) {
	group.Handle(http.MethodPut, path, handlers...)
}

func (group *RouterGroup) DELETE(path string, handlers ...HandlerFunc) {
	group.Handle(http.MethodDelete, path, handlers...)
}

func (group *RouterGroup) HEAD(path string, handlers ...HandlerFunc) {
	group.Handle(http.MethodHead, path, handlers...)
}

func (group *RouterGroup) OPTIONS(path string, handlers ...HandlerFunc) {
	group.Handle(http.MethodOptions, path, handlers...)
}

func (group *RouterGroup) PATCH(path string, handlers ...HandlerFunc) {
	group.Handle(http.MethodPatch, path, handlers...)
}

func (group *RouterGroup) calcAbsolutePath(relativePath string) string {
	var absolutePath string
	if relativePath == "/" {
		absolutePath = group.fullPath
	} else if group.fullPath == "/" {
		absolutePath = relativePath
	} else {
		absolutePath = group.fullPath + relativePath
	}
	return absolutePath
}

func (group *RouterGroup) Static(relativePath string, root string) {
	absolutePath := group.calcAbsolutePath(relativePath)
	handler := func(c *Context) {
		file, err := os.Open(root + strings.Replace(c.Path, absolutePath, "", 1))
		if err != nil {
			notFoundHandler(c)
		} else {
			io.Copy(c.Writer, file)
		}
		defer file.Close()
	}
	group.GET(relativePath+"/*", handler)
}
