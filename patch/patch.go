package patch

import (
	"bytes"
	"os"
	"path/filepath"

	"capcut-toolbox/util"
)

var (
	vipEntrance = []byte{0, 'v', 'i', 'p', '_', 'e', 'n', 't', 'r', 'a', 'n', 'c', 'e', 0}
	proFortnite = []byte{0, 'p', 'r', 'o', '_', 'f', 'o', 'r', 't', 'n', 'i', 't', 'e', 0}
)

func IsPatched(exeDir string) bool {
	dllPath := filepath.Join(exeDir, "VECreator.dll")
	data, err := os.ReadFile(dllPath)
	if err != nil {
		return false
	}
	return bytes.Contains(data, proFortnite)
}

func On(exeDir string) error {
	dllPath := filepath.Join(exeDir, "VECreator.dll")
	onPath := dllPath + "_On.dll"
	offPath := dllPath + "_Off.dll"

	util.KillCapCut()

	wmDir := filepath.Join(exeDir, "Resources", "watermark")
	os.RemoveAll(wmDir)

	if _, err := os.Stat(onPath); err == nil {
		return copyFile(onPath, dllPath)
	}

	data, err := os.ReadFile(dllPath)
	if err != nil {
		return err
	}

	onData := bytes.ReplaceAll(data, vipEntrance, proFortnite)
	offData := bytes.ReplaceAll(data, proFortnite, vipEntrance)

	os.WriteFile(onPath, onData, 0644)
	os.WriteFile(offPath, offData, 0644)

	return copyFile(onPath, dllPath)
}

func Off(exeDir string) error {
	dllPath := filepath.Join(exeDir, "VECreator.dll")
	offPath := dllPath + "_Off.dll"

	util.KillCapCut()

	if _, err := os.Stat(offPath); err == nil {
		return copyFile(offPath, dllPath)
	}

	return os.ErrNotExist
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
