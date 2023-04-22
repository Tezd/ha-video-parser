package image

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"ha-video-parser/pkg/utils"
)

const (
	MAX_RGB               = "rgb(255,255,255)"
	THRESHOLD_RGB_BLACK   = "rgb(3,3,2)"
	THRESHOLD_RGB_WHITE   = "rgb(3,3,3)"
	MIN_RGB               = "rgb(0,0,0)"
	MIN_DIFFERENCE        = 0.001
	BLACK_WHITE_THRESHOLD = 5
)

func NonEmptyUnique(input string, output string) []string {

	textFrames, err := pickFramesWithTextOnly(input, output)
	if err != nil {
		return nil
	}

	sort.Strings(textFrames)

	seen := make(map[uint64]string)
	var uniques []string

	for _, frame := range textFrames {
		fileName := filepath.Base(frame)
		frameId, err := strconv.ParseUint(fileName[:len(fileName)-len(filepath.Ext(fileName))], 10, 64)
		if err != nil {
			log.Printf("Failed to parse file name [%s] with error [%s]", fileName, err.Error())
			return nil
		}

		seen[frameId] = frame

		//there was no previous frame
		if _, ok := seen[frameId-1]; !ok {
			uniques = append(uniques, frame)
			continue
		}

		maskFile, err := createMaskFromFrame(seen[frameId-1], output)
		if err != nil {
			log.Printf("Failed to create mask file [%s] with error [%s]", seen[frameId-1], err.Error())
			return nil
		}

		if areSimilar(output, maskFile, frame) {
			continue
		}

		uniques = append(uniques, frame)
	}

	return uniques
}

func areSimilar(output, mask, frame string) bool {

	combinations := fmt.Sprintf("%s/combinations", output)
	if err := os.MkdirAll(combinations, os.ModePerm); err != nil {
		log.Printf("Failed to create directory for combinations [%s]", err.Error())
		return false
	}

	combination := fmt.Sprintf("%s/%s", combinations, filepath.Base(frame))
	r := utils.Exec("magick", "convert", frame, mask, "-composite", combination)

	if r.ExitCode != 0 {
		log.Printf("Failed to combine images [%s]", r.StdErr)
		return false
	}

	return isAllBlack(combination)
}

func createMaskFromFrame(frame, output string) (string, error) {
	masks := fmt.Sprintf("%s/masks", output)
	if err := os.MkdirAll(masks, os.ModePerm); err != nil {
		log.Printf("Failed to create directory for masks [%s]", err.Error())
		return "", err
	}

	frameFile := filepath.Base(frame)

	maskFile := fmt.Sprintf("%s/%s", masks, frameFile)

	r := utils.Exec(
		"magick",
		"convert",
		frame,
		"-threshold",
		fmt.Sprintf("%d%%", BLACK_WHITE_THRESHOLD),
		"-transparent",
		MIN_RGB,
		"-fill",
		MIN_RGB,
		"-opaque",
		MAX_RGB,
		maskFile,
	)

	if r.ExitCode != 0 {
		return "", errors.New(r.StdErr)
	}

	return maskFile, nil
}

func GetMaxRgb() {
	//magick identify -colorspace RGB -format '%[fx:round(255 * maxima.r)] %[fx:round(255 * maxima.g)] %[fx:round(255 * maxima.b)]' 000000780.png
}

func pickFramesWithTextOnly(input, output string) ([]string, error) {
	texts := fmt.Sprintf("%s/texts", output)
	if err := os.MkdirAll(texts, os.ModePerm); err != nil {
		log.Printf("Failed to create director for uniques [%s]", err.Error())
		return nil, err
	}

	crops, err := os.ReadDir(input)
	if err != nil {
		log.Printf("Failed to list director [%s] with error [%s]", input, err.Error())
		return nil, err
	}

	var filePaths []string
	for _, crop := range crops {
		if crop.IsDir() {
			continue
		}

		filePath := fmt.Sprintf("%s/%s", input, crop.Name())
		if isAllBlack(filePath) {
			continue
		}

		filePaths = append(filePaths, filePath)
	}

	for _, filePath := range filePaths {
		f, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed opening source [%s] with error [%s]", f.Name(), err.Error())
			continue
		}

		outpath := fmt.Sprintf("%s/%s", texts, filepath.Base(filePath))
		o, err := os.Create(outpath)
		if err != nil {
			log.Printf("Failed to create destination [%s] with error [%s]", outpath, err.Error())
			continue
		}
		defer o.Close()
		if _, err := io.Copy(o, f); err != nil {
			log.Printf("Failed to copy [%s]", err.Error())
		}
	}

	return filePaths, nil
}

func isAllBlack(path string) bool {
	result := utils.Exec("magick", "identify", "-format", "%[fx:mean.g]", path)
	if result.ExitCode != 0 {
		log.Printf("Failed to identify image [%s] with code[%d] and error [%s]", path, result.ExitCode, result.StdErr)
		return false
	}

	mean, err := strconv.ParseFloat(result.StdOut, 64)
	if err != nil {
		log.Printf("Failed to parse mean for image [%s] output [%s]", path, result.StdOut)
		return false
	}

	return mean <= MIN_DIFFERENCE //Truncate glare effects and artifacts
}
