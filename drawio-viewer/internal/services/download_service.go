package services

import (
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"

    "drawio-viewer/internal/models"

    "github.com/google/uuid"
)

type DownloadService struct{}

func NewDownloadService() *DownloadService {
    return &DownloadService{}
}

func convertToRawURL(rawURL string) string {
    reGitHubBlob := regexp.MustCompile(`^https://github\.com/([^/]+/[^/]+)/blob/(.+)$`)
    if m := reGitHubBlob.FindStringSubmatch(rawURL); m != nil {
        return "https://raw.githubusercontent.com/" + m[1] + "/" + m[2]
    }

    if strings.Contains(rawURL, "/-/blob/") {
        return strings.Replace(rawURL, "/-/blob/", "/-/raw/", 1)
    }

    return rawURL
}

// extractToken extrai e remove o private_token da URL, retornando (urlLimpa, token).
// Suporta: ?private_token=TOKEN ou &private_token=TOKEN
func extractToken(rawURL string) (string, string) {
    parsed, err := url.Parse(rawURL)
    if err != nil {
        return rawURL, ""
    }
    q := parsed.Query()
    token := q.Get("private_token")
    if token == "" {
        return rawURL, ""
    }
    q.Del("private_token")
    parsed.RawQuery = q.Encode()
    return parsed.String(), token
}

func isHTMLContent(body []byte) bool {
    preview := strings.ToLower(strings.TrimSpace(string(body[:min(512, len(body))])))
    return strings.HasPrefix(preview, "<!doctype") ||
        strings.HasPrefix(preview, "<html") ||
        strings.Contains(preview, "<title>sign in</title>") ||
        strings.Contains(preview, "<title>login</title>")
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func (d *DownloadService) DownloadFromURL(rawURL string) (models.Diagram, error) {
    parsed, err := url.Parse(rawURL)
    if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
        return models.Diagram{}, fmt.Errorf("URL inválida: deve usar http:// ou https://")
    }

    effectiveURL := convertToRawURL(rawURL)
    effectiveURL, token := extractToken(effectiveURL)

    req, err := http.NewRequest("GET", effectiveURL, nil)
    if err != nil {
        return models.Diagram{}, fmt.Errorf("erro ao criar requisição: %w", err)
    }

    if token != "" {
        req.Header.Set("PRIVATE-TOKEN", token)
    }

    client := &http.Client{Timeout: 15 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return models.Diagram{}, fmt.Errorf("erro ao baixar arquivo: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
        return models.Diagram{}, fmt.Errorf(
            "acesso negado (%s). Para repositórios privados, adicione ?private_token=SEU_TOKEN à URL",
            resp.Status,
        )
    }

    if resp.StatusCode != http.StatusOK {
        return models.Diagram{}, fmt.Errorf("resposta inesperada do servidor: %s", resp.Status)
    }

    contentType := resp.Header.Get("Content-Type")
    if strings.Contains(contentType, "text/html") {
        return models.Diagram{}, fmt.Errorf(
            "servidor retornou HTML (possível redirect de login). Para repositórios privados, adicione ?private_token=SEU_TOKEN à URL",
        )
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return models.Diagram{}, fmt.Errorf("erro ao ler resposta: %w", err)
    }

    if isHTMLContent(body) {
        return models.Diagram{}, fmt.Errorf(
            "conteúdo recebido é HTML (repositório privado?). Adicione ?private_token=SEU_TOKEN à URL",
        )
    }

    id := uuid.New().String()
    tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("drawio-viewer-%s.drawio", id))

    if err = os.WriteFile(tempPath, body, 0644); err != nil {
        return models.Diagram{}, fmt.Errorf("erro ao salvar arquivo temporário: %w", err)
    }

    name := filepath.Base(parsed.Path)
    if name == "" || name == "." {
        name = "diagrama.drawio"
    }

    return models.Diagram{
        ID:      id,
        Name:    name,
        XMLPath: tempPath,
        IsTemp:  true,
    }, nil
}