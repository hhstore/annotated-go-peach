# Peach v0.9.2 阅读笔记:

![](./public/img/favicon.ico)


## peach 源码结构:


```

-> % tree . -L 2
.
├── Dockerfile
├── LICENSE
├── README.md
├── _note.md
├── bindata.sh
├── cmd                         // 入口目录:
│   ├── cmd.go                  // 
│   ├── init.go                 // 创建 peach 初始化目录.
│   └── web.go                  // 启动 web 服务器
├── conf
│   ├── app.ini
│   └── locale
├── docker
│   ├── README.md
│   ├── build.sh
│   ├── nsswitch.conf
│   ├── s6
│   └── start.sh
├── models
│   ├── doc.go
│   ├── markdown.go             // 解析并渲染 md 文件.
│   ├── protect.go
│   └── toc.go                  // 核心: [难点: GetDoc() 递归实现] github.com/peachdocs/peach/models/toc.go:201
├── modules
│   ├── bindata
│   ├── middleware
│   └── setting
├── peach.go                    // 全局入口: main()
├── public
│   ├── config.codekit
│   ├── css
│   ├── fonts
│   ├── img
│   ├── js
│   └── less
├── routers                     // 路由部分
│   ├── docs.go
│   ├── home.go
│   ├── protect.go
│   └── search.go
└── templates                   // 模板文件
    ├── 404.html
    ├── base.html
    ├── disqus.html
    ├── docs.html
    ├── duoshuo.html
    ├── footer.html
    ├── home.html
    ├── navbar.html
    └── search.html

18 directories, 32 files


```


## 阅读顺序:

- 0. 根本上是: 根据 HTTP 请求执行顺序, 自顶向下, 逐层阅读.
- 1. 首先阅读 :/peach/peach.go
    - 找到cmd.Web 和 cmd.New
- 2. 由上, 阅读: 
    - /peach/cmd/web.go
    - /peach/cmd/init.go
- 3. 由web.go 中 runWeb() 中 m.Get()部分, 阅读:
    - /peach/routers/home.go  // 路由部分






