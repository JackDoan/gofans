package hwmonHelp

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func FindDeviceByName(name string) (string, error) {
	baseDir := "/sys/class/hwmon"
	hwmons, _ := os.ReadDir(baseDir)
	for _, d := range hwmons {
		h := filepath.Join(baseDir, d.Name())
		n, err := os.ReadFile(filepath.Join(h, "name"))
		if err != nil {
			continue
		}
		if strings.TrimSpace(string(n)) == name {
			return h, nil
		}
	}
	return "", errors.New("device not found")
}

// OfNameAndLabel finds a sensor given hwmon name and label.
//
// For example, if you have /sys/class/hwmon/hwmon4/name containing
// "k10temp", and /sys/class/hwmon/hwmon4/temp1_label containing
// "Tctl", the former would be name, and the latter would be label.
func OfNameAndLabel(name string, label string) string {
	baseDir := "/sys/class/hwmon"
	files, _ := os.ReadDir(baseDir)
	for _, file := range files {
		n, _ := os.ReadFile(filepath.Join(baseDir, file.Name(), "name"))
		if strings.TrimSpace(string(n)) == name {
			baseDir := filepath.Join(baseDir, file.Name())
			files, _ := os.ReadDir(filepath.Join(baseDir))
			for _, file := range files {
				if strings.HasSuffix(file.Name(), "_label") {
					l, _ := os.ReadFile(filepath.Join(baseDir, file.Name()))
					if strings.TrimSpace(string(l)) == label {
						filename := file.Name()
						filename = strings.TrimSuffix(filename, "_label") + "_input"
						return filepath.Join(baseDir, filename)
					}
				}
			}
		}
	}
	return ""
}

func OfNameAndInput(name string, label string) string {
	baseDir := "/sys/class/hwmon"
	files, _ := os.ReadDir(baseDir)
	for _, file := range files {
		n, _ := os.ReadFile(filepath.Join(baseDir, file.Name(), "name"))
		if strings.TrimSpace(string(n)) == name {
			baseDir := filepath.Join(baseDir, file.Name())
			files, _ := os.ReadDir(filepath.Join(baseDir))
			for _, file := range files {
				if strings.HasSuffix(file.Name(), "_input") {
					l := strings.TrimSuffix(file.Name(), "_input")
					if strings.TrimSpace(l) == label {
						return filepath.Join(baseDir, file.Name())
					}
				}
			}
		}
	}
	return ""
}
