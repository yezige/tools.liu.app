package core

import (
	"github.com/yezige/tools.liu.app/config"
	"github.com/yezige/tools.liu.app/youtube"
)

var (
	Version string
)

type PageConfig struct {
	Site        *config.SectionSite `json:"site"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Keywords    string              `json:"keywords"`
	Tags        []string            `json:"tags"`
	Permalink   string              `json:"permalink"`
	BodyName    string              `json:"body_name"`
	ErrorMsg    string              `json:"error_msg"`
	AdsConfig   []AdsConfig         `json:"adsconfig"`
	Version     string              `json:"version"`
}
type AdsConfig struct {
	Format          string `default:"auto"`
	LayoutKey       string `default:""`
	Slot            string
	Class           string `default:""`
	Style           string `default:""`
	WidthResponsive string `default:"true"`
	Client          string
}
type PageIndexConfig struct {
	Page          PageConfig         `json:"page"`
	PopularList   []SelectionPopular `json:"popular_list"`
	DownloadTotal int64              `json:"download_total"`
}

type SelectionPopular struct {
	Title string           `json:"popular_title"`
	Items []*youtube.Video `json:"popular_item"`
}

type PageSearchConfig struct {
	Page       PageConfig             `json:"page"`
	SearchList []*youtube.SearchVideo `json:"search_list"`
	Q          string                 `json:"q"`
}

type PageDownloadConfig struct {
	Page              PageConfig                 `json:"page"`
	DownloadList      *[]youtube.SelectionFormat `json:"download_list"`
	Info              *youtube.Video             `json:"info"`
	DownloadTotal     int64                      `json:"download_total"`
}
type PageAboutConfig struct {
	Page PageConfig `json:"page"`
}
type PageTagConfig struct {
	Page PageConfig `json:"page"`
	Name string     `json:"name"`
}
type PageSponsorConfig struct {
	Page PageConfig `json:"page"`
	Name string     `json:"name"`
}

func GetPage(pc PageConfig) PageConfig {
	pc.Version = Version
	return pc
}
