package config

const ConfigTemplate = `# 配置文件
# 请根据注释进行配置，不要删除任意一项

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
    bandori_station_token: "" # BanG Dream! 车站令牌
    forward_response: false # 是否转发响应
    response_content: "" # 响应内容，只有在转发响应为 true 时有效

  verify_player: # 验证玩家配置
    use_proxy: false # 是否使用代理

  backend: # 后端配置
    url: "http://tsugubot.com:8080" # 后端地址，默认为山本服务器后端地址，若有自建后端服务器可填入
    use_proxy: false # 是否使用代理

  user_data_backend: # 用户数据后端配置
    url: "http://tsugubot.com:8080" # 用户数据后端地址，默认为山本服务器后端地址，若有自建后端服务器可填入
    use_proxy: false # 是否使用代理

  functions: # 功能配置
    # 功能开关
    # 若关闭则不会处理相关功能
    help: true # 帮助文档
    car_forward: true # 车牌转发
    switch_gacha_simulate: true # 开关本群抽卡模拟
    switch_car_forward: true # 是否允许指令开启车牌转发
    bind_player: true # 绑定玩家
    change_main_server: true # 切换主服务器
    change_server_list: true # 切换服务器列表
    player_status: true # 玩家状态
    ycm: true # 有车吗
    search_player: true # 玩家信息
    search_card: true # 查卡
    card_illustration: true # 查卡面
    search_character: true # 查角色
    search_event: true # 查活动
    search_song: true # 查歌曲
    search_chart: true # 查谱面
    song_meta: true # 查询分数表
    event_stage: true # 查活动试炼
    search_gacha: true # 查卡池
    ycx: true # ycx
    ycx_all: true # ycxall
    lsycx: true # lsycx
    gacha_simulatie: true # 抽卡模拟
  
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

  car_config: # 车牌转发配置
    car: # 有效车牌关键词
      - "q1"
      - "q2"
      - "q3"
      - "q4"
      - "Q1"
      - "Q2"
      - "Q3"
      - "Q4"
      - "缺1"
      - "缺2"
      - "缺3"
      - "缺4"
      - "差1"
      - "差2"
      - "差3"
      - "差4"
      - "3火"
      - "三火"
      - "3把"
      - "三把"
      - "打满"
      - "清火"
      - "奇迹"
      - "中途"
      - "大e"
      - "大分e"
      - "exi"
      - "大分跳"
      - "大跳"
      - "大a"
      - "大s"
      - "大分a"
      - "大分s"
      - "长途"
      - "生日车"
      - "军训"
      - "禁fc"
  
    fake: # 无效车牌关键词
      - "114514"
      - "野兽"
      - "恶臭"
      - "1919"
      - "下北泽"
      - "粪"
      - "糞"
      - "臭"
      - "11451"
      - "xiabeize"
      - "雀魂"
      - "麻将"
      - "打牌"
      - "maj"
      - "麻"
      - "["
      - "]"
      - "断幺"
      - "qq.com"
      - "腾讯会议"
      - "master"
      - "疯狂星期四"
      - "离开了我们"
      - "日元"
      - "av"
      - "bv"

satori: # Satori 配置
  version: 1 # Satori 版本，目前只有 1
  path: "" # Satori 部署路径，可以为空，如果不为空需要以 / 开头
  token: "" # 鉴权令牌，如果不设置则不会进行鉴权
  host: "127.0.0.1" # 主机地址
  port: 5140 # 端口`
