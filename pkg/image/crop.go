package image

import (
	"fmt"
	"log"

	"ha-video-parser/pkg/utils"
)

func BatchCrop(input, output string) bool {

	result := utils.Exec("magick", "mogrify", "-crop", "640x480+0+420", "-path", output, fmt.Sprintf("%s/*", input))
	if result.ExitCode != 0 {
		log.Printf("ImageMagick failed to crop frames for [%s] with code [%d] and stderr [%s]", input, result.ExitCode, result.StdErr)
		return false
	}

	return true
}
