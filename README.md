# BitFS-contract

BitFS 跨项目共享契约层（只放协议和数据库结构事实，不放运行时实现）。

## 目录

- `proto/`: 跨进程协议契约（.proto）
- `ent/v1/schema/`: 当前数据库契约真源（v1）
- `ent/v1/gen/`: 由 v1 schema 生成的 Go 代码
- `ent/v2/`: 预留给下一次大版本
- `pkg/`: 轻量 Go 契约类型（错误码、协议 ID、DB 边界策略）
- `scripts/`: 生成与校验脚本

## 约束

- 系统内唯一 ID 语义：`pubkey_hex`
- 错误信息用英文，注释与说明用中文
- 本仓库只承载契约，不承载业务实现
