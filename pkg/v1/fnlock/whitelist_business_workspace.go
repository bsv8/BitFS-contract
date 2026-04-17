package fnlock

// 工作区业务能力入口。
var whitelistBusinessWorkspace = []LockedFunction{
	{
		ID:               "bitfs.clientapp.workspace.kernel_list",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "workspaceManager.List",
		Signature:        "func (m *workspaceManager) List() ([]workspaceItem, error)",
		ObsControlAction: "workspace.list",
		Note:             "workspace 列表能力单点入口，调用侧统一从 workspace manager 获取结果。",
	},
	{
		ID:               "bitfs.clientapp.workspace.kernel_add",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "workspaceManager.Add",
		Signature:        "func (m *workspaceManager) Add(path string, maxBytes uint64) (workspaceItem, error)",
		ObsControlAction: "workspace.add",
		Note:             "workspace 新增能力单点入口。",
	},
	{
		ID:               "bitfs.clientapp.workspace.kernel_update",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "workspaceManager.UpdateByPath",
		Signature:        "func (m *workspaceManager) UpdateByPath(workspacePath string, maxBytes *uint64, enabled *bool) (workspaceItem, error)",
		ObsControlAction: "workspace.update",
		Note:             "workspace 更新能力单点入口。",
	},
	{
		ID:               "bitfs.clientapp.workspace.kernel_delete",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "workspaceManager.DeleteByPath",
		Signature:        "func (m *workspaceManager) DeleteByPath(workspacePath string) error",
		ObsControlAction: "workspace.delete",
		Note:             "workspace 删除能力单点入口。",
	},
	{
		ID:               "bitfs.clientapp.workspace.kernel_sync_once",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "workspaceManager.SyncOnce",
		Signature:        "func (m *workspaceManager) SyncOnce(ctx context.Context) (map[string]sellerSeed, error)",
		ObsControlAction: "workspace.sync_once",
		Note:             "workspace 单次扫描同步能力单点入口。",
	},
}
