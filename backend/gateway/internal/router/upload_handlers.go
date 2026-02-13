package router

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// setupUploadRoutes registers file upload routes.
func (r *Router) setupUploadRoutes() {
	r.mux.HandleFunc("POST /api/v1/files/upload", r.handleFileUpload)
	r.mux.HandleFunc("GET /api/v1/files/{path...}", r.handleFileDownload)
	r.mux.HandleFunc("DELETE /api/v1/files/{path...}", r.handleFileDelete)
}

func (r *Router) handleFileUpload(w http.ResponseWriter, req *http.Request) {
	// Max 20MB
	if err := req.ParseMultipartForm(20 << 20); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "파일 크기 초과 (최대 20MB)"})
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "파일을 찾을 수 없습니다"})
		return
	}
	defer file.Close()

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	category := req.FormValue("category") // profiles, reports, medical, community
	if category == "" {
		category = "general"
	}

	// Generate unique path
	ext := filepath.Ext(header.Filename)
	filePath := fmt.Sprintf("%s/%s/%s%s", category, time.Now().Format("2006/01/02"), uuid.New().String(), ext)

	// S3 업로드
	if r.s3Client != nil {
		if err := r.s3Client.Upload(req.Context(), filePath, file, header.Size, contentType); err != nil {
			log.Printf("[gateway] S3 업로드 실패: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "파일 업로드에 실패했습니다"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"path":         filePath,
			"filename":     header.Filename,
			"size":         header.Size,
			"content_type": contentType,
			"storage":      "s3",
		})
	} else {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"path":         filePath,
			"filename":     header.Filename,
			"size":         header.Size,
			"content_type": contentType,
			"message":      "S3 미설정 — 파일 저장 비활성화",
		})
	}
}

func (r *Router) handleFileDownload(w http.ResponseWriter, req *http.Request) {
	path := req.PathValue("path")
	if path == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "파일 경로 필요"})
		return
	}

	if r.s3Client != nil {
		// Presigned URL 방식
		presignedURL, err := r.s3Client.GetPresignedURL(req.Context(), path, 15*time.Minute)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "파일을 찾을 수 없습니다"})
			return
		}
		http.Redirect(w, req, presignedURL, http.StatusTemporaryRedirect)
	} else {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "S3 미설정 — 파일 다운로드 비활성화",
		})
	}
}

func (r *Router) handleFileDelete(w http.ResponseWriter, req *http.Request) {
	path := req.PathValue("path")
	if path == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "파일 경로 필요"})
		return
	}

	if r.s3Client != nil {
		if err := r.s3Client.Delete(req.Context(), path); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "파일 삭제에 실패했습니다"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"deleted": path,
		})
	} else {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "S3 미설정 — 파일 삭제 비활성화",
		})
	}
}
