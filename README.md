# faw-vw-auto
一汽大众定时签到获取积分，Release版本只编译了macOS（M系列芯片）环境的程序包，其他平台可自行下载编译。
<div>
    <img src="IMG_1.PNG" width="40%" alt="预览1">
    <img src="IMG_2.PNG" width="40%" alt="预览2">
</div>

## 有哪些功能
- 每日签到
- 连续签到7天，自动打开盲盒奖励；获取额外积分
- 支持Bark推送签到结果至iPhone

## 如何使用
#### 配置文件
无论哪种平台，都需要配置一个`env.yaml`文件来存储账号信息和签到参数。
配置env.yaml文件，参数如下：
- mobile: 手机号（一汽大众APP注册手机号）
- password: 密码密文（需要先用浏览器登录一下Web版的应用，然后看下https://vw.faw-vw.com/api/business/cpoint/registeOrLogin。在Request Body里就有这个参数）
- WebDid：设备ID（同上，在Request Body里的did参数）
- BarkPushServerURL：Bark推送地址（可选），这个参数是使用Bark来推送签到结果至iPhone，如果不需要推送可以不填写。

- securityCode：客户端签到接口必须参数。这个参数可以通过抓包App的登录接口（https://oneapp-api.faw-vw.com/account/login/loginByPassword/v1）获取。
- did：设备ID。这个参数是客户端应用的设备ID，可以通过抓包App的登录接口（https://oneapp-api.faw-vw.com/account/login/loginByPassword/v1）获取。

### 环境要求
- macOS（M系列芯片）环境，Release版本已编译。
1. 从Release里下载最新版本的`faw-vw-auto`。并解压至任意目录，并填写配置文件env.yaml里的配置项。
2. 运行`faw-vw-auto`， 执行命令如‘nohup ./faw-vw-auto-darwin &’，即可自动签到。

- Docker环境
1. 确保Docker已安装并运行。
2. 下载Docker镜像：`docker pull hytaoist/faw-vw-auto:latest`
3. 运行Docker容器：`docker run -d --name faw-vw-auto -v /path/to/your/env.yaml:/app/env.yaml hytaoist/faw-vw-auto:latest`
   - 注意替换`/path/to/your/env.yaml`为你的env.yaml文件的实际路径。

- 其他平台（Linux、Windows等）可自行下载源代码编译。

## 注意事项
- 本项目仅供学习交流使用，请勿用于商业用途。
- 因涉及账户密码，请勿将env.yaml文件上传至公共仓库，以免泄露个人信息。