package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yezige/tools.liu.app/config"
	"github.com/yezige/tools.liu.app/core"
	"github.com/yezige/tools.liu.app/redis"
	"github.com/yezige/tools.liu.app/youtube"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func YoutubeHandler(c *gin.Context) {
	// 更新总下载量
	redis.New().IncrBy("youtube:download:total", 1)

	// 默认页面配置
	data := core.PageIndexConfig{
		Page: core.PageConfig{
			Site:        &config.SectionSite{},
			Title:       I("title") + " | Home",
			Description: I("desc"),
			Keywords:    I("keywords"),
			Tags:        strings.Split(I("tags"), "|"),
			Permalink:   "/youtube",
			BodyName:    I("body-name-youtube"),
		},
		DownloadTotal: redis.New().GetInt64("youtube:download:total"),
	}

	cfg, err := config.GetConfig()
	if err != nil {
		data.Page.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "pages/error/error.tmpl", data)
		return
	}
	data.Page.Site = &cfg.Site

	// 获取当前美国流行视频列表
	popularUS, err := youtube.Popular("US", "0")
	if err != nil {
		data.Page.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "pages/error/error.tmpl", data)
		return
	}
	data.PopularList = append(data.PopularList, core.SelectionPopular{
		Title: I("popular-us"),
		Items: popularUS.Items,
	})

	// // 获取当前香港流行视频列表
	// popularHK, err := youtube.Popular("HK", "0")
	// if err != nil {
	// 	data.ErrorMsg = err.Error()
	// 	c.HTML(http.StatusOK, "pages/error/error.tmpl", data)
	// 	return
	// }
	// data.PopularList = append(data.PopularList, core.SelectionPopular{
	// 	Title: I("popular-hk"),
	// 	Items: popularHK.Items,
	// })

	// 定义多个热门分类列表
	var popularCategoryDict map[string]string = map[string]string{
		"1":  I("popular-film-animation"),
		"20": I("popular-gaming"),
		// "24": I("popular-entertainment"),
		"10": I("popular-music"),
		// "17": I("popular-sports"),
	}

	for k, v := range popularCategoryDict {
		// 获取当前中国流行视频列表
		popular, err := youtube.Popular("US", k)
		if err != nil {
			data.Page.ErrorMsg = err.Error()
			c.HTML(http.StatusOK, "pages/error/error.tmpl", data)
			return
		}
		data.PopularList = append(data.PopularList, core.SelectionPopular{
			Title: v,
			Items: popular.Items,
		})
	}

	// 定义 AdsConfig
	// 首页 第五个video
	data.Page.AdsConfig = append(data.Page.AdsConfig, core.AdsConfig{
		Format:    "fluid",
		LayoutKey: "-7c+dd+2b+o-1d",
		Slot:      "1153918191",
		Client:    cfg.Site.GoogleAdsense,
		Class:     "video ads",
	})
	// 首页 第二个H2
	data.Page.AdsConfig = append(data.Page.AdsConfig, core.AdsConfig{
		Slot:   "9245845411",
		Client: cfg.Site.GoogleAdsense,
		Class:  "popular_videos_box",
	})

	c.HTML(http.StatusOK, "pages/youtube/index.tmpl", data)
}

func YoutubeDownloadHandler(c *gin.Context) {
	// 更新总下载量
	redis.New().IncrBy("youtube:download:total", 1)

	// 默认页面配置
	data := core.PageDownloadConfig{
		Page: core.PageConfig{
			Site:        &config.SectionSite{},
			Title:       I("title") + " | Download",
			Description: I("desc"),
			Keywords:    I("keywords"),
			Tags:        strings.Split(I("tags"), "|"),
			Permalink:   "/youtube/download",
			BodyName:    I("body-name-youtube"),
		},
		DownloadTotal: redis.New().GetInt64("youtube:download:total"),
	}

	cfg, err := config.GetConfig()
	if err != nil {
		data.Page.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "pages/error/error.tmpl", data)
		return
	}
	data.Page.Site = &cfg.Site

	// 视频id不能为空
	videoID := c.Query("id")
	if videoID == "" {
		data.Page.ErrorMsg = "Video id is empty"
		c.HTML(http.StatusOK, "pages/error/error.tmpl", data)
		return
	}
	// 是否使用缓存
	noCache, _ := strconv.ParseBool(c.DefaultQuery("nocache", "0"))

	// 改为从Download获取
	// // 通过视频id获取视频信息
	// info, err := youtube.GetInfo(videoID)
	// if err != nil {
	// 	data.Page.ErrorMsg = err.Error()
	// 	c.HTML(http.StatusOK, "pages/error/error.tmpl", data)
	// 	return
	// }
	// data.Info = *info.Items[0]

	// // 根据视频信息更新Title
	// data.Page.Title = data.Info.Snippet.Title + " | " + data.Page.Title
	// data.Page.Description = " | " + data.Info.Snippet.Description + " | " + data.Page.Description
	// data.Page.Keywords = " | " + data.Info.Snippet.Description + " | " + data.Page.Keywords

	// 通过视频id获取视频下载信息
	down, err := youtube.Download(videoID, noCache)
	if err != nil {
		data.Page.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "pages/error/error.tmpl", data)
		return
	}

	data.DownloadList = down.Format
	data.PlayabilityStatus = down.Info.PlayabilityStatus.Status
	data.Info = down.Video

	// 根据视频信息更新Title
	data.Page.Title = data.Info.Snippet.Title + " | " + data.Page.Title
	data.Page.Description = " | " + data.Info.Snippet.Description + " | " + data.Page.Description
	data.Page.Keywords = " | " + data.Info.Snippet.Description + " | " + data.Page.Keywords
	// 下载-尾
	data.Page.AdsConfig = append(data.Page.AdsConfig, core.AdsConfig{
		Slot:   "2798667502",
		Client: cfg.Site.GoogleAdsense,
		Class:  "adsbygoogle",
	})

	c.HTML(http.StatusOK, "pages/youtube/download.tmpl", data)
}

