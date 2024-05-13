<div align="center">

![tsugu-bot-go logo](https://github.com/WindowsSov8forUs/tsugu-bot-go/blob/main/logo/tsugu-bot-go.png)

# tsugu-bot-go

_✨ 一个 Go 实现的 [TsuguBanGDreamBot](https://github.com/Yamamoto-2/tsugu-bangdream-bot) 前端整合机器人应用 ✨_

</div>

<p align="center">

<a href="https://github.com/Yamamoto-2/tsugu-bangdream-bot">
  <img src="https://img.shields.io/badge/tsugu bangdream bot-v2 api-FFEE88" alt="license">
</a>

<a href="https://github.com/WindowsSov8forUs/tsugu-bot-go">
  <img src="https://img.shields.io/github/v/release/WindowsSov8forUs/tsugu-bot-go" alt="Latest Release Version">
</a>

<a href="https://github.com/WindowsSov8forUs/tsugu-bot-go/blob/main/LICENSE">
  <img src="https://img.shields.io/github/license/WindowsSov8forUs/tsugu-bot-go" alt="License">
</a>

<a href="https://golang.org/dl/">
  <img src="https://img.shields.io/github/go-mod/go-version/WindowsSov8forUs/tsugu-bot-go" alt="Go Version">
</a>

</p>

## 引用

本项目参考了这些项目进行编写

- [`tsugu-python-frontend`](https://github.com/kumoSleeping/tsugu-python-frontend)

本项目引用了如下项目

- [`satori-protocol-go/satori-model-go`](https://github.com/satori-protocol-go/satori-model-go)

## 说明

本项目是一个 Go 语言编写的机器人应用，通过与聊天平台进行通讯交互，并向 Tsugu 官方处理后端和用户数据后端、自建 Tsugu 后端或本地数据库发送请求进行数据处理，实现几乎全部的 **[TsuguBanGDreamBot](https://github.com/Yamamoto-2/tsugu-bangdream-bot)** (以下简称为 `Tsugu` ) 功能。

> 目前仅实现了与 Satori 协议聊天平台进行交互，且由于历史遗留问题暂时并未实现平台间冲突隔离，请注意

### 实现的功能

目前 **tsugu-bot-go** 已经实现了 Tsugu 的几乎全部功能，且支持对所有功能的全局开关配置，同时添加了针对机器人消息处理针对性的配置，以处理当使用场景出现冲突时的防冲突处理。

### 待实现的功能

- [ ] 多群组的机器人总开关
- [ ] 运行过程中的功能开关配置
- [x] 指令的别名配置

## 使用

从最新的 [Release](https://github.com/WindowsSov8forUs/tsugu-bot-go/releases) 中选择适合自己的版本，下载并运行。

这个过程中，建议为 tsugu-bot-go 单独创建一个文件夹。

> Windows 系统在运行 `exe` 文件后，将会生成一个 `bat` 脚本，随后直接双击脚本运行即可。

> 目前仅 Windows 64位系统经过测试，其他系统若运行出现问题请立马 [告知作者](https://github.com/WindowsSov8forUs/tsugu-bot-go/issues)

成功运行后，初次运行时 tsugu-bot-go 将会生成一个配置文件 `config.yml` 并退出运行。

修改配置文件后重新运行即可。若配置无误，将会在终端看到 “连接成功” 的输出。

### 支持的聊天平台

- [x] Satori 协议聊天平台 ([`chronocat`](https://github.com/chrononeko/chronocat) 等)
- [ ] Lagrange ([`LagrangeGo`](https://github.com/LagrangeDev/LagrangeGo) 正在开发中...)
- [ ] 其他可能支持的聊天平台...

## 配置

> 目前对于聊天平台连接的配置仅有 Satori 协议配置

当 tsugu-bot-go 成功初次运行后，将会在同级目录下生成一个 `config.yml` 文件。编辑该文件即可以进行配置。

配置结构如下：

```yaml
# 日志等级
# 可选项：
#   - 0：关闭日志
#   - 1：仅输出致命错误日志
#   - 2：输出致命错误日志和错误日志
#   - 3：输出致命错误日志、错误日志和警告日志
#   - 4：输出致命错误日志、错误日志、警告日志和信息日志
#   - 5：输出致命错误日志、错误日志、警告日志、信息日志和调试日志
#   - 6/7：输出所有日志
log_level: 4

tsugu: # Tsugu 机器人配置
  ...

satori: # Satori 配置
  ...
```

### Tsugu 机器人配置

该结构下包含多个针对 Tsugu 表现的配置，绝大多数配置无需特意更改，保持默认值即可。

配置结构如下：

```yaml
tsugu: # Tsugu 机器人配置
  require_at: false # 是否需要 at 才能触发指令
  reply: true # 是否回复消息
  at: false # 是否 at 发送者
  no_space: false # 是否不需要指令头后空格
  timeout: 10 # 超时时间，单位为秒
  proxy: "" # 代理地址，如果不使用代理请留空
  use_easy_bg: true # 是否使用简单背景，若关闭可能会减慢处理速度
  compress: true # 是否启用压缩，若关闭可能会增加流量消耗

  # 用户数据路径
  # 若不使用，将会使用远程用户数据库
  # 若不使用请设置为空
  user_database_path: ""

  # 禁止模拟抽卡
  # 填入禁止使用抽卡模拟的群号
  ban_gacha_simulate: []

  car_station: # 车站配置
    ...

  verify_player: # 验证玩家配置
    ...

  backend: # 后端配置
    ...

  user_data_backend: # 用户数据后端配置
    ...

  functions: # 功能配置
    # 功能开关
    # 若关闭则不会处理相关功能
    ...

  car_config: # 车牌转发配置
    ...
```

- `require_at`: 指定是否需要 at 才能触发指令；当为 `true` 时，仅 at 机器人的消息可以触发指令。
- `reply`: 指定发送消息时是否回复对应的消息。
- `at`: 指定发送消息时是否 @ 对应的发送者，在部分平台中回复并不总是携带 @ 。
- `no_space`: 指定是否启用无需空格触发大部分指令，启用这将方便一些用户使用习惯，但会增加 bot 误判概率，仍然建议使用空格。
- `timeout`: 指定后端请求超时时间；当为 `0` 时，将长时间等待后端响应直到连接关闭。建议设置范围在 `15-30` 之间。
- `proxy`: 代理服务器地址，所有的请求当需要经过代理服务器时都会经过该服务器地址。当需要通过代理服务器才能够访问后端服务器时，通过该项配置配合子配置项中的 `use_proxy` 进行配置。
- `use_easy_bg`: 是否使用简单背景。使用简单背景可以一定程度上减少后端响应耗时，若对生成图片没有特殊要求建议开启。
- `compress`: 是否压缩图片。压缩图片可以一定程度上减少响应传输所用流量以及图片发送所需耗时，若对生成图片的质量没有特殊要求建议开启。
- `user_database_path`: 本地用户数据库路径；当为空时将使用 `user_data_backend` 中配置的用户数据库后端。需要注意的是，各用户数据库以及本地用户数据库之间不会进行数据统一，请酌情配置。
- `ban_gacha_simulate`: 禁止使用 `抽卡模拟` 功能的群组，填入群组 ID 以配置。可以在运行过程中配置，但无法保存至文件中，且在文件中的配置更改不会立马反馈到应用中，因此每次更改都需要重启 tsugu-bot-go 以应用配置。

<details>
<summary>车站配置 car_station</summary>

#### 车站配置 `car_station`

`car_station` 子配置项用于对车牌转发进行配置。需要注意的是，车牌转发依然受 `require_at` 控制，因此当 `require_at` 为 `true` 时，仍然需要 at 机器人才能够进行车牌转发。此时建议开启 `forward_response` 选项，以即时得知车牌是否转发成功。

配置结构如下：

```yaml
  car_station: # 车站配置
    bandori_station_token: "" # BanG Dream! 车站令牌
    forward_response: false # 是否转发响应
    response_content: "" # 响应内容，只有在转发响应为 true 时有效
```

- `bandori_station_token`: 车站转发所需令牌。若没有自己的令牌，可不填，将会默认使用 Tsugu 的令牌。
- `forward_response`: 是否在转发成功后进行响应。若为 `true` 则将会在转发成功后回复车牌所在消息，否则将会保持静默。
- `response_content`: 仅当 `forward_response` 为 `true` 时有效，可用于自定义转发成功后的回复消息。此时若此配置留空，则会回复默认的转发成功消息。

</details>

<details>
<summary>验证玩家配置 verify_player</summary>

#### 验证玩家配置 `verify_player`

`verify_player` 子配置项用于对验证玩家进行配置。目前只有 `use_proxy` 唯一一个配置项，用于控制当使用本地用户数据库时是否使用代理服务器。

由于当使用本地用户数据库时，验证玩家需要访问 **[bestdori](https://bestdori.com/)** ，因此若您的机器人所在网络环境无法直接访问 bestdori 或访问不稳定，且需要使用本地用户数据库，建议启用该配置。

配置结构如下：

```yaml
  verify_player: # 验证玩家配置
    use_proxy: false # 是否使用代理
```

</details>

<details>
<summary>后端配置 backend</summary>

#### 后端配置 `backend`

`backend` 子配置项用于对后端进行配置。一般不需要更改该子配置。

配置结构如下：

```yaml
  backend: # 后端配置
    url: "http://tsugubot.com:8080" # 后端地址，默认为山本服务器后端地址，若有自建后端服务器可填入
    use_proxy: false # 是否使用代理
```

- `url`: 后端地址。默认为 Tsugu 官方后端地址，若有自建后端服务器且官方后端访问不稳定可以配置。
- `use_proxy`: 是否使用代理。若机器人所在网络环境访问后端受限，可启用该配置。

</details>

<details>
<summary>用户数据后端配置 user_data_backend</summary>

#### 用户数据后端配置 `user_data_backend`

`user_data_backend` 子配置项用于对用户数据后端进行配置。一般不需要更改该子配置。若配置了 `user_database_path` 选项，则表明启用了本地用户数据库，该子配置将无效。

配置结构如下：

```yaml
  user_data_backend: # 用户数据后端配置
    url: "http://tsugubot.com:8080" # 用户数据后端地址，默认为山本服务器后端地址，若有自建后端服务器可填入
    use_proxy: false # 是否使用代理
```

- `url`: 后端地址。默认为 Tsugu 官方后端地址，若有自建后端服务器且官方后端访问不稳定可以配置。
- `use_proxy`: 是否使用代理。若机器人所在网络环境访问后端受限，可启用该配置。

</details>

<details>
<summary>功能启用配置 functions</summary>

#### 功能启用配置 `functions`

`functions` 子配置项用于对 tsugu-bot-go 所启用的功能进行配置。其中每个子配置项都是一个功能，设置为 `true` 则在全局启用该功能，否则在全局关闭该功能。默认为全部开启，一般不需要进行更改。

配置结构如下：

```yaml
  functions: # 功能配置
    # 功能开关
    # 若关闭则不会处理相关功能
    help: true # 帮助文档
    car_forward: true # 车牌转发
    change_main_server: true # 切换主服务器
    switch_car_forward: true # 是否允许指令开启车牌转发
    bind_player: true # 绑定玩家
    change_server_list: true # 切换服务器列表
    player_status: true # 玩家状态
    card_illustration: true # 查卡面
    player: true # 玩家信息
    gacha_simulatie: true # 抽卡模拟
    gacha: true # 查卡池
    event: true # 查活动
    song: true # 查歌曲
    song_meta: true # 查询分数表
    character: true # 查角色
    chart: true # 查谱面
    ycx: true # ycx
    ycx_all: true # ycxall
    lsycx: true # lsycx
    ycm: true # 有车吗
    card: true # 查卡
```

- `help`: 帮助消息，发送指定功能的帮助信息。
- `car_forward`: 车牌转发。
- `change_main_server`: 切换主服务器。
- `switch_car_forward`: 开启/关闭个人车牌转发。
- `bind_player`: 绑定玩家。
- `change_server_list`: 设置服务器列表。
- `player_status`: 玩家状态。
- `card_illustration`: 查卡面。
- `player`: 查询玩家信息。
- `gacha_simulatie`: 抽卡模拟。
- `gacha`: 查卡池。
- `event`: 查活动。
- `song`: 查歌曲。
- `song_meta`: 查分数表。
- `character`: 查角色。
- `chart`: 查谱面。
- `ycx`: 活动指定等级预测线。
- `ycx_all`: 活动全部预测线。
- `lsycx`: 活动的历史预测线。
- `ycm`: 有车吗？/查询车站车牌号。
- `card`: 查卡。

</details>

<details>
<summary>功能启用配置 functions</summary>

#### 指令别名配置 `command_alias`

`command_alias` 子配置项用于对 tsugu-bot-go 的指令进行别名配置。其中每个子配置项都是一个字符串列表。

配置结构如下：

```yaml
  command_alias: # 指令别名
    switch_gacha_simulate: [] # 开关本群抽卡模拟
    open_car_forward: [] # 开启车牌转发
    close_car_forward: [] # 关闭车牌转发
    bind_player: [] # 绑定玩家
    unbind_player: [解绑玩家] # 解绑玩家
    change_main_server: [服务器模式, 切换服务器] # 切换主服务器
    change_server_list: [默认服务器] # 设置默认服务器
    player_status: [] # 玩家状态
    ycm: [有车吗, 车来] # 有车吗
    search_player: [查询玩家] # 玩家信息
    search_card: [查卡牌] # 查卡
    card_illustration: [查卡插画, 查插画] # 查卡面
    search_character: [] # 查角色
    search_event: [] # 查活动
    search_song: [] # 查歌曲
    search_chart: [] # 查谱面
    song_meta: [查分数表, 查询分数榜, 查分数榜] # 查询分数表
    event_stage: [查stage, 查舞台, 查festival, 查5v5] # 查活动试炼
    search_gacha: [] # 查卡池
    ycx: [] # ycx
    ycx_all: [myycx] # ycxall
    lsycx: [] # lsycx
    gacha_simulatie: [] # 抽卡模拟
```

</details>

<details>
<summary>车牌转发配置 car_config</summary>

#### 车牌转发配置 `car_config`

`car_config` 子配置项用于对车牌转发功能进行配置。两个子配置项均为字符串数组。

该子配置项自带一定数量默认项，一般情况下不需要进行额外配置。

配置结构如下：

```yaml
  car_config: # 车牌转发配置
    car: # 有效车牌关键词
      ...
  
    fake: # 无效车牌关键词
      ...
```

- `car`: 有效车牌关键词。若消息中不含该数组内任一关键词则将被视为无效车牌，不予转发。
- `fake`: 无效车牌关键词。若消息中含有该数组内任一关键词则将被视为无效车牌，不予转发。

</details>

### 平台连接配置

在配置文件中，剩余的配置结构都将被视为平台连接配置。不同的配置结构名将被应用在不同的平台连接配置中。

若某个平台的配置项不全，则将视为不连接对应平台。

> 目前仅适配了 Satori 协议聊天平台

<details>
<summary>Satori 协议聊天平台</summary>

#### Satori 协议聊天平台

`satori` 配置项将用于对 **[Satori 协议](https://satori.js.org/zh-CN/)** 聊天平台进行配置。具体的配置内容请从对应的 Satori 平台内获得。

配置结构如下：

```yaml
satori: # Satori 配置
  version: 1 # Satori 版本，目前只有 1
  path: "" # Satori 部署路径，可以为空，如果不为空需要以 / 开头
  token: "" # 鉴权令牌，如果不设置则不会进行鉴权
  host: "http://127.0.0.1" # 主机地址
  port: 5140 # 端口
```

- `version`: Satori 协议版本号，输入对应版本号的数字即可。
    > 目前仅支持 `v1` 版本
- `path`: Satori 平台设置的机器人部署路径，若未设置则留空即可。
- `token`: Satori 平台与机器人应用连接所需的鉴权令牌，若无需鉴权则留空即可。若需要鉴权，留空则无法进行连接。
- `host`: Satori 平台与机器人应用连接的地址，若处于相同网络环境下则填入 `"http://127.0.0.1"` 即可。
- `port`: Satori 平台与机器人应用连接的端口，输入端口对应 `1-65535` 的数字即可。

</details>
