package main

import (
	"embed"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed public
var public embed.FS

var exts = map[string]bool{
	".mobi": true,
	".epub": true,
	".text": true,
	".azw3": true,
	".pdf":  true,
	".azw":  true,
	".txt":  true,
	".doc":  true,
	".docx": true,
}

type FileItem struct {
	Name string
	Path string
	Size string
}

type Template struct {
	templates *template.Template
}

func NewTemplate() *Template {
	tt, err := template.ParseFS(public, "public/*.gohtml")
	if err != nil {
		panic(err)
	}
	return &Template{templates: tt}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, fmt.Sprintf("%s.gohtml", name), data)
}

func main() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("访问地址: ")
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println(fmt.Sprintf("http://%s:1614", ipnet.IP.String()))
			}
		}
	}

	// Echo instance
	e := echo.New()
	e.Renderer = NewTemplate()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	pwd, _ := os.Getwd()
	if len(os.Args) == 2 {
		pwd = os.Args[1]
	}
	//e.GET("/static/*", func(c echo.Context) error {
	//	c.Attachment()
	//})
	e.Static("/static", pwd)
	var fileList []FileItem
	filepath.Walk(pwd, func(filepath string, info os.FileInfo, err error) error {
		if !info.IsDir() && exts[path.Ext(filepath)] {
			item := FileItem{
				Name: info.Name(),
				Path: strings.ReplaceAll(filepath, pwd, ""),
				Size: FormatBytesLength(info.Size()),
			}
			if runtime.GOOS == "windows" {
				item.Path = strings.ReplaceAll(item.Path, "\\", "/")
			}
			fileList = append(fileList, item)
		}
		return nil
	})
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/static/") {
				c.Response().Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", path.Base(c.Request().URL.Path)))
				m, err := mimetype.DetectFile(path.Join(pwd, strings.ReplaceAll(c.Request().URL.Path, "/static", "")))
				if err != nil {
					return err
				}
				fmt.Println(m.String(), m.Extension(), path.Base(c.Request().URL.Path))
				c.Response().Header().Set("content-type", m.String())
			}
			return next(c)
		}
	})

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", fileList)
	})

	e.POST("/upload", func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(fmt.Sprintf("%s/%s", pwd, file.Filename))
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		scheme := c.Scheme()
		host := c.Request().Host
		base := scheme + "://" + host
		return c.Redirect(http.StatusFound, base)
	})

	if runtime.GOOS == "windows" {
		dir := path.Dir(os.Args[0])
		Run(dir, "cmd", "/c", "start", "http://localhost:1614/")
	}
	// Start server
	e.Logger.Fatal(e.Start(":1614"))
}

func FormatBytesLength(length int64) string {
	if length < 1024*1024 {
		return fmt.Sprintf("%d K", length/(1024))
	} else {
		return fmt.Sprintf("%d M", length/(1024*1024))
	}
}

func Run(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
