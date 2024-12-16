# B 站图片上传工具使用说明

## 功能简介

这是一个用于批量上传图片到哔哩哔哩（B 站）的命令行工具。它可以自动将指定目录下的图片文件上传到 B 站服务器，并生成上传结果报告。

## 使用前准备

1. 获取 B 站 Cookie

   1. 登录 bilibili.com
   2. 打开浏览器开发者工具（F12）
   3. 在开发者工具中找到 Network 标签
   4. 刷新页面，在请求中找到包含完整 cookie 的请求
   5. 复制完整的 cookie 字符串

2. 配置文件设置
   创建 `config.json` 文件，填入以下内容：

```json
{
    "cookie": "你的B站cookie字符串",
    "input_dir": "图片所在目录",
    "output_file": "上传结果保存的JSON文件"
}
```

配置文件示例：

```json
{
    "cookie": "buvid3=A9BDA8DF-DFCA-B8D2-AA14-7BE19F9329B590148infoc; rpdid=|(umu)Y~m~kk0J'u~k|kJYl~R;",
    "input_dir": "./imgs",
    "output_file": "upload_results.json"
}
```

配置项说明：

- cookie: B 站登录凭证
- input_dir: 待上传图片所在的目录路径
- output_file: 上传结果保存的 JSON 文件路径

## 使用方法

### 基本用法

```bash
./uploader -config config.json
```

uploader 为编译后的可执行文件

命令行参数

-config: 指定配置文件路径，默认为 config.json

支持的图片格式

- JPG/JPEG
- PNG
- GIF
- WebP

### 输出结果说明

程序会在指定的输出文件中生成 JSON 格式的上传结果，包含以下信息：

```json
[
    {
        "local_path": "图片本地路径",
        "remote_url": "上传成功后的B站图片URL",
        "success": true,
        "error": ""
    },
    {
        "local_path": "失败图片路径",
        "remote_url": "",
        "success": false,
        "error": "错误信息"
    }
]
```

### 注意事项

1. 确保 cookie 有效且未过期
2. 上传间隔为 1 秒，以避免触发 B 站频率限制
   建议图片大小不超过 5MB
   请确保有足够的磁盘空间存储结果文件

### 常见问题

1. 提示"无法从 cookie 中提取 bili_jct 值"
   检查 cookie 是否完整
   确认 cookie 中包含 bili_jct 字段
   重新登录 B 站获取新的 cookie

2. 上传失败
   检查网络连接
   确认 cookie 未过期
   检查图片格式是否支持
   确认图片大小是否超限

### 技术支持

如遇到问题，请提供以下信息

- 错误信息截图
- 配置文件内容（请去除敏感信息）
- 运行环境信息

### 免责声明

请遵守 B 站相关规定使用本工具
上传的图片内容请符合相关法律法规
本工具仅供学习交流使用
