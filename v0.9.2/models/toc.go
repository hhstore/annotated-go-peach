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

package models

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/Unknwon/com"
	"github.com/mschoch/blackfriday-text"
	"github.com/russross/blackfriday"
	"gopkg.in/ini.v1"

	"github.com/peachdocs/peach/modules/setting"
)

/*
链表结构:

 */
type Node struct {
	Name  string // Name in TOC
	Title string // Name in given language
	text  []byte // Clean text without formatting
	runes []rune

	Plain         bool // Root node without content
	DocumentPath  string
	FileName      string // Full path with .md extension
	Nodes         []*Node	// 链表结构
	LastBuildTime int64
}

func (n *Node) SetText(text []byte) {
	n.text = text
	n.runes = []rune(string(n.text))
}

func (n *Node) Text() []byte {
	return n.text
}

var textRender = blackfridaytext.TextRenderer()
var (
	docsRoot = "data/docs"
	HTMLRoot = "data/html"
)


/*
	关键方法:
	- 解析 md 文件, 提取文件头部的信息


	md 文件示例格式:
	==================================

	---
	name: 创建文档仓库
	---

	# 创建文档仓库

	每一个 Peach 文档仓库都包含两部分内容：

	==================================


 */
func parseNodeName(name string, data []byte) (string, []byte) {
	data = bytes.TrimSpace(data)

	// 出错处理:
	// - md 文件内容长度<3, 且头部无 ---, 直接返回空
	if len(data) < 3 || string(data[:3]) != "---" {
		return name, []byte("")
	}

	// 匹配 --- 尾, 匹配不到, 返回空
	endIdx := bytes.Index(data[3:], []byte("---")) + 3
	if endIdx == -1 {
		return name, []byte("")
	}

	// 切分标题部分:
	// - 标题部分(--- 字符对, 中间的部分),
	// - 根据 换行符 作切分
	opts := strings.Split(strings.TrimSpace(string(string(data[3:endIdx]))), "\n")

	title := name

	// 遍历标题块
	for _, opt := range opts {
		infos := strings.SplitN(opt, ":", 2)	// 根据冒号作切分
		if len(infos) != 2 {
			continue
		}

		switch strings.TrimSpace(infos[0]) {
		case "name":
			title = strings.TrimSpace(infos[1])		// 提取文档标题
		}
	}

	return title, data[endIdx+3:]	// 返回: 文档标题 + 文档内容
}

/*
	关键方法:
		- 解析 markdown 文件, 并渲染成 HTML 页面.
 */
func (n *Node) ReloadContent() error {
	data, err := ioutil.ReadFile(n.FileName)
	if err != nil {
		return err
	}

	// [解析 md 文件, 提取 文档标题+文档内容]
	n.Title, data = parseNodeName(n.Name, data)		// 关键方法 [本模块内实现]
	n.Plain = len(bytes.TrimSpace(data)) == 0

	if !n.Plain {
		n.SetText(bytes.ToLower(blackfriday.Markdown(data, textRender, 0)))
		data = markdown(data)	// 解析并渲染 markdown 数据. [github.com/peachdocs/peach/models/markdown.go:40]
	}

	return n.GenHTML(data)		// 生成 HTML 页面
}

// HTML2JS converts []byte type of HTML content into JS format.
func HTML2JS(data []byte) []byte {
	s := string(data)
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\"", `\"`, -1)
	return []byte(s)
}

// 生成 HTML 文件
func (n *Node) GenHTML(data []byte) error {
	var htmlPath string
	if setting.Docs.Type.IsLocal() {
		htmlPath = path.Join(HTMLRoot, strings.TrimPrefix(n.FileName, setting.Docs.Target))
	} else {
		htmlPath = path.Join(HTMLRoot, strings.TrimPrefix(n.FileName, docsRoot))
	}
	htmlPath = strings.Replace(htmlPath, ".md", ".js", 1)

	buf := new(bytes.Buffer)
	buf.WriteString("document.write(\"")
	buf.Write(HTML2JS(data))	// 格式化字符串
	buf.WriteString("\")")

	n.LastBuildTime = time.Now().Unix()

	// 写数据到 htmlPath.tmp 中
	if err := com.WriteFile(htmlPath+".tmp", buf.Bytes()); err != nil {
		return err
	}
	os.Remove(htmlPath)	// 删除旧的 htmlPath 文件
	return os.Rename(htmlPath+".tmp", htmlPath)	// 重命名 htmlPath.tmp 为 htmlPath
}



