package lock

import (
	"os"
	"os/exec"
	"path/filepath"

	"capcut-toolbox/util"
)

func IsLocked(appDir string) bool {
	pi := filepath.Join(appDir, "ProductInfo.xml")
	_, err := os.Stat(pi)
	return err == nil
}

func On(appDir string) error {
	util.KillCapCut()

	// Set version in configure.ini
	ini := filepath.Join(appDir, "configure.ini")
	f, err := os.OpenFile(ini, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	f.WriteString("[Version]\nlast_version=1.0.0.0\n")
	f.Close()

	// Delete User Data folder
	os.RemoveAll(filepath.Join(appDir, "User Data"))

	// Create readonly lock file
	createReadonly(filepath.Join(appDir, "ProductInfo.xml"))
	return nil
}

func Off(appDir string) error {
	util.KillCapCut()

	removeReadonly(filepath.Join(appDir, "ProductInfo.xml"))
	os.Remove(filepath.Join(appDir, "ProductInfo.xml"))
	return nil
}

func createReadonly(path string) {
	os.MkdirAll(filepath.Dir(path), 0755)
	os.WriteFile(path, []byte{}, 0644)
	exec.Command("attrib", "+r", path).Run()
}

func removeReadonly(path string) {
	exec.Command("attrib", "-r", path).Run()
}
