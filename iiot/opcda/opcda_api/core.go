package opcda_api

import (
	"sync"

	"github.com/huskar-t/opcda"
)

var (
	GServerMap sync.Map // name=>*ServerItem
	GGroup     = "group01"
)

type ServerItem struct {
	Host   string
	Prog   string
	Server *opcda.OPCServer
	Group  *opcda.OPCGroup
	Handle []uint32
	Tag    []string
}

func recursiveBrowse(browser *opcda.OPCBrowser, results *[]string, filterSystem bool) {
	// A. 提取并过滤叶子节点 (Tags)
	browser.ShowLeafs(false)
	leafCount := browser.GetCount()
	for i := 0; i < leafCount; i++ {
		shortName, _ := browser.Item(i)

		// 过滤逻辑：如果开启过滤且名称以 _ 开头，则跳过
		if filterSystem && len(shortName) > 0 && shortName[0] == '_' {
			continue
		}

		// 只有通过过滤的才获取 ItemID
		fullID, _ := browser.GetItemID(shortName)
		if fullID != "" {
			*results = append(*results, fullID)
		}
	}

	// B. 提取并过滤分支 (Folders)
	browser.ShowBranches()
	branchCount := browser.GetCount()
	branches := make([]string, 0, branchCount) // 预分配容量但初始长度为0
	for i := 0; i < branchCount; i++ {
		bName, _ := browser.Item(i)

		// 同样过滤文件夹：如果文件夹以 _ 开头，其内部的点位通常也是系统点位，直接切断递归
		if filterSystem && len(bName) > 0 && bName[0] == '_' {
			continue
		}
		branches = append(branches, bName)
	}

	// C. 递归进入分支
	for _, bName := range branches {
		if err := browser.MoveDown(bName); err == nil {
			// 传递过滤参数到下一层
			recursiveBrowse(browser, results, filterSystem)
			browser.MoveUp()
		}
	}
}