//*********************************************************
// 			关键模块: Toc
//
//
//*********************************************************

// Toc represents table of content in a specific language.
type Toc struct {
	RootPath string
	Lang     string		// 语言
	Nodes    []*Node	// 节点
	Pages    []*Node
}


/*
说明:
	- 解析 toc 文件树.(递归实现)
	- 本函数实现, 不好理解, 需画个目录树辅助理解.
	- 注意: 内部有个递归处理.


示例目录结构:

-> % tree zh-CN
zh-CN
├── advanced
│   ├── README.md
│   └── config_cheat_sheet.md
├── faqs
│   └── README.md
├── howto
│   ├── README.md
│   ├── documentation.md
│   ├── extension.md
│   ├── navbar.md
│   ├── pages.md
│   ├── protect_resources.md
│   ├── static_resources.md
│   ├── templates.md
│   ├── upgrade.md
│   └── webhook.md
└── intro
    ├── README.md
    ├── getting_started.md
    ├── installation.md
    └── roadmap.md

 */
// GetDoc should only be called by top level toc.
func (t *Toc) GetDoc(name string) (*Node, bool) {
	name = strings.TrimPrefix(name, "/")

	// Returns first node whatever avaiable as default.
	if len(name) == 0 {
		if len(t.Nodes) == 0 ||
			t.Nodes[0].Plain {
			return nil, false
		}
		return t.Nodes[0], false	// 默认返回: 首个节点
	}

	infos := strings.Split(name, "/")        // 根据/, 切分

	// Dir node. [目录节点: 处理目录]
	if len(infos) == 1 {
		for i := range t.Nodes {
			if t.Nodes[i].Name == infos[0] {
				return t.Nodes[i], false
			}
		}
		return nil, false
	}

	// File node. [文件节点: 处理文件]
	for i := range t.Nodes {	// 遍历 节点集
		if t.Nodes[i].Name == infos[0] {

			// 遍历目标节点的子节点集.
			for j := range t.Nodes[i].Nodes {	// 子节点集
				if t.Nodes[i].Nodes[j].Name == infos[1] {
					if com.IsFile(t.Nodes[i].Nodes[j].FileName) {
						return t.Nodes[i].Nodes[j], false
					}

					// If not default language, try again.
					n, _ := Tocs[setting.Docs.Langs[0]].GetDoc(name)	// todo: 递归处理
					return n, true
				}
			}
		}
	}

	return nil, false
}


// 全文检索的存储结果格式
type SearchResult struct {
	Title string
	Path  string
	Match string
}

func (n *Node) adjustRange(start int) (int, int) {
	start -= 20
	if start < 0 {
		start = 0
	}

	length := len(n.runes)
	end := start + 230
	if end > length {
		end = length
	}
	return start, end
}


/*
功能:
	- 全文搜索(自己实现)
说明:
	- 学习该实现方法, 很简单. 关键词匹配+计数.

 */
func (t *Toc) Search(q string) []*SearchResult {
	if len(q) == 0 {
		return nil
	}
	q = strings.ToLower(q)

	results := make([]*SearchResult, 0, 5)

	// Dir node.
	for i := range t.Nodes {
		if idx := bytes.Index(t.Nodes[i].Text(), []byte(q)); idx > -1 {
			// 关键词匹配+计数
			start, end := t.Nodes[i].adjustRange(utf8.RuneCount(t.Nodes[i].Text()[:idx]))
			results = append(results, &SearchResult{
				Title: t.Nodes[i].Title,
				Path:  t.Nodes[i].Name,
				Match: string(t.Nodes[i].runes[start:end]),
			})
		}
	}

	// File node.
	for i := range t.Nodes {
		for j := range t.Nodes[i].Nodes {
			if idx := bytes.Index(t.Nodes[i].Nodes[j].Text(), []byte(q)); idx > -1 {
				// 关键词匹配+计数
				start, end := t.Nodes[i].Nodes[j].adjustRange(utf8.RuneCount(t.Nodes[i].Nodes[j].Text()[:idx]))
				results = append(results, &SearchResult{
					Title: t.Nodes[i].Nodes[j].Title,
					Path:  path.Join(t.Nodes[i].Name, t.Nodes[i].Nodes[j].Name),
					Match: string(t.Nodes[i].Nodes[j].runes[start:end]),
				})
			}
		}
	}

	return results
}


/******************************************
	            TOC 部分
	说明:
	-


*******************************************/
var (
	tocLocker = sync.Mutex{}	// 加锁
	Tocs      map[string]*Toc
)

