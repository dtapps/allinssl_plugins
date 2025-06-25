# 插件库

AllinSSL 第三方插件库

## 功能模块

| 模块名称 | 功能描述 | 对应文件 |
|---------|---------|---------|
| 天翼云 | CDN加速 | ctyun/cdn |
| 天翼云 | 全站加速 | ctyun/icdn |
| 天翼云 | 边缘安全加速平台 | ctyun/accessone |
| 天翼云 | 证书管理服务 | ctyun/ccms |
| Nginx Proxy Manager | Proxy Hosts | nginx_proxy_manager/proxy_host |
| Nginx Proxy Manager | SSL Certificates | nginx_proxy_manager/certificate |

## 天翼云 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在AllinSSL后台「添加授权API」选择 ctyun 插件
3. 配置参数（JSON格式）：
```json
   {
     "access_key": "您的AccessKey",
     "secret_key": "您的SecretKey"
   }
```

## Nginx Proxy Manager 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在AllinSSL后台「添加授权API」选择 nginx_proxy_manager 插件
3. 配置参数（JSON格式）：
```json
   {
     "url": "您的网站包含端口",
     "email": "您的邮箱",
     "password": "您密码"
   }
```
