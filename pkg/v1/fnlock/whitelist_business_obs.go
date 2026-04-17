package fnlock

// 观测日志能力入口。
var whitelistBusinessObs = []LockedFunction{
	{
		ID:        "bftp.obs.info",
		Module:    ModuleBFTP,
		Package:   "./pkg/obs",
		Symbol:    "Info",
		Signature: "func Info(service, name string, fields map[string]any)",
		Note:      "系统流程日志统一入口。",
	},
	{
		ID:        "bftp.obs.debug",
		Module:    ModuleBFTP,
		Package:   "./pkg/obs",
		Symbol:    "Debug",
		Signature: "func Debug(service, name string, fields map[string]any)",
		Note:      "调试日志统一入口。",
	},
	{
		ID:        "bftp.obs.business",
		Module:    ModuleBFTP,
		Package:   "./pkg/obs",
		Symbol:    "Business",
		Signature: "func Business(service, name string, fields map[string]any)",
		Note:      "业务日志统一入口。",
	},
	{
		ID:        "bftp.obs.important",
		Module:    ModuleBFTP,
		Package:   "./pkg/obs",
		Symbol:    "Important",
		Signature: "func Important(service, name string, fields map[string]any)",
		Note:      "关键业务日志统一入口。",
	},
	{
		ID:        "bftp.obs.error",
		Module:    ModuleBFTP,
		Package:   "./pkg/obs",
		Symbol:    "Error",
		Signature: "func Error(service, name string, fields map[string]any)",
		Note:      "错误日志统一入口。",
	},
}
