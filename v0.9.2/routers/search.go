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

package routers

import (
	"github.com/peachdocs/peach/models"
	"github.com/peachdocs/peach/modules/middleware"
	"github.com/peachdocs/peach/modules/setting"
)

// 搜索功能路由实现
func Search(ctx *middleware.Context) {
	ctx.Data["Title"] = ctx.Tr("search")

	toc := models.Tocs[ctx.Locale.Language()]
	if toc == nil {
		toc = models.Tocs[setting.Docs.Langs[0]]
	}

	q := ctx.Query("q")		// 提取查询关键词
	if len(q) == 0 {
		ctx.Redirect(setting.Page.DocsBaseURL)
		return
	}

	ctx.Data["Keyword"] = q
	ctx.Data["Results"] = toc.Search(q)		// 调用搜索方法, 把搜索结果存入 Results 字段.

	ctx.HTML(200, "search")
}
