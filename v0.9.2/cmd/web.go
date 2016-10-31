// Copyright 2015 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cmd

import (
	"fmt"
	"net/http"

	"github.com/Unknwon/log"
	"github.com/codegangsta/cli"
	"github.com/go-macaron/i18n"
	"github.com/go-macaron/pongo2"
	"gopkg.in/macaron.v1"

	// 项目模块:
	"github.com/peachdocs/peach/models"                       // 数据库模型
	"github.com/peachdocs/peach/modules/middleware"           // 中间件
	"github.com/peachdocs/peach/modules/setting"              // 配置部分
	"github.com/peachdocs/peach/routers"                      // 路由部分
)

/*
	关键入口:
	- 从命令行获取参数, 执行服务启动操作
	- Action: 执行启动操作
 */
var Web = cli.Command{
	Name:   "web",
	Usage:  "Start Peach web server",
	Action: runWeb,		// todo: 关键入口
	Flags: []cli.Flag{
		stringFlag("config, c", "custom/app.ini", "Custom configuration file path"),
	},
}

/*
	关键入口:
	- 创建项目实例, 并启动服务

 */
func runWeb(ctx *cli.Context) {
	if ctx.IsSet("config") {
		setting.CustomConf = ctx.String("config")
	}
	setting.NewContext()
	models.NewContext()

	log.Info("Peach %s", setting.AppVer)

	//---------------------------------------
	//          关键入口:
	//---------------------------------------
	// 创建web 服务对象
	m := macaron.New()

	// 日志
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())

	// 静态资源处理
	m.Use(macaron.Statics(macaron.StaticOptions{
		SkipLogging: setting.ProdMode,
	}, "custom/public", "public", models.HTMLRoot))

	m.Use(i18n.I18n(i18n.Options{
		Files:       setting.Docs.Locales,
		DefaultLang: setting.Docs.Langs[0],
	}))

	// 模板资源:
	tplDir := "templates"
	if setting.Page.UseCustomTpl {
		tplDir = "custom/templates"
	}
	m.Use(pongo2.Pongoer(pongo2.Options{
		Directory: tplDir,		// 注意: 模板处理部分 pongo2模块 [github.com/go-macaron/pongo2/pongo2.go:194]
	}))

	// 服务中间件:
	m.Use(middleware.Contexter())  // todo: 留意该自定义组件.

	//---------------------------------------
	//          路由配置:
	//---------------------------------------
	m.Get("/", routers.Home)		// 首页部分路由: [github.com/peachdocs/peach/routers/home.go:30]
	m.Get("/docs", routers.Docs)	// 文档部分路由: [github.com/peachdocs/peach/routers/docs.go:44]
	m.Get("/docs/images/*", routers.DocsStatic)			// 文档图片路由: [github.com/peachdocs/peach/routers/docs.go:80]
	m.Get("/docs/*", routers.Protect, routers.Docs)
	m.Post("/hook", routers.Hook)		// 钩子路由: 自动拉取最新资源, 更新文档. [github.com/peachdocs/peach/routers/docs.go:105]
	m.Get("/search", routers.Search)	// 搜索页面路由: [github.com/peachdocs/peach/routers/search.go:24]
	m.Get("/*", routers.Pages)			// [页面遍历搜索, 找到即渲染页面]github.com/peachdocs/peach/routers/home.go:39

	m.NotFound(routers.NotFound)

	listenAddr := fmt.Sprintf("0.0.0.0:%d", setting.HTTPPort)	// 设置服务 IP + 端口
	log.Info("%s Listen on %s", setting.Site.Name, listenAddr)
	log.Fatal("Fail to start Peach: %v", http.ListenAndServe(listenAddr, m))   // 启动服务
}
