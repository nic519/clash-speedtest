# Clash-SpeedTest

基于 Clash/Mihomo 核心的测速工具，快速测试你的节点速度。

Features:
1. 无需额外的配置，直接将 Clash/Mihomo 配置本地文件路径或者订阅地址作为参数传入即可
2. 支持 Proxies 和 Proxy Provider 中定义的全部类型代理节点，兼容性跟 Mihomo 一致
3. 不依赖额外的 Clash/Mihomo 进程实例，单一工具即可完成测试
4. 代码简单而且开源，不发布构建好的二进制文件，保证你的节点安全

<img width="1346" height="682" alt="Image" src="https://github.com/user-attachments/assets/9fea1d47-251f-4c49-b059-05b5962d4e72" />

## Prerequisites/注意事项

### OpenWRT 环境
在 OpenWRT 环境下使用本工具时，建议临时关闭 OpenClash/Clash/Mihomo 等代理服务，以避免路由冲突影响测速结果的准确性。或者给 OpenClash/Clash/Mihomo 配置进程规则绕过代理：
```
rules:
  - PROCESS-NAME,clash-speedtest,DIRECT
```

### Windows CMD 用户
在 Windows CMD 中使用时，如果订阅地址包含 `&` 字符，必须使用双引号而非单引号：
```bash
# 正确
> clash-speedtest -c "https://domain.com/api/v1/client/subscribe?token=secret&flag=meta"

# 错误
> clash-speedtest -c 'https://domain.com/api/v1/client/subscribe?token=secret&flag=meta'
```

## 使用方法

