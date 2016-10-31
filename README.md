# annotated-go-peach

- 本项目对 peach 源码作详细注解.
- 可参考阅读笔记和提交记录阅读.

## peach 简介:

- [peach - github](https://github.com/peachdocs/peach)
- [peach - 官网介绍](https://peachdocs.org/)
- 官方说明: peach 是一款支持多语言、实时同步以及全文搜索功能的 Web 文档服务器.
- 使用 go 语言开发, 是 macaron 框架作者[Unknwon](https://github.com/Unknwon) 应用 macaron 的最佳实践.
- 特点: 源码结构组织良好, 代码很少, 且质量很高.


## peach 源码版本:

- [v0.9.2](./v0.9.2)
    - 为当前最新 release 版本(且是第一个 release).
    - 作者 release 频率不高.

## peach 源码结构:

```

-> % tree v0.9.2 -L 1 

v0.9.2
├── README.md
├── bindata.sh
├── cmd            # 项目次级入口
├── conf
├── models
├── modules
├── peach.go       # 项目全局入口: main()
├── public
├── routers        # 路由部分
└── templates      # 模板文件


```


## 阅读建议:

- 阅读源码之前, 参考官方文档, 在本机运行起来, 方便和源码具体细节对照.
- 从 `peach.go` 入口文件读起. 自顶向下, 逐层阅读.
