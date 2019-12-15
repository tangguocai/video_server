### golang开发视频网站
利用golang原生http和httprouter开发的视频网站，包含api、scheduler、stream_server和web四大模块。
####api
业务数据与数据库的交互；处理视频的上传、下载和播放请求转发到stream_server模块。
####stream_server
处理视频的上传、下载和播放；限制请求流量，保护后端。
####scheduler
处理视频的删除、软删除、定期删除。
####web
前后端交互的模块，将请求转发到相应的请求处理。
