# faw-vw-auto
一汽大众定时签到获取积分，适用于一汽大众APP。

## 有哪些功能
- 每日签到
- 连续签到7天，自动打开盲盒奖励；获取额外积分
- 支持Bark推送签到结果至iPhone

## 如何使用
1. 从Release里下载最新版本的`faw-vw-auto`。并解压至任意目录。
2. 配置env.yaml文件，参数如下：
- mobile: 手机号(一汽大众APP注册手机号)
- password: 密码
- securityCode: 安全码（需要先用抓包APP登录的接口https://oneapp-api.faw-vw.com/account/login/loginByPassword/v1，在Request Body里就有这个参数）
- did：设备ID（需要先用抓包APP的接口，就可以看到这个参数）
- BarkPushServerURL：Bark推送地址（可选），这个参数是使用Bark来推送签到结果至iPhone，如果不需要推送可以不填写。
3. 运行`faw-vw-auto`， 执行命令如‘nohup ./faw-vw-auto-darwin &’，即可自动签到。

## 注意事项
- 本项目仅供学习交流使用，请勿用于商业用途。
- 因涉及账户密码，请勿将env.yaml文件上传至公共仓库，以免泄露个人信息。