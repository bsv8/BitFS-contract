package fnlock

// sqlite actor 能力边界入口。
var whitelistBusinessSQLiteActor = []LockedFunction{
	{
		ID:        "bftp.sqliteactor.open",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "Open",
		Signature: "func Open(path string, debug bool) (*Opened, error)",
		Note:      "sqlite 单 owner actor 的标准打开入口。",
	},
	{
		ID:        "bftp.sqliteactor.new",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "New",
		Signature: "func New(db *sql.DB) (*Actor, error)",
		Note:      "actor 构造入口，串行访问边界依赖该签名。",
	},
	{
		ID:        "bftp.sqliteactor.actor_do",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "Actor.Do",
		Signature: "func (a *Actor) Do(ctx context.Context, fn func(*sql.DB) error) error",
		Note:      "sqlite 非事务执行主入口。",
	},
	{
		ID:        "bftp.sqliteactor.actor_tx",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "Actor.Tx",
		Signature: "func (a *Actor) Tx(ctx context.Context, fn func(*sql.Tx) error) error",
		Note:      "sqlite 事务执行主入口。",
	},
	{
		ID:        "bftp.sqliteactor.do_value",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "DoValue",
		Signature: "func DoValue[T any](ctx context.Context, a *Actor, fn func(*sql.DB) (T, error)) (T, error)",
		Note:      "sqlite 非事务泛型返回值入口。",
	},
	{
		ID:        "bftp.sqliteactor.tx_value",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "TxValue",
		Signature: "func TxValue[T any](ctx context.Context, a *Actor, fn func(*sql.Tx) (T, error)) (T, error)",
		Note:      "sqlite 事务泛型返回值入口。",
	},
}
