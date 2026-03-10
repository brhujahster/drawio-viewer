# Product Requirement Document (PRD) - Drawio Go Viewer (POC)

## 1. Visão Geral
O **Drawio Go Viewer** é uma aplicação desktop leve para Windows que permite a visualização rápida de diagramas do draw.io (`.drawio`), suportando arquivos locais e remotos via URL. É construída com **Wails v2** (Go no backend, Svelte no frontend) e empacotada como um único executável `.exe`.

## 2. Objetivos
- Prover uma ferramenta rápida para visualização de arquiteturas e fluxogramas sem necessidade de abrir o navegador ou o editor completo.
- Validar a integração do Go (Wails v2) com a biblioteca de renderização oficial do draw.io.

## 3. Stack Tecnológica

| Camada | Tecnologia | Versão |
|---|---|---|
| Runtime desktop | Wails | v2.x |
| Backend | Go | 1.21+ |
| Frontend framework | Svelte | 4.x |
| Estilização | Tailwind CSS | 3.x |
| Renderização de diagramas | draw.io `viewer.min.js` | latest (CDN ou local) |
| Build tool frontend | Vite | 4.x (incluído no template Wails+Svelte) |

## 4. Arquitetura em Camadas
/ ├── main.go # Ponto de entrada, inicializa Wails app ├── app.go # Struct App: registra hooks OnStartup/OnShutdown ├── wails.json # Config da app (nome, ícone, ID) ├── internal/ │ ├── services/ │ │ ├── file_service.go # Abertura de arquivos locais via diálogo nativo │ │ └── download_service.go # Download de XML via URL para pasta temp │ └── models/ │ └── diagram.go # Struct Diagram compartilhada entre camadas └── frontend/ ├── index.html # HTML base, inclui viewer.min.js ├── src/ │ ├── App.svelte # Componente raiz, layout principal │ ├── components/ │ │ ├── TopBar.svelte # Botões "Abrir Local" e "Abrir URL" │ │ ├── TabBar.svelte # Lista de abas com diagramas abertos │ │ └── DiagramCanvas.svelte # Renderiza XML via API do viewer │ └── stores/ │ └── diagrams.js # Svelte store: lista de diagramas abertos └── public/ └── viewer.min.js # Script oficial draw.io (copiado localmente)


## 5. Contrato da API Go → Frontend (Wails Bindings)

Funções expostas ao JavaScript via Wails:

```go
// Abre diálogo nativo de seleção de arquivo e retorna o Diagram
func (f *FileService) OpenLocalFile() (models.Diagram, error)

// Faz download do XML da URL, salva em temp e retorna o Diagram
func (f *FileService) DownloadFromURL(url string) (models.Diagram, error)

type Diagram struct {
    ID      string `json:"id"`       // UUID gerado no momento da abertura
    Name    string `json:"name"`     // Nome do arquivo (ex: "arquitetura.drawio")
    XMLPath string `json:"xmlPath"`  // Caminho absoluto do arquivo no disco
    IsTemp  bool   `json:"isTemp"`   // true = arquivo temporário (download via URL)
}