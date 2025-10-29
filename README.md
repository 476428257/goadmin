goadmin 是一个前后端分离的管理后台项目。
- 后端：使用 Golang 开发，基于 Gin 框架 + GORM，集成 Redis 做缓存处理，内置完善的中间件与 RBAC 权限体系。
- 前端：基于开源项目 vue-manage-system 进行二次开发，采用 Vue3 + Pinia + Element Plus + TypeScript + Vite 技术栈。

本仓库包含后端服务说明，并整合前端使用与集成指南，方便快速落地一套通用的后台管理系统。

## 功能概览

- 管理员与角色权限（RBAC）
- 菜单与路由管理（含多级菜单）
- 内容管理（文章、分类等，示例可扩展）
- 操作日志审计
- 系统配置与字典
- 登录鉴权（JWT）与令牌刷新
- 文件上传、表单校验、XSS 防护

---

## 技术栈

后端（Golang）：
- Web 框架：Gin
- ORM：GORM
- 数据库：MySQL
- 缓存：Redis
- 身份认证：JWT
- 日志：Zap
- 配置：Viper
- 验证：Validator

前端（Vue）：
- 基础：Vue 3 + TypeScript + Vite
- UI：Element Plus
- 状态：Pinia
- 路由：Vue Router
- 二开来源：vue-manage-system（感谢 lin-xin）

---

## 项目结构（后端）

```
.
├── config/             # 配置文件和配置管理
├── internal/
│   ├── controller/     # 控制器层（包含业务逻辑）
│   ├── middleware/     # 中间件集合
│   ├── model/          # 数据模型定义
│   └── router/         # 路由配置
├── pkg/
│   ├── database/       # 数据库连接管理
│   └── logger/         # 日志管理
├── sql/                # 数据库 SQL 文件
├── logs/               # 日志目录
├── config.yaml         # 主配置文件
├── go.mod
└── main.go             # 入口
```

如果前端位于单独仓库或子目录，请参考下方“前端使用”章节进行搭建与联调。

---

## 后端快速开始

1) 环境要求
- Go 1.19+（推荐最新稳定版）
- MySQL 
- Redis 

2) 拉取代码
```bash
git clone https://github.com/476428257/goadmin.git
cd ./goadmin/server
```

3) 初始化数据库,sql文件夹下有sql文件

4) 安装依赖
```bash
go mod tidy
```

5) 启动方式
- 方式 A：热重载开发
```bash
# 安装 fresh（一次性）
go install github.com/pilu/fresh@latest
# 确保 $GOPATH/bin 已加入 PATH 后，直接运行
fresh
```
- 方式 B：直接运行
```bash
go run main.go
```
- 方式 C：编译运行
```bash
go build -o goadmin-server
./goadmin-server   # Windows 在当前目录执行 goadmin-server.exe
```

服务默认监听：http://localhost:8080 （可在 config.yaml 的 server.port 调整）

## 前端使用（基于 vue-manage-system 二开）

技术栈：Vue 3 + Pinia + Element Plus + TypeScript + Vite

环境要求：
- Node.js 14.18+（Vite 3 要求）
- npm/yarn/pnpm（任选其一）

常用脚本：
```bash
# 安装依赖
npm install
# 启动开发
npm run dev
# 构建生产
npm run build
# 预览构建
npm run preview
```
页面展示：

<img width="1905" height="867" alt="1029-1" src="https://github.com/user-attachments/assets/45fc702a-2619-42d1-8416-4c4494b6c5a1" />
<img width="1904" height="867" alt="1029-2" src="https://github.com/user-attachments/assets/38fc6d64-011c-415e-9564-cedd67a30701" />
<img width="1891" height="880" alt="1029-3" src="https://github.com/user-attachments/assets/3891368e-bfa1-4eea-b292-53c4e5b0e7db" />


## 许可证

本项目采用 MIT 许可证开源，后端与前端均可自由使用、修改与分发。欢迎提交 Issue / PR 共建 goadmin！
