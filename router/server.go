package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	ginapi "github.com/appleboy/gin-status-api"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/golang-queue/queue"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yezige/tools.liu.app/config"
	"github.com/yezige/tools.liu.app/handler"
	"github.com/yezige/tools.liu.app/logx"
)

func rootHandler(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/youtube")
}

func routerEngine(cfg *config.ConfYaml, q *queue.Queue) *gin.Engine {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	isTerm := isatty.IsTerminal(os.Stdout.Fd())
	if isTerm {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stdout,
				NoColor: false,
			},
		)
	}

	// set server mode
	gin.SetMode(cfg.Core.Mode)

	r := gin.Default()

	r.HTMLRender = loadTemplates("templates", "pages")

	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output:    logx.LogAccess.Out,
		Formatter: ginLogFormatter,
	}))

	r.Use(gin.Recovery())
	// 设置语言
	r.Use(handler.SetI18n())

	r.Static("/static", "static")
	r.StaticFile("/sitemap.xml", "sitemap.xml")

	r.GET("/", rootHandler)
	r.GET("/youtube", handler.YoutubeHandler)
	r.GET("/youtube/download", handler.YoutubeDownloadHandler)
	r.GET("/youtube/search", handler.YoutubeSearchHandler)
	r.GET("/about", handler.AboutHandler)
	r.GET("/tags", handler.TagsHandler)
	r.GET("/sponsor", handler.SponsorHandler)

	r.GET(cfg.API.StatGoURI, ginapi.GinHandler)
	r.GET("/api/youtube/popular", handler.ApiYoutubePopularHandler)
	r.GET("/api/youtube/search", handler.ApiYoutubeSearchHandler)
	r.GET("/api/youtube/download", handler.ApiYoutubeDownloadHandler)
	r.GET("/api/youtube/download/setlink", handler.ApiYoutubeDownloadSetLinkHandler)
	r.GET("/api/youtube/download/getlink", handler.ApiYoutubeDownloadGetLinkHandler)

	return r
}

func loadTemplates(templatesDir string, pagesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/**/*.tmpl")
	if err != nil {
		logx.LogError.Infoln(err.Error())
	}

	pages, err := filepath.Glob(pagesDir + "/**/*.tmpl")
	if err != nil {
		logx.LogError.Infoln(err.Error())
	}

	logx.LogAccess.Infoln("检测到页面：")
	// Generate our templates map from our templates/ and pages/ directories
	for _, page := range pages {
		var pageSlice = []string{page}

		files := append(pageSlice, layouts...)

		logx.LogAccess.Infoln(filepath.ToSlash(filepath.Dir(page)) + "/" + filepath.Base(page))

		r.AddFromFilesFuncs(
			filepath.ToSlash(filepath.Dir(page))+"/"+filepath.Base(page),
			template.FuncMap{
				"jsonEncode": func(obj interface{}) string {
					j, _ := json.Marshal(obj)
					return string(j)
				},
				"sizeFormat": func(size int64) string {
					return handler.SizeFormat(size)
				},
				"bpsFormat": func(bps int) string {
					return handler.BpsFormat(bps)
				},
				"numFormat": func(num int64) string {
					return handler.NumFormat(num)
				},
				"I": func(key string, args ...interface{}) string {
					return handler.I(key, args...)
				},
				"default": func(key string, args ...string) string {
					return handler.Default(key, args...)
				},
			},
			files...,
		)
	}

	return r
}

// 自定义日志格式
var ginLogFormatter = func(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}
	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v %s \"%s\"\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.Request.Proto,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}
