package apidoc

import (
	"html/template"

	"fibera/docs"
	"fibera/internal/config"

	"github.com/gofiber/fiber/v3"
	scalar "github.com/yokeTH/gofiber-scalar/scalar/v3"
	"go.uber.org/zap"
)

// RegisterRoutes 注册 API 文档路由，接口定义由 swag 注解生成，页面由 Scalar 渲染
func RegisterRoutes(app *fiber.App, conf *config.Config) {
	if !conf.Docs.Enable {
		return
	}

	docsConf := conf.Docs.Normalize()

	// 标题与版本以运行时配置为准，配置缺失时保留注解中的缺省值
	docs.SwaggerInfo.Title = docsConf.Title
	if conf.App.Version != "" {
		docs.SwaggerInfo.Version = conf.App.Version
	}

	handler := scalar.New(scalar.Config{
		BasePath:          docsConf.BasePath,
		Path:              docsConf.Path,
		Title:             docsConf.Title,
		FileContentString: docs.SwaggerInfo.ReadDoc(),
		Theme:             scalar.Theme(docsConf.Theme),
		CustomStyle:       template.CSS(docsConf.CustomStyle),
		ProxyUrl:          docsConf.ProxyURL,
		CacheAge:          docsConf.CacheAge,
		FallbackCacheAge:  docsConf.FallbackCacheAge,
		ForceOffline:      forceOffline(docsConf.ForceOffline),
	})

	// UI 根路径与子路径分别注册，子路径用于内置脚本与 doc.json
	app.Get(docsConf.Path, handler)
	app.Get(docsConf.Path+"/*", handler)

	zap.L().Info("api documentation enabled", zap.String("path", docsConf.Path))
}

// forceOffline 将布尔配置转换为 Scalar 要求的指针类型
func forceOffline(offline bool) *bool {
	if offline {
		return scalar.ForceOfflineTrue
	}
	return scalar.ForceOfflineFalse
}
