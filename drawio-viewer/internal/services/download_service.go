package services

import (
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "time"

    "drawio-viewer/internal/models"

    "github.com/google/uuid"
)

type DownloadService struct{}

func NewDownloadService() *DownloadService {
    return &DownloadService{}
}

func (d *DownloadService) DownloadFromURL(rawURL string) (models.Diagram, error) {
    parsed, err := url.Parse(rawURL)
    if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
        return models.Diagram{}, fmt.Errorf("URL inválida: deve usar http:// ou https://")
    }

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(rawURL)
    if err != nil {
        return models.Diagram{}, fmt.Errorf("erro ao baixar arquivo: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return models.Diagram{}, fmt.Errorf("resposta inesperada: %s", resp.Status)
    }

    id := uuid.New().String()
    fileName := fmt.Sprintf("drawio-viewer-%s.drawio", id)
    tempPath := filepath.Join(os.TempDir(), fileName)

    out, err := os.Create(tempPath)
    if err != nil {
        return models.Diagram{}, fmt.Errorf("erro ao criar arquivo temporário: %w", err)
    }
    defer out.Close()

    if _, err = io.Copy(out, resp.Body); err != nil {
        return models.Diagram{}, fmt.Errorf("erro ao salvar arquivo: %w", err)
    }

    return models.Diagram{
        ID:      id,
        Name:    filepath.Base(parsed.Path),
        XMLPath: tempPath,
        IsTemp:  true,
    }, nil
}