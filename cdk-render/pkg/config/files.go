package config

import (
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
)

func ReadData(dirPath string) (map[string]string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		logrus.Error("Failed to read dir")
		return nil, err
	}

	resp := map[string]string{}
	for _, file := range files {
		if !file.IsDir() {
			filePath := path.Join(dirPath, file.Name())

			fInfo, err := os.Lstat(filePath)
			if err != nil {
				return nil, err
			}

			if fInfo.Mode()&os.ModeSymlink != 0 {
				fp, err := os.Readlink(filePath)
				if err != nil {
					return nil, err
				}
				filePath = path.Join(dirPath, fp)
				fInfo, err = os.Stat(filePath)
				if err != nil {
					return nil, err
				}
			}

			var content []byte
			if !fInfo.IsDir() {
				logrus.Infof("Reading file %s", file.Name())
				content, err = os.ReadFile(filePath)
				if err != nil {
					return nil, err
				}
			}

			resp[file.Name()] = strings.TrimSuffix(string(content), "\n")
		}
	}

	return resp, nil
}
