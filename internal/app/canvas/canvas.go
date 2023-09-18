package canvas

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/canvas-cli/internal/pkg/canvashttp"
)

func Start(config *Config) {
	var wg sync.WaitGroup

	cc := canvashttp.NewClient(config.AccessToken,
		canvashttp.WithHost(config.Host),
		canvashttp.WithPathApiPrefix(config.ApiPath),
	)
	idCodeMap := cc.GetCurrentlyEnrolledCourses()
	nusCodeToFileInfo := getNusCodeFileMap(cc, idCodeMap)

	filterFiles(config, nusCodeToFileInfo)

	for code, fileInfos := range nusCodeToFileInfo {
		wg.Add(1)
		go func(code string, fileInfos []canvashttp.FileInfo) {
			defer wg.Done()
			var (
				pl      = log.With("nus_code", code)
				i   int = 0
				max int = len(fileInfos)
			)
			for _, fileInfo := range fileInfos {
				i++
				if !cc.DownloadFile(code, fileInfo, i, max) {
					continue
				}
				pl.Info(fmt.Sprintf("[%d/%d] download completed", i, max), "file", fileInfo.DisplayName)
			}
		}(code, fileInfos)
	}
	wg.Wait()
}

func getNusCodeFileMap(cc *canvashttp.CanvasClient, enrolledCourses []canvashttp.CourseInfo) map[string][]canvashttp.FileInfo {
	nusCodeFileInfo := make(map[string][]canvashttp.FileInfo)
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
func filterFiles(config *Config, nusCodeFileInfo map[string][]canvashttp.FileInfo) {
	maxSizeBytes := config.MaxSizeMb * 1_000_000
	for nusCode, infos := range nusCodeFileInfo {
		var filteredInfo []canvashttp.FileInfo
		// TODO: add more filter here if needed
		for _, info := range infos {
			if info.Size > maxSizeBytes {
				log.Debug("filtered file out...", "filename", info.DisplayName, "size", info.Size, "max_size", maxSizeBytes)
				continue
			}
			filteredInfo = append(filteredInfo, info)
		}

		nusCodeFileInfo[nusCode] = filteredInfo
	}
}
