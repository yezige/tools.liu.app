package youtube

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/yezige/tools.liu.app/config"
	"github.com/yezige/tools.liu.app/logx"
	"github.com/yezige/tools.liu.app/redis"
	"github.com/yezige/tools.liu.app/request"

	youtubedl "github.com/yezige/youtube/v2"
)

const (
	YoutubePopular = "https://youtube.googleapis.com/youtube/v3/videos"
	YoutubeSearch  = "https://youtube.googleapis.com/youtube/v3/search"
)

var (
	videoRegexpList = []*regexp.Regexp{
		regexp.MustCompile(`(?:v|embed|shorts|watch\?v)(?:=|/)([^"&?/=%]{11})`),
		regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
		regexp.MustCompile(`([^"&?/=%]{11})`),
	}
)

type PopularResult struct {
	Items         []*Video `json:"items"`
	Etag          string   `json:"etag"`
	Kind          string   `json:"kind"`
	NextPageToken string   `json:"nextPageToken"`
}

type DownloadResult struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Author      string               `json:"author"`
	Format      *[]SelectionFormat   `json:"format"`
	Playability SelectionPlayability `json:"playability"`
}

type Video struct {
	ID         string            `json:"id"`
	Snippet    SectionSnippet    `json:"snippet"`
	Statistics SectionStatistics `json:"statistics"`
}
type SectionSnippet struct {
	CategoryID           string           `json:"categoryId"`
	Description          string           `json:"description"`
	PublishedAt          string           `json:"publishedAt"`
	Title                string           `json:"title"`
	ChannelId            string           `json:"channelId"`
	ChannelTitle         string           `json:"channelTitle"`
	DefaultAudioLanguage string           `json:"defaultAudioLanguage"`
	Tags                 []string         `json:"tags"`
	Thumbnails           SectionThumbnail `json:"thumbnails"`
}

type SectionThumbnail struct {
	Default SectionThumbnailDetail `json:"default"`
	Medium  SectionThumbnailDetail `json:"medium"`
	High    SectionThumbnailDetail `json:"high"`
}

type SectionThumbnailDetail struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
}

type SectionStatistics struct {
	ViewCount     string `json:"viewCount"`
	LikeCount     string `json:"likeCount"`
	FavoriteCount string `json:"favoriteCount"`
	CommentCount  string `json:"commentCount"`
}

type SearchResult struct {
	Items         []*SearchVideo `json:"items"`
	Etag          string         `json:"etag"`
	Kind          string         `json:"kind"`
	NextPageToken string         `json:"nextPageToken"`
}

type SearchVideo struct {
	ID         SelectionID       `json:"id"`
	Snippet    SectionSnippet    `json:"snippet"`
	Statistics SectionStatistics `json:"statistics"`
}
type SelectionID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

type SelectionFormat struct {
	F   youtubedl.Format `json:"f"`
	Ext string           `json:"ext"`
}

type SelectionPlayability struct {
	Status string `json:"status"`
}

func Popular(regionCode string, videoCategoryId string) (result *PopularResult, err error) {
	// 先查询redis
	popular, err := redis.New().Get("youtube:popular:" + regionCode + ":" + videoCategoryId)
	if err == nil {
		if err := json.Unmarshal([]byte(popular), &result); err != nil {
			return nil, err
		}
		// 修改 i.ytimg.com 为 ytimg.liu.app
		SetVideoURL(result.Items)
		return result, nil
	}

	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	// 调接口查询当前热门视频信息
	res, err := request.New(YoutubePopular).SetParams(map[string]interface{}{
		"part":            "snippet,contentDetails,statistics",
		"chart":           "mostPopular",
		"regionCode":      regionCode,
		"videoCategoryId": videoCategoryId,
		"key":             conf.API.GoogleAPIKey,
		"maxResults":      "10",
	}).Get().GetBody()

	if err != nil {
		return nil, err
	}

	// 存储到redis
	if err := redis.New().SetTTL("youtube:popular:"+regionCode+":"+videoCategoryId, res, time.Hour*1); err != nil {
		return nil, err
	}

	// 接口调用结果解析到结构体
	if err := json.Unmarshal([]byte(res), &result); err != nil {
		return nil, err
	}

	// 修改 i.ytimg.com 为 ytimg.liu.app
	SetVideoURL(result.Items)

	return result, nil
}

func Search(key string) (result *SearchResult, err error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	// 先查询redis
	search, err := redis.New().Get("youtube:search:" + key)
	if err == nil {
		if err := json.Unmarshal([]byte(search), &result); err != nil {
			return nil, err
		}
		// 修改 i.ytimg.com 为 ytimg.liu.app
		SetSearchVideo(result.Items)
		return result, err
	}

	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	// 调接口搜索视频信息
	res, err := request.New(YoutubeSearch).SetParams(map[string]interface{}{
		"part":       "snippet",
		"q":          key,
		"key":        conf.API.GoogleAPIKey,
		"maxResults": "10",
	}).Get().GetBody()

	if err != nil {
		return nil, err
	}

	// 存储到redis
	if err := redis.New().SetTTL("youtube:search:"+key, res, time.Hour*1); err != nil {
		return nil, err
	}

	// 接口调用结果解析到结构体
	if err := json.Unmarshal([]byte(res), &result); err != nil {
		return nil, err
	}

	// 修改 i.ytimg.com 为 ytimg.liu.app
	SetSearchVideo(result.Items)

	return result, nil
}

