
---

**`docs/tasks.md`** — versão completa reescrita:

```markdown
# Plano de Desenvolvimento - Tasks

> **Convenção de "Done":** uma task só é marcada como concluída quando a funcionalidade está implementada **e** testada manualmente conforme o critério descrito.

---

## Fase 1: Setup e Estrutura Base

- [ ] **Task 1.1:** Instalar Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`) e dependências de compilação no Windows (GCC via TDM-GCC ou MSYS2, Node.js 18+).
  - *Done quando:* `wails doctor` retorna sem erros.

- [ ] **Task 1.2:** Inicializar projeto Wails com template Svelte: `wails init -n drawio-viewer -t svelte`.
  - *Done quando:* `wails dev` sobe a janela desktop com a tela padrão do template.

- [ ] **Task 1.3:** Configurar `wails.json` com nome (`"Drawio Viewer"`), ID (`"com.drawio.viewer"`) e ícone da aplicação.

- [ ] **Task 1.4:** Criar a estrutura de pastas `internal/services/` e `internal/models/` conforme definido no PRD.

- [ ] **Task 1.5:** Instalar e configurar Tailwind CSS no frontend Svelte (via `npm install -D tailwindcss` + `vite.config.js`).
  - *Done quando:* uma classe Tailwind aplicada em `App.svelte` afeta o visual.

- [ ] **Task 1.6:** Baixar `viewer.min.js` do repositório oficial do draw.io e colocar em `frontend/public/viewer.min.js`.
  - Fonte: `https://github.com/jgraph/drawio/blob/dev/src/main/webapp/js/viewer.min.js`

---

## Fase 2: Backend (Go)

- [ ] **Task 2.1:** Criar `internal/models/diagram.go` com a struct `Diagram` (campos: `ID`, `Name`, `XMLPath`, `IsTemp`).

- [ ] **Task 2.2:** Implementar `internal/services/file_service.go` com o método `OpenLocalFile()`:
  - Usa `runtime.OpenFileDialog` do Wails com filtro `*.drawio`.
  - Retorna `models.Diagram` preenchido com o caminho do arquivo selecionado.
  - *Done quando:* ao chamar a função via `wails dev`, o diálogo abre e o caminho retorna corretamente.

- [ ] **Task 2.3:** Implementar `internal/services/download_service.go` com o método `DownloadFromURL(url string)`:
  - Valida se a URL tem esquema `http://` ou `https://`.
  - Faz `http.Get(url)` com timeout de 10 segundos.
  - Salva o conteúdo em `os.TempDir()` com nome `drawio-viewer-<uuid>.drawio`.
  - Retorna `models.Diagram` com `IsTemp: true`.
  - *Done quando:* dado uma URL válida de um `.drawio` público, o arquivo aparece em `%TEMP%`.

- [ ] **Task 2.4:** Criar estrutura thread-safe em `app.go` para gerenciar a lista de arquivos temporários ativos:
  ```go
  type App struct {
      ctx      context.Context
      tempFiles []string
      mu        sync.Mutex
  }
  func (a *App) registerTempFile(path string) { ... }

  Fase 3: Frontend (UI & Render)
 Task 3.1: Incluir viewer.min.js no frontend/index.html via tag <script>.

Done quando: o objeto GraphViewer está disponível no console do wails dev.
 Task 3.2: Criar Svelte store em frontend/src/stores/diagrams.js:
 // store que mantém array de { id, name, xmlPath, isTemp }
export const diagrams = writable([]);
export const activeTabId = writable(null);

Task 3.3: Criar componente TopBar.svelte com:

Botão "Abrir Local" → chama binding FileService.OpenLocalFile().
Botão "Abrir URL" → abre modal/input para colar URL → chama FileService.DownloadFromURL(url).
Ambos adicionam o Diagram retornado ao store diagrams.
 Task 3.4: Criar componente TabBar.svelte:

Renderiza uma aba por item no store diagrams.
Aba ativa é controlada por activeTabId.
Botão × em cada aba remove o diagrama do store.
 Task 3.5: Criar componente DiagramCanvas.svelte:

Recebe xmlPath como prop.
Lê o conteúdo XML via fetch local (ou binding Go auxiliar ReadFile(path)).
Renderiza usando a API do viewer

Done quando: um arquivo .drawio de teste é renderizado corretamente no canvas.
 Task 3.6: Implementar controles de Zoom (botões + e -) e seletor de páginas integrados à API do viewer.min.js.

Fase 4: Integração e Refinamento
 Task 4.1: Implementar binding auxiliar ReadFile(path string) (string, error) em Go para que o frontend possa ler o XML do disco.

 Task 4.2: Adicionar tratamento de erros no frontend:

Exibir toast/notificação para: URL inválida, arquivo não encontrado, XML malformado.
Erros retornados pelo Go via binding devem ser capturados no bloco catch do JS.
 Task 4.3: Estilização final com Tailwind CSS:

Tema escuro (bg-gray-900, text-gray-100).
TopBar fixa no topo, TabBar abaixo, DiagramCanvas ocupa o restante da altura (h-screen).
 Task 4.4: Garantir que ao fechar uma aba com IsTemp: true, o arquivo temporário seja deletado imediatamente (não apenas no shutdown) — adicionar binding DeleteTempFile(path string).

 Fase 5: Build e Entrega
 Task 5.1: Gerar build de produção: wails build -platform windows/amd64.

Done quando: o .exe gerado em build/bin/ é executado sem depender do Wails dev server.
 Task 5.2: Smoke test completo:

 Abrir arquivo local → diagrama renderiza.
 Abrir via URL válida → diagrama renderiza.
 Abrir via URL inválida → mensagem de erro exibida.
 Abrir múltiplas abas → todas funcionam independentemente.
 Fechar app → verificar que %TEMP%\drawio-viewer-* não existe mais.
 Task 5.3: Verificar tamanho do .exe final (meta: abaixo de 50MB conforme RNF04).