# minres
Smart resource management, with caddy+seaweedfs

# 安装
go get github.com/zhangmingfeng/minres

# 前置条件
* redis (用来存储上传文件的信息, 实现断点续传以及文件的本地缓存)
* seaweedfs (小文件存储系统)
