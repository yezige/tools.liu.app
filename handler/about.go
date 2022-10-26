package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yezige/tools.liu.app/config"
	"github.com/yezige/tools.liu.app/core"
)

func AboutHandler(c *gin.Context) {
	conf := core.PageConfig{
		Site:        &config.SectionSite{},
		Title:       I("title") + " | About",
		Description: I("desc"),
		Keywords:    I("keywords"),
		Tags:        strings.Split(I("tags"), "|"),
		Permalink:   "/about",
		BodyName:    I("body-name-about"),
	}
	cfg, err := config.GetConfig()
	if err != nil {
		conf.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "pages/error/error.tmpl", conf)
		return
	}
	conf.Site = &cfg.Site

	c.HTML(http.StatusOK, "pages/about/index.tmpl", conf)
}
