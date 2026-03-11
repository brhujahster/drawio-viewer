package services

import (
    "path/filepath"

    "drawio-viewer/internal/models"

    "github.com/google/uuid"
    "github.com/wailsapp/wails/v2/pkg/runtime"
    goruntime "runtime"
    gocontext "context"
)

type FileService struct {
    ctx gocontext.Context
}

func NewFileService(ctx gocontext.Context) *FileService {
    return &FileService{ctx: ctx}
}

func (f *FileService) OpenLocalFile() (models.Diagram, error) {
    filePath, err := runtime.OpenFileDialog(f.ctx, runtime.OpenDialogOptions{
        Title: "Abrir arquivo Draw.io",
        Filters: []runtime.FileFilter{
            {DisplayName: "Draw.io Files (*.drawio)", Pattern: "*.drawio"},
        },
    })
    if err != nil {
        return models.Diagram{}, err
    }
    if filePath == "" {
        return models.Diagram{}, nil
    }

    _ = goruntime.GOOS

    return models.Diagram{
        ID:      uuid.New().String(),
        Name:    filepath.Base(filePath),
        XMLPath: filePath,
        IsTemp:  false,
    }, nil
}