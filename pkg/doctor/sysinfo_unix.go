//go:build !windows

package doctor

import "syscall"

func getDiskSpace(path string) (uint64, uint64) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return 0, 0
	}
	free := stat.Bavail * uint64(stat.Bsize)
	total := stat.Blocks * uint64(stat.Bsize)
	return free, total
}
