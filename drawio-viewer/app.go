package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"drawio-viewer/internal/models"
	"drawio-viewer/internal/services"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx             context.Context
	tempFiles       []string
	mu              sync.Mutex
	downloadService *services.DownloadService
}

func NewApp() *App {
	return &App{
		downloadService: services.NewDownloadService(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, f := range a.tempFiles {
		os.Remove(f)
	}
	a.tempFiles = nil
}

func (a *App) registerTempFile(path string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.tempFiles = append(a.tempFiles, path)
}

func (a *App) OpenLocalFile() {
	go func() {
		ps := `Add-Type -AssemblyName System.Windows.Forms;` +
			`$d = New-Object System.Windows.Forms.OpenFileDialog;` +
			`$d.Title = 'Abrir arquivo Draw.io';` +
			`$d.Filter = 'Draw.io Files (*.drawio;*.xml)|*.drawio;*.xml|All Files (*.*)|*.*';` +
			`$d.Multiselect = $false;` +
			`if ($d.ShowDialog() -eq 'OK') { Write-Output $d.FileName }`

		cmd := exec.Command("powershell", "-NonInteractive", "-WindowStyle", "Hidden", "-Command", ps)
		out, err := cmd.Output()
		if err != nil {
			runtime.EventsEmit(a.ctx, "localfile:error", err.Error())
			return
		}

		filePath := strings.TrimSpace(string(out))
		if filePath == "" {
			runtime.EventsEmit(a.ctx, "localfile:cancelled")
			return
		}

		diagram := models.Diagram{
			ID:      uuid.New().String(),
			Name:    filepath.Base(filePath),
			XMLPath: filePath,
			IsTemp:  false,
		}
		runtime.EventsEmit(a.ctx, "localfile:opened", diagram)
	}()
}

func (a *App) DownloadFromURL(url string) (models.Diagram, error) {
	diagram, err := a.downloadService.DownloadFromURL(url)
	if err != nil {
		return models.Diagram{}, err
	}
	if diagram.IsTemp {
		a.registerTempFile(diagram.XMLPath)
	}
	return diagram, nil
}

func (a *App) ReadFile(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("caminho do arquivo não informado")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("arquivo não encontrado: %s", path)
		}
		return "", fmt.Errorf("erro ao ler arquivo: %w", err)
	}
	return string(data), nil
}

func (a *App) DeleteTempFile(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("erro ao deletar arquivo temporário: %w", err)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	filtered := make([]string, 0, len(a.tempFiles))
	for _, f := range a.tempFiles {
		if f != path {
			filtered = append(filtered, f)
		}
	}
	a.tempFiles = filtered
	return nil
}
