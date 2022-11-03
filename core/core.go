package core

import (
	"github.com/yezige/tools.liu.app/config"
	"github.com/yezige/tools.liu.app/youtube"
)

type PageConfig struct {
	Site          *config.SectionSite `json:"site"`
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	Keywords      string              `json:"keywords"`
	Tags          []string            `json:"tags"`
	Permalink     string              `json:"permalink"`
	BodyName      string              `json:"body_name"`
	PopularList   []SelectionPopular  `json:"popular_list"`
	ErrorMsg      string              `json:"error_msg"`
	DownloadTotal int64               `json:"download_total"`
	AdsConfig     []AdsConfig         `json:"adsconfig"`
}

type SelectionPopular struct {
	Title string           `json:"popular_title"`
	Items []*youtube.Video `json:"popular_item"`
}

type PageSearchConfig struct {
	Site        *config.SectionSite    `json:"site"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Keywords    string                 `json:"keywords"`
	Tags        []string               `json:"tags"`
	Permalink   string                 `json:"permalink"`
	BodyName    string                 `json:"body_name"`
	SearchList  []*youtube.SearchVideo `json:"search_list"`
	ErrorMsg    string                 `json:"error_msg"`
	Q           string                 `json:"q"`
}

type PageDownloadConfig struct {
	Site         *config.SectionSite        `json:"site"`
	Title        string                     `json:"title"`
	Description  string                     `json:"description"`
	Keywords     string                     `json:"keywords"`
	Tags         []string                   `json:"tags"`
	Permalink    string                     `json:"permalink"`
	BodyName     string                     `json:"body_name"`
	DownloadList *[]youtube.SelectionFormat `json:"download_list"`
	ErrorMsg     string                     `json:"error_msg"`
	Info         youtube.Video              `json:"info"`
}
type AdsConfig struct {
	Format          string `default:"auto"`
	LayoutKey       string `default:""`
	Slot            string
	Margin          string `default:"0"`
	WidthResponsive string `default:"true"`
	Client          string
}
