# 更新日志 
## (1.0.2) - 2025-07-01
### 修复 
* [certificates] 修复上传接口的响应类型不是JSON类型导致解析失败引起的提示错误 ([#4](https://github.com/dtapps/allinssl_plugins/issues/4))

## (1.0.1) - 2025-06-30
### 修复 
* [certificates] 修复判断证书唯一性时未考虑的问题，增加对比证书域名列表和过期时间
* [certificates] 修复没有正确反馈接口状态
