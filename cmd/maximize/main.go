//go:build windows

package main

import "winbinds/common"

func main() {
	if h := common.Target(); h != 0 {
		common.Maximize(h)
	}
}
