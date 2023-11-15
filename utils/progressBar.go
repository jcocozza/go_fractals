package utils

import (
	"fmt"
	"strings"
)

func ProgressBar(current, total int) {
	const barWidth = 50
	progress := int(float32(current) / float32(total) * float32(barWidth))

	fmt.Printf("\rProgress: [%s%s] %d%%", strings.Repeat("=", progress), strings.Repeat(" ", barWidth-progress), progress*2)
}