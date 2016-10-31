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

package middleware

import (
	"strings"

	"gopkg.in/macaron.v1"

	"github.com/peachdocs/peach/modules/setting"
)


/*
	说明:
	- 简单封装: macaron.Context
	- 方便扩展, 但目前的源码, 没有作扩展
 */
type Context struct {
	*macaron.Context
}


/*
	说明:
	- /peach/cmd/web.go 中引用
	- 添加了一些初始化参数

 */
func Contexter() macaron.Handler {
	return func(c *macaron.Context) {

		// 使用上面自定义的 Context
		ctx := &Context{
			Context: c,
		}
		c.Map(ctx)  //  todo: ?? 留意Map()实现

		ctx.Data["Link"] = strings.TrimSuffix(ctx.Req.URL.Path, ".html")

		// 设置全局配置参数
		ctx.Data["AppVer"] = setting.AppVer
		ctx.Data["Site"] = setting.Site
		ctx.Data["Page"] = setting.Page
		ctx.Data["Navbar"] = setting.Navbar
		ctx.Data["Asset"] = setting.Asset
		ctx.Data["Extension"] = setting.Extension
	}
}
