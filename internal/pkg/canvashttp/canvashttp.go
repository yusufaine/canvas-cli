package canvashttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"golang.org/x/exp/slices"
)

type CanvasClient struct {
	client      *http.Client
	accessToken string
}

const (
	CanvasHost = "canvas.nus.edu.sg"
	ApiPath    = "/api/v1"
)

func NewClient(token string, opts ...canvasClientOpts) *CanvasClient {
	if token == "" {
		panic("Canvas access token required")
	}

	cc := &CanvasClient{
		client:      http.DefaultClient,
		accessToken: token,
	}

	for _, opt := range opts {
		opt(cc)
	}

	return cc
}

// Returns a map of the canvas course ID to the NUS course code
// Returns an array of course info
func (c CanvasClient) GetCurrentlyEnrolledCourses() []CourseInfo {
	endpoint := url.URL{
		Scheme: "https",
		Host:   CanvasHost,
		Path:   fmt.Sprint(ApiPath + "/users/self/courses"),
		RawQuery: url.Values{
			"enrollment_state": {"active"},
		}.Encode(),
	}

	log.Debug("requesting enrolled courses", "url", endpoint.String())
	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		log.Error("unable to build request for currently enrolled courses", "error", err)
		os.Exit(1)
	}

	body, err := c.doRequest(req)
	if err != nil {
		log.Error("unable to resolve request for currently enrolled courses", "error", err)
		os.Exit(1)
	}

	var courses []CourseInfo
	if err := json.Unmarshal(body, &courses); err != nil {
		log.Error("unable to unmarshal currently enrolled courses response", "error", err)
		os.Exit(1)
	}

	nusCourseCodes := make([]string, 0, len(courses))
	for _, v := range courses {
		nusCourseCodes = append(nusCourseCodes, v.NusCode)
	}

	slices.Sort(nusCourseCodes)
	log.Info("obtained enrolled courses...", "courses", strings.Join(nusCourseCodes, ","))

	return courses
}

// Returns a map of the canvas course ID to the NUS course code
func (c CanvasClient) GetFilesInCourse(canvasId int, nusCode string) []FileInfo {
	endpoint := url.URL{
		Scheme: "https",
		Host:   CanvasHost,
		Path:   fmt.Sprintf("%s/courses/%d/files", ApiPath, canvasId),
	}

	log.Debug("requesting files...", "id", canvasId, "code", nusCode)
	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		log.Error("unable to build request for course files",
			"course_code", nusCode,
			"error", err,
		)
		os.Exit(1)
	}

	body, err := c.doRequest(req)
	if err != nil {
		log.Error("unable to resolve request for course files",
			"course_code", nusCode,
			"error", err,
		)
		os.Exit(1)
	}

	var courseFiles []FileInfo
	if err := json.Unmarshal(body, &courseFiles); err != nil {
		log.Error("unable to unmarshal course files",
			"course_code", nusCode,
			"error", err,
		)
		os.Exit(1)
	}

	log.Info("obtained course file info...", "course_code", nusCode)

	return courseFiles
}

func (c CanvasClient) DownloadFile(nusCode string, info FileInfo, current, max int) bool {
	path := fmt.Sprintf("%s/%s", nusCode, info.DisplayName)
	if _, err := os.Stat(path); err == nil {
		log.Info(fmt.Sprintf("[%d/%d] file exists, skipping...", current, max), "path", path)
		return false
	}
	req, err := http.NewRequest(http.MethodGet, info.DownloadLink, nil)
	if err != nil {
		log.Error("unable to build download request",
			"file_name", info.DisplayName,
			"nus_code", nusCode,
			"error", err,
		)
		return false
	}

	resp, err := c.client.Do(req)
	if err != nil || resp.StatusCode >= 400 {
		log.Error("unable to download",
			"code", resp.StatusCode,
			"file_name", info.DisplayName,
			"nus_code", nusCode,
			"error", err,
		)
		return false
	}

	if err := os.MkdirAll(filepath.Dir(path), 0775); err != nil {
		log.Error("unable to create file path",
			"path", path,
			"error", err,
		)
		return false
	}

	newFile, err := os.Create(path)
	if err != nil {
		log.Error("unable to create file",
			"path", path,
			"error", err,
		)
		return false
	}

	if _, err := io.Copy(newFile, resp.Body); err != nil {
		log.Error("unable to copy downloaded content to new file",
			"path", path,
			"error", err,
		)
		return false
	}

	return true
}

// Helper function add the bearer token to each request and resolve them
func (c *CanvasClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error response")
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}
