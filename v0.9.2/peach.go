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

// Peach is a web server for multi-language, real-time synchronization and searchable documentation.
package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"

	"github.com/peachdocs/peach/cmd"
	"github.com/peachdocs/peach/modules/setting"
)

const APP_VER = "0.9.2.1214"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "Peach"
	app.Usage = "Modern Documentation Server"
	app.Version = APP_VER
	app.Author = "Unknwon"
	app.Email = "u@gogs.io"


	/*
		全局关键:
		- 1. 注册2个命令服务. 通过命令行传参方式执行.
		- 2. 使用方法:
			- peach new -target=my.peach    // 在当前目录下, 创建peach工程目录.
			- peach web                     // 启动 web 服务器.
	 */
	app.Commands = []cli.Command{
		cmd.Web,	// todo: 命令1: 启动 web 服务 [全局关键入口]
		cmd.New,	// 命令2: 用来生成 peach 项目原型结构.
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)

	// 从命令行获取参数,并执行
	app.Run(os.Args)
}
