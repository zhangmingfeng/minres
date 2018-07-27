# minres
* Smart resource management, with caddy（https://github.com/mholt/caddy）+seaweedfs（https://github.com/chrislusf/seaweedfs）
* 直接使用weed储存图片，在获取图片使用width和height的时候，weed并没有将获取的这个缩略图缓存，导致每次获取同样尺寸的缩略图weed都会重新处理图片，从而导致前端获取图片时间很长。minres的思路是将用户请求的图片文件按照请求的尺寸缓存起来，当有同样的资源请求的时候直接取缓存。

# features
* weed支持的资源都支持
* 由于web服务基于caddy，支持自动https
* 提供了统一资源上传接口，支持断点续传
* 仅仅图片缩略图支持缓存，从而不用每次到weed请求，提高了性能
* 接口支持鉴权（下个版本支持）

# 安装
go get github.com/zhangmingfeng/minres

# 编译
go build -o minres minres.go

# 前置条件
* redis (用来存储上传文件的信息, 实现断点续传以及文件的本地缓存)
* seaweedfs (小文件存储系统)

# 配置文件
* 配置文件基于caddy配置，详情请看caddy配置：https://caddyserver.com/docs/caddyfile
* 配置文件内容如下： 
  ```
  localhost:61621 {
      cors / { #跨域支持
          origin            *
          methods           GET,POST,OPTIONS
          allowed_headers   X-Requested-With,X-User-Token
      }
      minres { #主配置
          wd_master 127.0.0.1:9333 # weed的master地址
          cache_path /home/zhangmingfeng/Desktop # 缓存文件的位置
      }
      redis { # redis配置
          addr 127.0.0.1:6379
          db 0
      }
  }
  ```
  
# 接口
## 获取上传参数，主要是实现断点续传，上传之前根据上送的信息判断之前是否有上传
* url: /params 
* 上送参数: 
  - fileName string 上传文件名
  - fileGroup string 文件分类组，默认default
  - fileSize int 文件大小，单位字节
  - fileTime int 文件最后修改时间戳
  - chunkSize int 文件上传分片大小，默认1024 * 1024,即1M
 
* 返回参数:
  - code int 返回码，200表示成功
  - msg string 返回信息
  - token string 文件上传token，下面的上传接口需要上送，安全机制
  - uploadUrl string 文件上传地址
  - chunkSize int 分片大小
  - chunk int 当前分片
  - chunks int 文件上传总分片
  - loaded int 文件已经上传的字节数，用于断点续传
  - fileGroup string 文件分类组
 
## 上传接口
* url: /upload
* 上送参数: 
  - token string 文件上传token
  - chunk int 当前上传的分片
  - fileHandle string 文件上传的域，默认'file'
  
* 返回参数
  - code int 返回码，200表示成功
  - msg string 返回信息
  - isFinished bool 文件上传是否完成
  - loaded int 文件已经上传的字节数
  - chunk int 当前分片
  - chunks int 文件上传总分片
  - file object 文件上传完成之后的信息
    - fid string 文件ID,唯一标识，后面请求文件都需要
    - url string 文件access url
    
## 资源aeecss接口
* url: /fetch/{fid}
* fid表示文件ID,文件上传成功之后返回的
* 可以在url后面增加query参数：
  - w int 仅仅在资源是图片的时候有用，定制返回的图片的宽度
  - h int 仅仅在资源是图片的时候有用，定制返回的图片的高度
  - m string 仅仅在资源是图片的时候有用，设置缩略图的属性，目前支持：fit, fill，具体可参见https://github.com/chrislusf/seaweedfs
  - dl bool 是否是下载，如果是下载，即使是浏览器可以直接打开的文件，也会当做附件下载
  
# 启动服务
在linux系统下，可以安装开机自启动，参见：https://caddyserver.com/docs/hook.service
以ubuntu为例：
sudo ./minres -conf /path/to/caddyfile -service install -name minres
service minres start

直接启动服务
./minres -conf /path/to/caddyfile
