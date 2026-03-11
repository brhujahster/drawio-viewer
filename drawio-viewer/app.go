package main

import (
    "context"
    "fmt"
    "sync"

    "drawio-viewer/internal/models"
    "drawio-viewer/internal/services"
)

type App struct {
    ctx             context.Context
    tempFiles       []string
    mu              sync.Mutex
    fileService     *services.FileService
    downloadService *services.DownloadService
}

func NewApp() *App {
    return &App{
        downloadService: services.NewDownloadService(),
    }
}

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    a.fileService = services.NewFileService(ctx)
}

func (a *App) registerTempFile(path string) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.tempFiles = append(a.tempFiles, path)
}

func (a *App) OpenLocalFile() (models.Diagram, error) {
    return a.fileService.OpenLocalFile()
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

func (a *App) Greet(name string) string {
    return fmt.Sprintf("Hello %s, It's show time!", name)
}