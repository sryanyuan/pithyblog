package blog

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ngaut/log"
	uuid "github.com/satori/go.uuid"
)

func packSitedata(srcdir, todir string) (string, error) {
	u := uuid.NewV4()
	packPath := todir
	packFilename := "sitedata_" + u.String() + ".zip"
	packFullPath := filepath.Join(packPath, packFilename)
	err := os.MkdirAll(packPath, 0777)
	if nil != err {
		return "", err
	}

	// Call zip to pack
	cmd := exec.Command("zip", "-r", packFullPath, srcdir)
	out, err := cmd.CombinedOutput()
	if nil != err {
		return "", err
	}

	log.Info(string(out))
	return packFullPath, nil
}