func YoutubeSearchHandler(c *gin.Context) {
	// 默认页面配置
	data := core.PageSearchConfig{
		Page: core.PageConfig{
			Site:        &config.SectionSite{},
			Title:       I("title") + " | Search",
			Description: I("desc"),
			Keywords:    I("keywords"),
			Tags:        strings.Split(I("tags"), "|"),
			Permalink:   "/youtube/search",
			BodyName:    I("body-name-youtube"),
		},
	}

	cfg, err := config.GetConfig()
	if err != nil {
		data.Page.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "pages/error/error.tmpl", data)
		return
	}
	data.Page.Site = &cfg.Site

	// 视频搜索关键字不能为空
	q := c.Query("q")
	if q == "" {
		c.Redirect(http.StatusMovedPermanently, "/youtube")
		return
	}
	// 视频链接直接跳下载页面
	id, err := youtube.ExtractVideoID(q)
	if err == nil {
		c.Redirect(http.StatusMovedPermanently, "/youtube/download?id="+id)
		return
	}

	data.Q = q
	data.Page.Title += " | " + q

	// 搜索视频
	search, err := youtube.Search(q)
	if err != nil {
		data.Page.ErrorMsg = err.Error()
		c.HTML(http.StatusOK, "pages/error/error.tmpl", data)
		return
	}
	data.SearchList = search.Items
	// 搜索-头
	data.Page.AdsConfig = append(data.Page.AdsConfig, core.AdsConfig{
		Slot:   "1820627751",
		Client: cfg.Site.GoogleAdsense,
		Class:  "adsbygoogle",
	})
	// 搜索-尾
	data.Page.AdsConfig = append(data.Page.AdsConfig, core.AdsConfig{
		Slot:   "8194464410",
		Client: cfg.Site.GoogleAdsense,
		Class:  "adsbygoogle",
	})

	c.HTML(http.StatusOK, "pages/youtube/search.tmpl", data)
}

// 获取当前热门视频接口
func ApiYoutubePopularHandler(c *gin.Context) {
	result, err := youtube.Popular("US", "0")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "成功",
		"data": result,
	})
}

// 查询视频接口
func ApiYoutubeSearchHandler(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	result, err := youtube.Search(key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "成功",
		"data": result,
	})
}

// 查询视频下载信息接口
func ApiYoutubeDownloadHandler(c *gin.Context) {
	videoID := c.DefaultQuery("id", "")
	noCache, _ := strconv.ParseBool(c.DefaultQuery("nocache", "0"))
	result, err := youtube.Download(videoID, noCache)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "成功",
		"data": result,
	})
}

// 设置视频下载链接接口
func ApiYoutubeDownloadSetLinkHandler(c *gin.Context) {
	key := c.DefaultPostForm("key", "")
	value := c.DefaultPostForm("value", "")
	err := youtube.DownloadSetLink(key, value)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "成功",
	})
}

// 获取视频下载链接接口
func ApiYoutubeDownloadGetLinkHandler(c *gin.Context) {
	key := c.DefaultQuery("key", "")
	result, err := youtube.DownloadGetLink(key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "成功",
		"data": result,
	})
}

// 解释视频大小
func SizeFormat(size int64) string {
	var unit string
	if size < 1024 {
		unit = "B"
	} else if size < 1024*1024 {
		unit = "KB"
		size /= 1024
	} else if size < 1024*1024*1024 {
		unit = "MB"
		size /= 1024 * 1024
	} else {
		unit = "GB"
		size /= 1024 * 1024 * 1024
	}
	return strconv.FormatInt(size, 10) + unit
}

// 解释视频码率
func BpsFormat(bps int) string {
	var unit string
	if bps < 1024 {
		unit = "bps"
	} else if bps < 1024*1024 {
		unit = "Kbps"
		bps /= 1024
	} else if bps < 1024*1024*1024 {
		unit = "Mbps"
		bps /= 1024 * 1024
	} else {
		unit = "Gbps"
		bps /= 1024 * 1024 * 1024
	}
	return strconv.Itoa(bps) + unit
}

// 换算为10,000,000的格式
func NumFormat(num int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", num)
}

// 模板默认值
func Default(key string, args ...string) string {
	if key != "" {
		return key
	}
	if len(args) == 1 {
		return args[0]
	}
	return ""
}
