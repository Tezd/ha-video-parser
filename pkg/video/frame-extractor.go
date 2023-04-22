package video

import (
	"fmt"
	"log"

	"ha-video-parser/pkg/utils"
)

func ExtractFrames(input, output string) bool {
	result := utils.Exec("ffmpeg", "-i", input, fmt.Sprintf("%s/%%09d.png", output))
	if result.ExitCode != 0 {
		log.Printf("FFMPEG failed to parse frames for [%s] with code [%d] and stderr [%s]", input, result.ExitCode, result.StdErr)
		return false
	}

	return true
}
