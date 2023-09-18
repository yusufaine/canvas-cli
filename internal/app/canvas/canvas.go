package canvas

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/canvas-cli/internal/pkg/canvasclient"
)

func Start(config *Config) {
	var wg sync.WaitGroup

	cc := canvasclient.NewClient(config.AccessToken,
		canvasclient.WithHost(config.Host),
		canvasclient.WithPathApiPrefix(config.ApiPath),
	)
	idCodeMap := cc.GetCurrentlyEnrolledCourses()
	nusCodeToFileInfo := getNusCodeFileMap(cc, idCodeMap)

	filterFiles(config, nusCodeToFileInfo)

	for code, fileInfos := range nusCodeToFileInfo {
		wg.Add(1)
		go func(code string, fileInfos []canvasclient.FileInfo) {
			defer wg.Done()
			var (
				pl      = log.With("nus_code", code)
				i   int = 0
				max int = len(fileInfos)
			)
			for _, fileInfo := range fileInfos {
				i++
				pl.Info("attempting to download...", "file", fileInfo.DisplayName)
				if !cc.DownloadFile(code, fileInfo, i, max) {
					continue
				}
				pl.Info(fmt.Sprintf("[%d/%d] download completed", i, max), "file", fileInfo.DisplayName)
			}
		}(code, fileInfos)
	}
	wg.Wait()
}

func getNusCodeFileMap(cc *canvasclient.CanvasClient, enrolledCourses []canvasclient.CourseInfo) map[string][]canvasclient.FileInfo {
	nusCodeFileInfo := make(map[string][]canvasclient.FileInfo)
	for _, course := range enrolledCourses {
		id := course.CanvasId
		code := course.NusCode
		log.Debug("getting file info for", "id", id, "code", code)
		files := cc.GetFilesInCourse(id, code)
		nusCodeFileInfo[code] = files
	}

	return nusCodeFileInfo
}

// Filters out files with a size larger than config.MaxSize
func filterFiles(config *Config, nusCodeFileInfo map[string][]canvasclient.FileInfo) {
	maxSizeBytes := config.MaxSizeMb * 1_000_000

	for nusCode, infos := range nusCodeFileInfo {
		var filteredInfo []canvasclient.FileInfo
		// TODO: add more filter here if needed
		for _, info := range infos {
			if !isWhitelistedExtension(config.ExtWhitelist, info.DisplayName) {
				continue
			}

			if isLargerThanMaxBytes(info.Size, maxSizeBytes, info.DisplayName) {
				continue
			}
			filteredInfo = append(filteredInfo, info)
		}

		nusCodeFileInfo[nusCode] = filteredInfo
	}
}

// Returns true if:
//  1. File is whitelisted, or
//  2. If whitelist map is empty, all extensions are allowed
func isWhitelistedExtension(whitelist map[string]struct{}, displayName string) bool {
	if len(whitelist) == 0 {
		return true
	}

	nameSplit := strings.Split(displayName, ".")
	ext := nameSplit[len(nameSplit)-1]
	_, ok := whitelist[ext]
	if !ok {
		log.Warn("extension not in whitelist, skipping", "filename", displayName)
	}
	return ok
}

// Returns true if the current file size is less than the maximum file size
func isLargerThanMaxBytes(fileSize, maxSizeBytes int, displayName string) bool {
	if fileSize > maxSizeBytes {
		log.Warn(fmt.Sprintf("file larger than %dMB", maxSizeBytes/1_000_000), "filename", displayName, "size", fileSize)
		return true
	}
	return false
}
