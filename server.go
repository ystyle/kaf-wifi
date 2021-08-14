package kafwifi

import (
	"embed"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed public
var public embed.FS

var (
	exts = map[string]bool{
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
	port   = 1614
	ipList []string
)

var version string

type FileItem struct {
	Name string
	Path string
	Size string
}

func initIpList() {
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
				ip := fmt.Sprintf("http://%s:%d", ipnet.IP.String(), port)
				fmt.Println(ip)
				ipList = append(ipList, ip)
			}
		}
	}
}

func Start() {
	fmt.Println("kaf-wifi: ", version)
	initIpList()
	// Echo instance
	e := echo.New()
	e.Logger.SetLevel(log.OFF)
	e.Renderer = newTemplate()

	// Middleware
	e.Use(middleware.Recover())

	pwd, _ := os.Getwd()
	if len(os.Args) >= 2 {
		pwd = os.Args[1]
	}
	e.Static("/static", pwd)

	// 添加处理mobi, azw3文件识别处理插件
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/static/") {
				// 转换中文在前端显示
				contentDisposition := fmt.Sprintf(`attachment; filename="%s"`, path.Base(c.Request().URL.Path))
				c.Response().Header().Add("Content-Disposition", url.PathEscape(contentDisposition))
				// 识别文件类型
				m, err := mimetype.DetectFile(path.Join(pwd, strings.ReplaceAll(c.Request().URL.Path, "/static", "")))
				if err != nil {
					return err
				}
				// 添加内容类型 mimetype， 以让kindle浏览器识别mobi, azw3文件
				c.Response().Header().Set("content-type", m.String())
			}
			return next(c)
		}
	})

	// 访问页面
	e.GET("/", func(c echo.Context) error {
		// 扫描目录下面的所有书籍文件
		var fileList []FileItem
		filepath.Walk(pwd, func(filepath string, info os.FileInfo, err error) error {
			if !info.IsDir() && exts[path.Ext(filepath)] {
				item := FileItem{
					Name: info.Name(),
					Path: strings.ReplaceAll(filepath, pwd, ""),
					Size: formatBytesLength(info.Size()),
				}
				if runtime.GOOS == "windows" {
					item.Path = strings.ReplaceAll(item.Path, "\\", "/")
				}
				fileList = append(fileList, item)
			}
			return nil
		})
		return c.Render(http.StatusOK, "index", map[string]interface{}{
			"version":  version,
			"fileList": fileList,
			"ipList":   ipList,
		})
	})
	e.GET("/favicon.png", func(c echo.Context) error {
		bs, err := public.ReadFile("public/favicon.png")
		if err != nil {
			return err
		}
		return c.Blob(http.StatusOK, "image/png", bs)
	})

	// 上传文件，由于kindle不支持上传，前端禁用了本功能
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

	// 在windows时自动打开浏览器
	if runtime.GOOS == "windows" {
		openUI()
	}
	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func openUI() {
	dir := path.Dir(os.Args[0])
	run(dir, "cmd", "/c", "start", fmt.Sprintf("http://localhost:%d/", port))
}
