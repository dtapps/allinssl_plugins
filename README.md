# 插件库

AllinSSL 第三方插件库

## 功能模块

| 模块名称            | 功能描述          | 对应文件                             | 开发版本 | 测试状态  |
| ------------------- | ----------------- | ------------------------------------ | -------- | --------- |
| 天翼云              | CDN 加速          | ctyun/cdn                            |          | ✅ 已测试 |
| 天翼云              | 全站加速          | ctyun/icdn                           |          | ✅ 已测试 |
| 天翼云              | 证书管理          | ctyun/ccms                           |          | ✅ 已测试 |
| 天翼云              | 边缘安全加速平台  | ctyun/accessone                      |          | ⚠️ 未测试 |
| Nginx Proxy Manager | SSL Certificates  | nginx_proxy_manager/certificate      | 2.12.3   | ✅ 已测试 |
| Nginx Proxy Manager | Proxy Hosts       | nginx_proxy_manager/proxy_host       | 2.12.3   | ✅ 已测试 |
| Nginx Proxy Manager | Redirection Hosts | nginx_proxy_manager/redirection_host |          | ❌ 计划中 |
| Nginx Proxy Manager | Streams           | nginx_proxy_manager/stream           |          | ❌ 计划中 |
| Proxmox VE          | 证书管理          | proxmox/certificate                  | 8.3.3    | ✅ 已测试 |
| Synology            | 证书管理          | synology/certificate                 | 7.2.2    | ✅ 已测试 |

## 天翼云 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在 AllinSSL 后台「添加授权 API」选择 ctyun 插件
3. 配置参数（JSON 格式）：

```json
{
  "access_key": "您的AccessKey",
  "secret_key": "您的SecretKey"
}
```

## Nginx Proxy Manager 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在 AllinSSL 后台「添加授权 API」选择 nginx_proxy_manager 插件
3. 配置参数（JSON 格式）：

```json
{
  "url": "您的网站，包含协议和端口",
  "email": "您的邮箱",
  "password": "您的密码"
}
```

## Proxmox VE 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在 AllinSSL 后台「添加授权 API」选择 proxmox 插件
3. 配置参数（JSON 格式）：

```json
{
  "url": "您的网站，包含协议和端口",
  "node": "您的节点名称",
  "user": "您的用户名和领域",
  "token_id": "您的令牌ID",
  "token_secret": "您的令牌密钥"
}
```

## Synology 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在 AllinSSL 后台「添加授权 API」选择 synology 插件
3. 配置参数（JSON 格式）：

```json
{
  "url": "您的网站，包含协议和端口",
  "username": "您的用户名",
  "password": "您的密码"
}
```

4.  由于 `隐私` 和 `双重验证` 问题，建议建个独立的账号给插件使用，开启双重验证会导致登录失败！
