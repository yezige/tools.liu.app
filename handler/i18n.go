package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"

	ginI18n "github.com/gin-contrib/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var AcceptLanguage = []language.Tag{language.Chinese, language.English}

func SetI18n() gin.HandlerFunc {
	return ginI18n.Localize(
		// 默认使用 yaml 格式解析语言文件
		ginI18n.WithBundle(&ginI18n.BundleCfg{
			RootPath:         "./lang",
			AcceptLanguage:   AcceptLanguage,
			DefaultLanguage:  language.Chinese,
			FormatBundleFile: "yaml",
			UnmarshalFunc:    yaml.Unmarshal,
		}),
		// 默认从Accept-Language或lng参数判断语言
		// 修改为优先从lng判断语言
		ginI18n.WithGetLngHandle(
			func(context *gin.Context, defaultLng string) string {
				if context == nil || context.Request == nil {
					return defaultLng
				}

				lng := context.Query("lng")
				if lng != "" {
					return lng
				}

				lng = context.GetHeader("Accept-Language")
				if lng != "" {
					return getAcceptLanguage(lng)
				}

				return defaultLng
			},
		),
	)
}

func I(key string, args ...string) string {
	if len(args) == 0 {
		if message := ginI18n.MustGetMessage(key); message != "" {
			return message
		} else {
			return key
		}
	}

	data := make(map[string]string, len(args))
	for k, v := range args {
		data["arg"+strconv.Itoa(k)] = v
	}
	if message := ginI18n.MustGetMessage(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	}); message != "" {
		return message
	} else {
		return key
	}
}

func getAcceptLanguage(acceptLanguageHeader string) string {
	var matcher = language.NewMatcher(AcceptLanguage)
	t, _, _ := language.ParseAcceptLanguage(acceptLanguageHeader)
	_, idx, _ := matcher.Match(t...)

	return AcceptLanguage[idx].String()
}
