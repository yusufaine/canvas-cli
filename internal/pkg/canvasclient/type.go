package canvasclient

type CourseInfo struct {
	CanvasId int    `json:"id"`
	NusCode  string `json:"course_code"`
}

type FileInfo struct {
	CanvasFileId    int    `json:"id"`
	CanvasFolderId  int    `json:"folder_id"`
	DisplayName     string `json:"display_name"`
	EscapedFileName string `json:"filename"`
	ContentType     string `json:"content-type"`
	DownloadLink    string `json:"url"`
	Size            int    `json:"size"`
}
