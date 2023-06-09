package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yezige/tools.liu.app/config"
	"github.com/yezige/tools.liu.app/core"
)

func SponsorHandler(c *gin.Context) {
	conf := core.PageSponsorConfig{
		Page: core.PageConfig{
			Site:        &config.SectionSite{},
			Title:       I("title") + " | Sponsor",
			Description: I("desc"),
			Keywords:    I("keywords"),
			Tags:        strings.Split(I("sponsors"), "|"),
			Permalink:   "/sponsor",
			BodyName:    I("body-name-sponsor"),
		},
		Name: c.Query("name"),
	}
	cfg, err := config.GetConfig()
	if err != nil {
		conf.Page.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "pages/error/error.tmpl", conf)
		return
	}
	conf.Page.Site = &cfg.Site

	c.HTML(http.StatusOK, "pages/sponsor/index.tmpl", conf)
}
