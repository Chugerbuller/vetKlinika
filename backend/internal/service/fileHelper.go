package service

import (



	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

)

// Data представляет файл с его метаданными
type Data struct {
	Content    []byte
	Name       string
	Size       int64
	MimeType   string
}

// Read загружает файл с диска
func Read(path string) (*Data, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &Data{
		Content:  content,
		Name:     filepath.Base(path),
		Size:     int64(len(content)),
		MimeType: getMimeType(path),
	}, nil
}

// Save сохраняет файл на диск
func Save(path string, data []byte) error {
	dir := filepath.Dir(path)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0644)
}

// ToReader конвертирует данные в io.Reader
func (d *Data) ToReader() io.Reader {
	return bytes.NewReader(d.Content)
}

// ToMultipart создает multipart часть для HTTP запроса
func (d *Data) ToMultipart(w *multipart.Writer, fieldName string) error {
	part, err := w.CreateFormFile(fieldName, d.Name)
	if err != nil {
		return err
	}
	_, err = part.Write(d.Content)
	return err
}

// NewFromRequest извлекает файл из HTTP запроса
func NewFromRequest(r *http.Request, fieldName string) (*Data, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &Data{
		Content:  content,
		Name:     header.Filename,
		Size:     header.Size,
		MimeType: header.Header.Get("Content-Type"),
	}, nil
}

// ToHTTPResponse отправляет файл как HTTP ответ
func (d *Data) ToHTTPResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", d.MimeType)
	w.Header().Set("Content-Length", string(rune(d.Size)))
	w.Header().Set("Content-Disposition", "attachment; filename="+d.Name)
	w.WriteHeader(http.StatusOK)
	w.Write(d.Content)
}

// getMimeType определяет MIME тип по расширению
func getMimeType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".txt":
		return "text/plain"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}