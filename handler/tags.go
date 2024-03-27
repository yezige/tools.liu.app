package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yezige/tools.liu.app/config"
	"github.com/yezige/tools.liu.app/core"
)

func TagsHandler(c *gin.Context) {
	conf := core.PageTagConfig{
		Page: core.GetPage(core.PageConfig{
			Site:        &config.SectionSite{},
			Title:       I("title") + " | Tags",
			Description: I("desc"),
			Keywords:    I("keywords"),
			Tags:        strings.Split(I("tags"), "|"),
			Permalink:   "/tags",
			BodyName:    I("body-name-tags"),
		}),
		Name: c.Query("name"),
	}
	cfg, err := config.GetConfig()
	if err != nil {
		conf.Page.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "pages/error/error.tmpl", conf)
		return
	}
	conf.Page.Site = &cfg.Site

	c.HTML(http.StatusOK, "pages/tags/index.tmpl", conf)
}