// 初始化
func initToc(localRoot string) (map[string]*Toc, error) {
	tocPath := path.Join(localRoot, "TOC.ini")
	if !com.IsFile(tocPath) {
		return nil, fmt.Errorf("TOC not found: %s", tocPath)
	}

	// Generate Toc.
	tocCfg, err := ini.Load(tocPath)
	if err != nil {
		return nil, fmt.Errorf("Fail to load TOC.ini: %v", err)
	}

	tocs := make(map[string]*Toc)
	for _, lang := range setting.Docs.Langs {
		toc := &Toc{
			RootPath: localRoot,
			Lang:     lang,
		}
		dirs := tocCfg.Section("").KeyStrings()
		toc.Nodes = make([]*Node, 0, len(dirs))
		for _, dir := range dirs {
			dirName := tocCfg.Section("").Key(dir).String()
			fmt.Println(dirName + "/")
			files := tocCfg.Section(dirName).KeyStrings()

			// Skip empty directory.
			if len(files) == 0 {
				continue
			}

			documentPath := path.Join(dirName, tocCfg.Section(dirName).Key(files[0]).String())
			dirNode := &Node{
				Name:         dirName,
				DocumentPath: documentPath,
				FileName:     path.Join(localRoot, lang, documentPath) + ".md",
				Nodes:        make([]*Node, 0, len(files)-1),
			}
			toc.Nodes = append(toc.Nodes, dirNode)

			for _, file := range files[1:] {
				fileName := tocCfg.Section(dirName).Key(file).String()
				fmt.Println(strings.Repeat(" ", len(dirName))+"|__", fileName)

				documentPath = path.Join(dirName, fileName)
				node := &Node{
					Name:         fileName,
					DocumentPath: documentPath,
					FileName:     path.Join(localRoot, lang, documentPath) + ".md",
				}
				dirNode.Nodes = append(dirNode.Nodes, node)
			}
		}

		// Single pages.
		pages := tocCfg.Section("pages").KeyStrings()
		toc.Pages = make([]*Node, 0, len(pages))
		for _, page := range pages {
			pageName := tocCfg.Section("pages").Key(page).String()
			fmt.Println(pageName)

			toc.Pages = append(toc.Pages, &Node{
				Name:         pageName,
				DocumentPath: pageName,
				FileName:     path.Join(localRoot, lang, pageName) + ".md",
			})
		}

		tocs[lang] = toc
	}
	return tocs, nil
}

/*
	关键方法:
		- 重载 docs 文档目录.
		- 拉取最新资源: git clone / git pull 远端最新文档资源.
		- 有几个细节, 需要留意.
 */
func ReloadDocs() error {
	tocLocker.Lock()			// 加锁处理
	defer tocLocker.Unlock()

	localRoot := setting.Docs.Target

	// Fetch docs from remote.
	if setting.Docs.Type.IsRemote() {
		localRoot = docsRoot

		absRoot, err := filepath.Abs(localRoot)
		if err != nil {
			return fmt.Errorf("filepath.Abs: %v", err)
		}

		// Clone new or pull to update.
		if com.IsDir(absRoot) {		// 目录已存在, 只需更新
			stdout, stderr, err := com.ExecCmdDir(absRoot, "git", "pull")	// git pull absRoot , 拉取最新代码
			if err != nil {
				return fmt.Errorf("Fail to update docs from remote source(%s): %v - %s", setting.Docs.Target, err, stderr)
			}
			fmt.Println(stdout)
		} else {		// 目录不存在, 创建目录, 并 git clone 下来.
			os.MkdirAll(filepath.Dir(absRoot), os.ModePerm)
			stdout, stderr, err := com.ExecCmd("git", "clone", setting.Docs.Target, absRoot)	// git clone 操作
			if err != nil {
				return fmt.Errorf("Fail to clone docs from remote source(%s): %v - %s", setting.Docs.Target, err, stderr)
			}
			fmt.Println(stdout)
		}
	}

	if !com.IsDir(localRoot) {
		return fmt.Errorf("Documentation not found: %s - %s", setting.Docs.Type, localRoot)
	}

	tocs, err := initToc(localRoot)		// 初始化 [github.com/peachdocs/peach/models/toc.go:366]
	if err != nil {
		return fmt.Errorf("initToc: %v", err)
	}
	initDocs(tocs, localRoot)	// [github.com/peachdocs/peach/models/doc.go:63]
	Tocs = tocs
	return reloadProtects(localRoot)	// [github.com/peachdocs/peach/models/protect.go:41]
}
