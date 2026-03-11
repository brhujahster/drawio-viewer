package services

import (
    "context"
    "path/filepath"

    "drawio-viewer/internal/models"

    "github.com/google/uuid"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type FileService struct {
    ctx context.Context
}

func NewFileService(ctx context.Context) *FileService {
    return &FileService{ctx: ctx}
}

func (f *FileService) OpenLocalFile() (models.Diagram, error) {
    filePath, err := runtime.OpenFileDialog(f.ctx, runtime.OpenDialogOptions{
        Title: "Abrir arquivo Draw.io",
        Filters: []runtime.FileFilter{
            {DisplayName: "Draw.io Files (*.drawio)", Pattern: "*.drawio;*.xml"},
        },
    })
    if err != nil {
        return models.Diagram{}, err
    }
    if filePath == "" {
        return models.Diagram{}, nil
    }

    return models.Diagram{
        ID:      uuid.New().String(),
        Name:    filepath.Base(filePath),
        XMLPath: filePath,
        IsTemp:  false,
    }, nil
}