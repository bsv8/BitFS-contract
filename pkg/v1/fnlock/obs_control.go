package fnlock

import (
	"sort"
	"strings"
	"sync"
)

var (
	obsControlActionToLockIDOnce sync.Once
	obsControlActionToLockID     map[string]string
)

func obsControlActionIndex() map[string]string {
	obsControlActionToLockIDOnce.Do(func() {
		obsControlActionToLockID = buildObsControlActionToLockID()
	})
	return obsControlActionToLockID
}

func buildObsControlActionToLockID() map[string]string {
	out := map[string]string{}
	// 这里只从主框架白名单生成 obs action 映射，不扫描模块私有白名单。
	for _, item := range Whitelist {
		action := strings.TrimSpace(item.ObsControlAction)
		lockID := strings.TrimSpace(item.ID)
		if action == "" || lockID == "" {
			continue
		}
		out[action] = lockID
	}
	return out
}

// ObsControlActions 返回 obs 控制动作列表（已排序）。
func ObsControlActions() []string {
	index := obsControlActionIndex()
	out := make([]string, 0, len(index))
	for action := range index {
		out = append(out, action)
	}
	sort.Strings(out)
	return out
}

// IsObsControlActionAllowed 判断 obs 控制动作是否在白名单内。
func IsObsControlActionAllowed(action string) bool {
	_, ok := ObsControlActionLockID(action)
	return ok
}

// ObsControlActionLockID 返回动作对应的函数签名锁 ID。
func ObsControlActionLockID(action string) (string, bool) {
	lockID, ok := obsControlActionIndex()[strings.TrimSpace(action)]
	return lockID, ok
}
