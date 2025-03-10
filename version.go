package main

import (
	"fmt"
	"strings"
)

func versionCompare(version1, version2 string) int {
	parts1 := strings.Split(version1, ".")
	parts2 := strings.Split(version2, ".")

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		num1 := 0
		_, err := fmt.Sscanf(parts1[i], "%d", &num1)
		if err != nil {
			return 0
		}

		num2 := 0
		_, err = fmt.Sscanf(parts2[i], "%d", &num2)
		if err != nil {
			return 0
		}

		if num1 < num2 {
			return -1
		} else if num1 > num2 {
			return 1
		}
	}

	if len(parts1) < len(parts2) {
		return -1
	} else if len(parts1) > len(parts2) {
		return 1
	}

	return 0
}
