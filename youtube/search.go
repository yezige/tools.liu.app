package youtube

import (
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unsafe"

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
	Video  *Video             `json:"video"`
	Info   *youtubedl.Video   `json:"info"`
	Format *[]SelectionFormat `json:"format"`
}

type Video struct {
	ID         string            `json:"id"`
	Snippet    SectionSnippet    `json:"snippet"`
	Statistics SectionStatistics `json:"statistics"`
}
type SearchVideo struct {
	ID         SelectionID       `json:"id"`
	Snippet    SectionSnippet    `json:"snippet"`
	Statistics SectionStatistics `json:"statistics"`
}
type Videos []*Video
type SearchVideos []*SearchVideo

type Videoer interface {
	SetThumbnail()
}

func (vs Videos) SetThumbnail() {
	for _, v := range vs {
		SetVideoThumbnail(&v.Snippet.Thumbnails.Default)
		SetVideoThumbnail(&v.Snippet.Thumbnails.High)
		SetVideoThumbnail(&v.Snippet.Thumbnails.Medium)
	}
}
func (vs SearchVideos) SetThumbnail() {
	for _, v := range vs {
		SetVideoThumbnail(&v.Snippet.Thumbnails.Default)
		SetVideoThumbnail(&v.Snippet.Thumbnails.High)
		SetVideoThumbnail(&v.Snippet.Thumbnails.Medium)
	}
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

type SetLink struct {
	Key   string `json:"key"`
	Value interface{} `json:"value"`
}

func Popular(regionCode string, videoCategoryId string) (result *PopularResult, err error) {
	// 先查询redis
	popular, err := redis.New().Get("youtube:popular:" + regionCode + ":" + videoCategoryId)
	if err == nil {
		if err := json.Unmarshal([]byte(popular), &result); err != nil {
			return nil, err
		}
		// 修改 i.ytimg.com 为 ytimg.liu.app
		SetVideoURL((*Videos)(unsafe.Pointer(&result.Items)))
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
	if err := redis.New().SetTTL("youtube:popular:"+regionCode+":"+videoCategoryId, res, time.Hour*24); err != nil {
		return nil, err
	}

	// 接口调用结果解析到结构体
	if err := json.Unmarshal([]byte(res), &result); err != nil {
		return nil, err
	}

	// 修改 i.ytimg.com 为 ytimg.liu.app
	SetVideoURL((*Videos)(unsafe.Pointer(&result.Items)))

	return result, nil
}

func Search(key string) (result *SearchResult, err error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}
	keyRedis := strings.ReplaceAll(key, " ", "_")

	// 先查询redis
	search, err := redis.New().Get("youtube:search:" + keyRedis)
	if err == nil {
		if err := json.Unmarshal([]byte(search), &result); err != nil {
			return nil, err
		}
		// 修改 i.ytimg.com 为 ytimg.liu.app
		SetVideoURL((*SearchVideos)(unsafe.Pointer(&result.Items)))
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
	if err := redis.New().SetTTL("youtube:search:"+keyRedis, res, time.Hour*1); err != nil {
		return nil, err
	}

	// 接口调用结果解析到结构体
	if err := json.Unmarshal([]byte(res), &result); err != nil {
		return nil, err
	}

	// 修改 i.ytimg.com 为 ytimg.liu.app
	SetVideoURL((*SearchVideos)(unsafe.Pointer(&result.Items)))

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
		SetVideoURL((*Videos)(unsafe.Pointer(&result.Items)))

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
	SetVideoURL((*Videos)(unsafe.Pointer(&result.Items)))

	return result, nil
}

func Download(id string, nocache bool) (result *DownloadResult, err error) {

	if id == "" {
		return nil, errors.New("id is empty")
	}

	// 先查询redis
	if !nocache {
		info, err := redis.New().Get("youtube:download:" + id)
		if err == nil {
			if err := json.Unmarshal([]byte(info), &result); err != nil {
				return nil, err
			}
			return result, err
		}
	}

	// 调用youtube-dl接口获取视频下载地址
	client := youtubedl.Client{Debug: false}

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
		v.URL_PROXY = SetDownloadUrl(v.URL, fileName)
		resSlice[k] = SelectionFormat{F: v, Ext: GetExtByMime(v.MimeType)}
	}

	video.Formats = nil

	result = &DownloadResult{
		Video: &Video{
			ID: video.ID,
			Snippet: SectionSnippet{
				Title:       video.Title,
				Description: video.Description,
				PublishedAt: video.PublishDate.Format("2023-04-12 14:38:01"),
				Thumbnails: SectionThumbnail{
					High: SectionThumbnailDetail{
						Url: video.Thumbnails[len(video.Thumbnails)-1].URL,
					},
				},
			},
		},
		Info:   video,
		Format: &resSlice,
	}

	// 修改 i.ytimg.com 为 ytimg.liu.app
	SetVideoURL(Videos{result.Video})

	// 存储到redis
	formatJson, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	if err := redis.New().SetTTL("youtube:download:"+id, formatJson, time.Minute*500); err != nil {
		return nil, err
	}

	return result, nil
}

func DownloadSetLink(param SetLink) (err error) {

	if param.Key == "" {
		return errors.New("key is empty")
	}

	// 存储到redis
	formatJson, err := json.Marshal(param.Value)
	if err != nil {
		return err
	}
	if err := redis.New().SetTTL("youtube:downloadlink:"+param.Key, formatJson, time.Hour*24); err != nil {
		return err
	}

	return nil
}

func DownloadGetLink(key string) (result map[string]interface{}, err error) {

	if key == "" {
		return nil, errors.New("key is empty")
	}

	info, err := redis.New().Get("youtube:downloadlink:" + key)
	if err == nil {
		if err := json.Unmarshal([]byte(info), &result); err != nil {
			return nil, err
		}
		return result, err
	}

	return nil, errors.New("link does not exist")
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
func SetVideoURL(videos Videoer) {
	// 修改 i.ytimg.com 为 ytimg.liu.app
	videos.SetThumbnail()
}
func SetVideoThumbnail(std *SectionThumbnailDetail) {
	std.Url = strings.Replace(std.Url, "i.ytimg.com", "ytimg.liu.app", 1)
}
func SetDownloadUrl(u string, fileName string) string {
	// 修改为ytdl.liu.app
	params := url.Values{}
	params.Add("u", u)
	params.Add("filename", fileName)
	return `https://ytdl.liu.app?` + params.Encode()
}
