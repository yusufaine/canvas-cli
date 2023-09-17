package canvas

import (
	"fmt"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/yusufaine/canvas-cli/internal/pkg/canvashttp"
	"github.com/yusufaine/canvas-cli/internal/pkg/stringset"
)

func Start(config *Config) {
	var wg sync.WaitGroup

	cc := canvashttp.NewClient(config.AccessToken)
	idCodeMap := cc.GetCurrentlyEnrolledCourses()
	nusCodeToFileInfo := getNusCodeFileMap(cc, idCodeMap)

	extBlacklist := stringset.FromElements("mp4")
	filterFiles(nusCodeToFileInfo, extBlacklist)

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

// Filters out files with the extensions in the blacklist.
func filterFiles(nusCodeFileInfo map[string][]canvashttp.FileInfo, extBlacklist *stringset.Stringset) {
	for nusCode, infos := range nusCodeFileInfo {
		var filteredInfo []canvashttp.FileInfo
		for _, info := range infos {
			split := strings.Split(info.DisplayName, ".")
			ext := split[len(split)-1]

			if extBlacklist.Contains(ext) {
				continue
			}

			filteredInfo = append(filteredInfo, info)
		}

		nusCodeFileInfo[nusCode] = filteredInfo
	}
}