```bash
# 支持从源码安装，或从 Release 里下载由 Github Action 自动构建的二进制文件
> go install github.com/nic519/clash-speedtest@latest

# 查看版本
> clash-speedtest -v

# 查看帮助
> clash-speedtest -h

Usage of clash-speedtest:
  -c string
        configuration file path, also support http(s) url
  -ua string
        User-Agent for fetching config from http(s) URL (default: mihomo kernel UA, e.g. mihomo/1.10.0)
  -f string
        filter proxies by name, use regexp (default ".*")
  -b string
        block proxies by keywords, use | to separate multiple keywords (example: -b 'rate|x1|1x')
  -server-url string
        server url or direct download url (default "https://dl.google.com/chrome/mac/universal/stable/GGRO/googlechrome.dmg")
  -latency-url string
        url used for latency testing, defaults to server-url target
  -latency-timeout duration
        timeout for each latency request, defaults to --timeout
  -probe-url string
        optional URL to probe through each proxy
  -probe-method string
        HTTP method for probe URL (default "GET")
  -probe-timeout duration
        timeout for probing each proxy (default 5s)
  -probe-fields string
        comma-separated JSON field mappings for probe output, e.g. ip=ip,country=country_name
  -speed-mode string
        speed test mode: fast, download, full (default "download")
  -download-size int
        download size for testing proxies (default 50MB)
  -upload-size int
        upload size for testing proxies (full mode only) (default 20MB)
  -timeout duration
        timeout for testing proxies (default 5s)
  -concurrent int
        download concurrent size (default 4)
  -proxy-concurrent int
        proxy test concurrent size (default 1)
  -output string
        output config file path (default "")
  -max-latency duration
        filter latency greater than this value (default 1s)
  -max-packet-loss float
        filter packet loss greater than this value(unit: %) (default 100)
  -min-download-speed float
        filter speed less than this value(unit: MB/s) (default 5)
  -min-upload-speed float
        filter upload speed less than this value(unit: MB/s, full mode only) (default 2)
  -rename
        rename nodes with IP location and speed
  -fast
        fast mode (alias for --speed-mode fast)
  -gist-token string
        GitHub personal access token for gist upload
  -gist-address string
        gist URL or ID for uploading output file (filename uses output basename)
  -repo-token string
        GitHub personal access token for repository file upload
  -repo-address string
        repository URL or owner/repo for uploading output file
  -repo-file-path string
        repository file path for uploading output file (default: output basename)
  -repo-branch string
        repository branch for uploading output file (default: repository default branch)

# 演示：

# 1. 测试全部节点，使用 HTTP 订阅地址
# 请在订阅地址后面带上 flag=meta 参数，否则无法识别出节点类型
> clash-speedtest -c 'https://domain.com/api/v1/client/subscribe?token=secret&flag=meta'

# 2. 测试香港节点，使用正则表达式过滤，使用本地文件
> clash-speedtest -c ~/.config/clash/config.yaml -f 'HK|港'
节点                                        	带宽          	延迟
Premium|广港|IEPL|01                        	484.80KB/s  	815.00ms
Premium|广港|IEPL|02                        	N/A         	N/A
Premium|广港|IEPL|03                        	2.62MB/s    	333.00ms
Premium|广港|IEPL|04                        	1.46MB/s    	272.00ms
Premium|广港|IEPL|05                        	3.87MB/s    	249.00ms

# 3. 当然你也可以混合使用
> clash-speedtest -c "https://domain.com/api/v1/client/subscribe?token=secret&flag=meta,/home/.config/clash/config.yaml"

# 4. 筛选出延迟低于 800ms 且下载速度大于 5MB/s 的节点，并输出到 filtered.yaml
> clash-speedtest -c "https://domain.com/api/v1/client/subscribe?token=secret&flag=meta" -output filtered.yaml -max-latency 800ms -min-speed 5
# 筛选后的配置文件可以直接粘贴到 Clash/Mihomo 中使用，或是贴到 Github\Gist 上通过 Proxy Provider 引用。

# 5. 使用 -rename 选项按照 IP 地区和下载速度重命名节点
> clash-speedtest -c config.yaml -output result.yaml -rename
# 重命名后的节点名称格式：🇺🇸 US 001 | ⬇️ 15.67MB/s
# 包含国旗 emoji、国家代码和下载速度

# 6. 快速测试模式
> clash-speedtest -f 'HK' -fast -c ~/.config/clash/config.yaml
# 此命令将只测试节点延迟，跳过其他测试项目，适用于：
# - 快速检查节点是否可用
# - 只需要检查延迟的场景
# - 需要快速得到测试结果的场景
🇭🇰 香港 HK-10 100% |██████████████████| (20/20, 13 it/min)
序号    节点名称                类型            延迟
1.      🇭🇰 香港 HK-01           Trojan          657ms
2.      🇭🇰 香港 HK-20           Trojan          649ms
3.      🇭🇰 香港 HK-15           Trojan          674ms
4.      🇭🇰 香港 HK-19           Trojan          649ms
5.      🇭🇰 香港 HK-12           Trojan          667ms

# 6.1 批量测试多个网站延迟并输出 CSV 报表
> scripts/multi_site_latency_report.sh -c ~/.config/clash/config.yaml -f 'HK|港'
# 默认会测试 YouTube / X / GitHub，并生成：
# - reports/YYYYMMDD-HHmm/YYYYMMDD-HHmm-details.csv
# - reports/YYYYMMDD-HHmm/YYYYMMDD-HHmm-summary.csv
#
# summary.csv 会按节点横向展开各网站延迟，便于后期汇总、对比。
#
# 也可以自定义站点：
> scripts/multi_site_latency_report.sh -c ~/.config/clash/config.yaml -f 'HK|港' \
    -s 'YouTube|https://www.youtube.com/generate_204' \
    -s 'X|https://x.com' \
    -s 'GitHub|https://github.com'

# 6.2 通过节点出口访问 probe URL，输出出口 IP / 地区 / ASN
> clash-speedtest -f 'HK|港' -fast -c ~/.config/clash/config.yaml \
    --latency-url "https://www.youtube.com/generate_204" \
    --latency-timeout "8s" \
    --probe-url "https://api.ip.sb/geoip/" \
    --probe-timeout "8s" \
    --probe-fields "ip=ip,country=country,country_code=country_code,region=region,city=city,asn=asn,org=organization"
#
# probe 请求会通过当前正在测试的代理节点发出。
# TSV 输出中会追加 Probe URL、Probe 延迟、Probe 状态、probe.ip、probe.country、probe.asn 等列。

# 7. 上传到 GitHub Gist
> clash-speedtest -c config.yaml -output result.yaml -gist-token "ghp_xxx" -gist-address "https://gist.github.com/user/abc123"
# 测试完成后，会将 result.yaml 上传到指定的 Gist，文件名与 -output 保持一致（去除目录前缀）
# gist-address 可以是完整的 Gist URL，也可以是 Gist ID（如 abc123）
# Gist/Repo 上传与远程配置 URL 加载默认遵循环境代理变量（HTTPS_PROXY/HTTP_PROXY）。

# 8. 上传到 GitHub 仓库文件（默认写入 output 文件名）
> clash-speedtest -c config.yaml -output result.yaml -repo-token "ghp_xxx" -repo-address "user/repo"
# 测试完成后，会将 result.yaml 上传到仓库默认分支下的 result.yaml

# 9. 上传到 GitHub 仓库指定分支与路径
> clash-speedtest -c config.yaml -output result.yaml -repo-token "ghp_xxx" -repo-address "https://github.com/user/repo" -repo-file-path "configs/subscriptions/result.yaml" -repo-branch "main"
```

## GitHub Token 创建与权限

### 1) 更新 Gist（`-gist-token`）

推荐使用 **Personal access tokens (classic)**：

1. 打开 GitHub `Settings` → `Developer settings` → `Personal access tokens` → `Tokens (classic)`。
2. 点击 `Generate new token (classic)`。
3. 仅勾选最小权限：`gist`。
4. 生成后复制 token，作为 `-gist-token` 传入。

最小权限结论：
- `gist`：必需（用于通过 API 更新 Gist 文件）。

### 2) 更新仓库文件（`-repo-token`）

可选两种 token：

#### A. Fine-grained PAT（推荐）

1. 打开 GitHub `Settings` → `Developer settings` → `Personal access tokens` → `Fine-grained tokens`。
2. `Repository access` 选择目标仓库（建议 `Only select repositories`）。
3. 在 `Repository permissions` 中设置：
   - `Contents`: **Read and write**（必需）
4. 生成后复制 token，作为 `-repo-token` 传入。

