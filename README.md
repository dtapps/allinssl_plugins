# 插件库

AllinSSL 第三方插件库

## 功能模块

| 模块名称            | 功能描述          | 对应文件            | 开发版本 | 测试状态  |
| ------------------- | ----------------- | ------------------- | -------- | --------- |
| 天翼云              | CDN 加速          | ctyun               |          | ✅ 已测试 |
| 天翼云              | 全站加速          | ctyun               |          | ✅ 已测试 |
| 天翼云              | 证书管理          | ctyun               |          | ✅ 已测试 |
| 天翼云              | 边缘安全加速平台  | ctyun               |          | ⚠️ 未测试 |
| Nginx Proxy Manager | SSL Certificates  | nginx_proxy_manager | 2.12.3   | ✅ 已测试 |
| Nginx Proxy Manager | Proxy Hosts       | nginx_proxy_manager | 2.12.3   | ✅ 已测试 |
| Nginx Proxy Manager | Redirection Hosts |                     |          | ❌ 计划中 |
| Nginx Proxy Manager | Streams           |                     |          | ❌ 计划中 |
| Proxmox VE          | 证书管理          | proxmox             | 8.3.3    | ✅ 已测试 |
| Synology            | 证书管理          | synology            | 7.2.2    | ✅ 已测试 |
| 堡塔云WAF            | 网站         | aawaf            | 6.1 / 6.2    | ✅ 已测试，推荐用官方 |
| 耗子面板            | 网站         | ratpanel            | 2.5.5    | ✅ 已测试 |
| 耗子面板            | 证书管理         | ratpanel            | 2.5.5    | ✅ 已测试 |
| 雷池WAF            | 防护应用（API更新接口不合理）         | safeline            | 8.10.1    | ✅ 已测试，推荐用官方 |
| 小皮面板            | 网站（API接口不合理）         | xppanel            | 1.3.10    | ⚠️ 未测试，不推荐使用 |


## 天翼云 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在 AllinSSL 后台「添加授权 API」选择 ctyun 插件
3. 配置参数（JSON 格式）：

```json
{
  "access_key": "AccessKey",
  "secret_key": "SecretKey"
}
```

## Nginx Proxy Manager 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在 AllinSSL 后台「添加授权 API」选择 nginx_proxy_manager 插件
3. 配置参数（JSON 格式）：

```json
{
  "url": "主机IP或域名，包含协议和端口，http://example.com 或 http://0.0.0.0:81",
  "email": "登录邮箱",
  "password": "登录密码"
}
```

## Proxmox VE 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在 AllinSSL 后台「添加授权 API」选择 proxmox 插件
3. 配置参数（JSON 格式）：

```json
{
  "url": "主机IP或域名，包含协议和端口，例如：https://example.com 或 https://0.0.0.0:8006",
  "node": "节点名称，例如：pve",
  "user": "用户名和领域，例如：root@pam",
  "token_id": "令牌 ID",
  "token_secret": "令牌 密钥"
}
```

## Synology 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在 AllinSSL 后台「添加授权 API」选择 synology 插件
3. 配置参数（JSON 格式）：

```json
{
  "url": "主机IP或域名，包含协议和端口，例如：https://example.com 或 https://0.0.0.0:5001",
  "username": "用户名",
  "password": "密码"
}
```

4.  由于 `隐私` 和 `双重验证` 问题，建议建个独立的账号给插件使用，开启双重验证会导致登录失败！

## 堡塔云WAF 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在 AllinSSL 后台「添加授权 API」选择 aawaf 插件
3. 配置参数（JSON 格式）：

```json
{
  "url": "主机IP或域名，包含协议和端口，https://example.com 或 https://0.0.0.0:8379",
  "api_key": "接口密钥",
}
```

## 雷池WAF 使用示例

1. 将模块文件放入 AllinSSL 插件目录(plugins)
2. 在 AllinSSL 后台「添加授权 API」选择 aawaf 插件
3. 配置参数（JSON 格式）：

```json
{
  "url": "主机IP或域名，包含协议和端口，https://example.com 或 https://0.0.0.0:9443",
  "api_token": "接口密钥",
}
```
