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

	"github.com/russross/blackfriday"
)

var (
	tab    = []byte("\t")
	spaces = []byte("    ")
)

type MarkdownRender struct {
	blackfriday.Renderer
}

// tab 符号转换成4空格
func (mr *MarkdownRender) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	var tmp bytes.Buffer
	mr.Renderer.BlockCode(&tmp, text, lang)
	out.Write(bytes.Replace(tmp.Bytes(), tab, spaces, -1))	// 把 tab 替换成4空格, 然后写出到 out 里.
}

// markdown 处理:
//	- 解析并渲染 markdown
func markdown(raw []byte) []byte {
	htmlFlags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	renderer := &MarkdownRender{
		Renderer: blackfriday.HtmlRenderer(htmlFlags, "", ""),
	}

	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS

	return blackfriday.Markdown(raw, renderer, extensions)	// 调用第三方包, 解析并渲染 markdown.
}