#### B. Tokens (classic)

- 更新**公开仓库**文件：至少 `public_repo`。
- 更新**私有仓库**文件：至少 `repo`。

最小权限结论：
- Fine-grained PAT：`Contents: Read and write`。
- Classic PAT：公有仓库 `public_repo`，私有仓库 `repo`。

### 常见权限问题

- `401 Unauthorized`：token 无效、过期，或复制时有空格/换行。
- `403 Forbidden`：token 权限不足，或目标分支启用了保护策略（可能禁止直接 push/commit）。
- `404 Not Found`：仓库地址/路径/分支不正确，或 token 对该仓库不可见。

> 安全建议：不要把 token 提交到仓库；优先通过环境变量或 CI Secret 注入。

## 测速原理

通过 HTTP GET 请求下载指定大小的文件，默认使用 https://dl.google.com/chrome/mac/universal/stable/GGRO/googlechrome.dmg 进行测试，计算下载时间得到下载速度。因为 speed.cloudflare.com 容易返回 403，所以默认不再使用它作为测速入口。

当 server-url 不带 path 时 (使用 https://speed.cloudflare.com 或自建测速服务)，使用 /__down 和 /__up 完成下载与上传测试。
当 server-url 带 path 时，会被识别为直接下载地址，只进行下载测速。

延迟测试默认访问 `server-url` 对应的目标。也可以通过 `latency-url` 单独指定要测试访问延迟的网址，例如只筛选香港节点并测试访问 YouTube 的延迟：

```shell
clash-speedtest -c config.yaml -f 'HK|港' --speed-mode fast \
  --latency-url "https://www.youtube.com/generate_204" \
  --latency-timeout "8s"
```

`latency-timeout` 控制每一次延迟探测请求最多等待多久，默认跟随 `timeout`。`max-latency` 只负责在输出配置时筛掉平均延迟过高的节点，不再作为访问 YouTube、X、GitHub 等目标站点时的 HTTP 请求超时时间。

`proxy-concurrent` 控制同时测试多少个代理节点。它和 `concurrent` 不同：`concurrent` 只控制单个节点测速时的下载分片并发；`proxy-concurrent` 控制节点级 worker pool，可以避免慢节点或坏节点把整个测试队列串行堵住。默认值为 `1`，以保持旧版串行测速语义；需要加速时可以显式调到 `2`、`4` 或 `8`，并发越高越可能因为本机网络、目标站点限流或同机场资源争抢而影响结果。

probe 测试和延迟测试使用同一个代理拨号链路：工具会先为当前节点创建 HTTP client，再通过这个 client 访问 `probe-url`。因此访问 `https://api.ip.sb/geoip/` 时拿到的是该代理节点的出口信息，而不是本机公网 IP。`probe-fields` 用于把 JSON 响应字段映射到 TSV 输出列，换 IP 检测后端时通常只需要调整 `probe-url` 和 `probe-fields`。

如果你确认 https://speed.cloudflare.com 可以访问并希望测试上传，请显式设置为 full 模式，例如：
```shell
clash-speedtest --server-url "https://speed.cloudflare.com" --speed-mode full
```
或者你也可以自己搭建一个测速服务器，用来测试下载和上传速度：

```shell
# 在您需要进行测速的服务器上安装和启动测速服务器
> go install github.com/faceair/clash-speedtest/download-server@latest
> download-server

# 此时在本地使用 http://your-server-ip:8080 作为 server-url 即可
> clash-speedtest --server-url "http://your-server-ip:8080" --speed-mode full
```


测试结果：
1. 带宽 是指下载指定大小文件的速度，即一般理解中的下载速度。当这个数值越高时表明节点的出口带宽越大。
2. 延迟 是指 HTTP GET 请求拿到第一个字节的的响应时间，即一般理解中的 TTFB。当这个数值越低时表明你本地到达节点的延迟越低，可能意味着中转节点有 BGP 部署、出海线路是 IEPL、IPLC 等。

请注意带宽跟延迟是两个独立的指标，两者并不关联：
1. 可能带宽很高但是延迟也很高，这种情况下你下载速度很快但是打开网页的时候却很慢，可能是是中转节点没有 BGP 加速，但出海线路带宽很充足。
2. 可能带宽很低但是延迟也很低，这种情况下你打开网页的时候很快但是下载速度很慢，可能是中转节点有 BGP 加速，但出海线路的 IEPL、IPLC 带宽很小。

## License

[GPL-3.0](LICENSE)

## 重新编译

这个项目默认不提交构建产物。修改 Go 代码后，可以在仓库根目录重新编译本机二进制：

```bash
go test ./speedtester ./output ./tui
go build -ldflags "-X main.version=0.1.3 -X main.commit=$(git rev-parse --short HEAD)" -o ~/go/bin/clash-speedtest .
~/go/bin/clash-speedtest -v
```

如果要给 Latency Compass 桌面 app 使用，`main.version` 必须和桌面 app 的 `package.json` 版本一致。当前桌面 app 要求：

```text
clash-speedtest version 0.1.3
```

版本不一致时，桌面 app 会把依赖状态标记为不可用，并提示先重新编译或安装匹配版本。
