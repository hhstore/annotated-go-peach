/*
模块说明:
	- 上层引用模块:
		- runWeb(): github.com/peachdocs/peach/cmd/web.go:53
	- 包含:
		- Home(): 首页渲染
		- Pages(): 找不到的页面, 遍历匹配.
	- 依赖模块:
		- github.com/peachdocs/peach/models/toc.go:93
		- github.com/peachdocs/peach/routers/docs.go:32

阅读经验:
	- 顺藤摸瓜, 逐层向下读.
	- 阅读要按照逻辑顺序, 切忌割裂开内在组织, 盲目的一个个文件读.(效率不高)

 */
package routers

import (
	"fmt"
	"strings"

	"github.com/Unknwon/com"

	"github.com/peachdocs/peach/models"
	"github.com/peachdocs/peach/modules/middleware"
	"github.com/peachdocs/peach/modules/setting"
)

// 首页路由
//	- 注意 ctx 类型
func Home(ctx *middleware.Context) {
	if !setting.Page.HasLandingPage {
		ctx.Redirect(setting.Page.DocsBaseURL)
		return
	}

	ctx.HTML(200, "home")
}

/*
	页面生成:
	- 依赖模块:
		- github.com/peachdocs/peach/models/toc.go:93
		- github.com/peachdocs/peach/routers/docs.go:32


 */
func Pages(ctx *middleware.Context) {
	toc := models.Tocs[ctx.Locale.Language()]
	if toc == nil {
		toc = models.Tocs[setting.Docs.Langs[0]]	// 默认文档语言类型
	}

	pageName := strings.ToLower(strings.TrimSuffix(ctx.Req.URL.Path[1:], ".html"))

	// 遍历生成 HTML 页面
	for i := range toc.Pages {
		if toc.Pages[i].Name == pageName {
			page := toc.Pages[i]
			langVer := toc.Lang
			if !com.IsFile(page.FileName) {
				ctx.Data["IsShowingDefault"] = true
				langVer = setting.Docs.Langs[0]
				page = models.Tocs[langVer].Pages[i]
			}
			if !setting.ProdMode {
				page.ReloadContent()	// 生成 HTML 页面 [github.com/peachdocs/peach/models/toc.go:93]
			}

			ctx.Data["Title"] = page.Title
			ctx.Data["Content"] = fmt.Sprintf(`<script type="text/javascript" src="/%s/%s?=%d"></script>`,
				langVer, page.DocumentPath+".js", page.LastBuildTime)
			ctx.Data["Pages"] = toc.Pages

			renderEditPage(ctx, page.DocumentPath)		// 页面渲染 [github.com/peachdocs/peach/routers/docs.go:32]
			ctx.HTML(200, "docs")
			return
		}
	}

	NotFound(ctx)	// 404页面
}

func NotFound(ctx *middleware.Context) {
	ctx.Data["Title"] = "404"
	ctx.HTML(404, "404")
}
