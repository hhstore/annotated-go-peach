# Peach v0.9.2 导读:


![](./public/img/favicon.ico)

## 项目源码结构:


```

-> % tree ./v0.9.2 -L 2 

./v0.9.2
├── README.md
├── bindata.sh
├── cmd
│   ├── cmd.go
│   ├── init.go
│   └── web.go
├── conf
│   ├── app.ini
│   └── locale
├── models
│   ├── doc.go
│   ├── markdown.go
│   ├── protect.go
│   └── toc.go
├── modules
│   ├── bindata
│   ├── middleware
│   └── setting
├── peach.go
├── public
│   ├── config.codekit
│   ├── css
│   ├── fonts
│   ├── img
│   ├── js
│   └── less
├── routers
│   ├── docs.go
│   ├── home.go
│   ├── protect.go
│   └── search.go
└── templates
    ├── 404.html
    ├── base.html
    ├── disqus.html
    ├── docs.html
    ├── duoshuo.html
    ├── footer.html
    ├── home.html
    ├── navbar.html
    └── search.html

```