func GetInfo(id string) (result *PopularResult, err error) {

	// 先查询redis
	info, err := redis.New().Get("youtube:info:" + id)
	if err == nil {
		if err := json.Unmarshal([]byte(info), &result); err != nil {
			return nil, err
		}
		// 修改 i.ytimg.com 为 ytimg.liu.app
		SetVideoURL(result.Items)

		return result, err
	}

	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	// 调接口查询视频信息
	res, err := request.New(YoutubePopular).SetParams(map[string]interface{}{
		"part": "snippet,contentDetails,statistics",
		"id":   id,
		"key":  conf.API.GoogleAPIKey,
	}).Get().GetBody()

	if err != nil {
		return nil, err
	}

	// 存储到redis
	if err := redis.New().SetTTL("youtube:info:"+id, res, time.Hour*24); err != nil {
		return nil, err
	}

	// 接口调用结果解析到结构体
	if err := json.Unmarshal([]byte(res), &result); err != nil {
		return nil, err
	}

	// 修改 i.ytimg.com 为 ytimg.liu.app
	SetVideoURL(result.Items)

	return result, nil
}

func Download(id string) (result *DownloadResult, err error) {

	// 先查询redis
	info, err := redis.New().Get("youtube:download:" + id)
	if err == nil {
		if err := json.Unmarshal([]byte(info), &result); err != nil {
			return nil, err
		}
		return result, err
	}

	// 调用youtube-dl接口获取视频下载地址
	client := youtubedl.Client{Debug: true}

	video, err := client.GetVideo(id)
	if err != nil {
		logx.LogError.Infoln(err)
		return nil, err
	}

	// 排序
	video.Formats.Sort()

	// 计算后缀
	resSlice := make([]SelectionFormat, len(video.Formats))
	for k, v := range video.Formats {
		v.URL, err = client.GetStreamURL(video, &v)
		if err != nil {
			logx.LogError.Infoln(err)
		}
		fileName := url.QueryEscape(video.Title) + "." + GetExtByMime(v.MimeType)
		v.URL = SetDownloadUrl(v.URL, fileName)
		resSlice[k] = SelectionFormat{F: v, Ext: GetExtByMime(v.MimeType)}
	}

	result = &DownloadResult{
		Title:       video.Title,
		Description: video.Description,
		Author:      video.Author,
		Format:      &resSlice,
		Playability: SelectionPlayability{
			Status: video.PlayabilityStatus.Status,
		},
	}

	// 存储到redis
	formatJson, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	if err := redis.New().SetTTL("youtube:download:"+id, formatJson, time.Minute*10); err != nil {
		return nil, err
	}

	return result, nil
}

func GetExtByMime(MimeType string) string {
	mime := map[string]string{
		"video/mp4":  "mp4",
		"video/webm": "webm",
		"video/ogg":  "ogv",
	}
	for k, v := range mime {
		if strings.Contains(MimeType, k) {
			return v
		}
	}
	return "mp4"
}

func ExtractVideoID(videoID string) (string, error) {
	if strings.Contains(videoID, "youtu") || strings.ContainsAny(videoID, "\"?&/<%=") {
		for _, re := range videoRegexpList {
			if isMatch := re.MatchString(videoID); isMatch {
				subs := re.FindStringSubmatch(videoID)
				videoID = subs[1]
			}
		}
		return videoID, nil
	}
	return "", errors.New("invalid video id")
}
func SetVideoURL(vs []*Video) {
	// 修改 i.ytimg.com 为 ytimg.liu.dev
	for _, item := range vs {
		item.Snippet.Thumbnails.Default.Url = strings.Replace(item.Snippet.Thumbnails.Default.Url, "i.ytimg.com", "ytimg.liu.dev", 1)
		item.Snippet.Thumbnails.High.Url = strings.Replace(item.Snippet.Thumbnails.High.Url, "i.ytimg.com", "ytimg.liu.dev", 1)
		item.Snippet.Thumbnails.Medium.Url = strings.Replace(item.Snippet.Thumbnails.Medium.Url, "i.ytimg.com", "ytimg.liu.dev", 1)
	}
}
func SetSearchVideo(vs []*SearchVideo) {
	// 修改 i.ytimg.com 为 ytimg.liu.dev
	for _, item := range vs {
		item.Snippet.Thumbnails.Default.Url = strings.Replace(item.Snippet.Thumbnails.Default.Url, "i.ytimg.com", "ytimg.liu.dev", 1)
		item.Snippet.Thumbnails.High.Url = strings.Replace(item.Snippet.Thumbnails.High.Url, "i.ytimg.com", "ytimg.liu.dev", 1)
		item.Snippet.Thumbnails.Medium.Url = strings.Replace(item.Snippet.Thumbnails.Medium.Url, "i.ytimg.com", "ytimg.liu.dev", 1)
	}
}
func SetDownloadUrl(u string, fileName string) string {
	// 修改为ytdl.liu.dev
	params := url.Values{}
	params.Add("u", u)
	params.Add("filename", fileName)
	return `https://ytdl.liu.app?` + params.Encode()
}
