<script>
  import { diagrams, activeTabId } from '../stores/diagrams.js';
  import { showToast } from '../stores/toast.js';

  let showUrlInput = false;
  let urlValue = '';
  let loadingLocal = false;
  let loadingUrl = false;

  async function openLocal() {
    loadingLocal = true;
    try {
      const diagram = await window['go']['main']['App']['OpenLocalFile']();
      if (!diagram || !diagram.xmlPath) return;
      diagrams.update(d => [...d, diagram]);
      activeTabId.set(diagram.id);
    } catch (e) {
      showToast('Erro ao abrir arquivo: ' + e, 'error');
    } finally {
      loadingLocal = false;
    }
  }

  async function openUrl() {
    const url = urlValue.trim();
    if (!url) return;

    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      showToast('URL inválida: deve começar com http:// ou https://', 'error');
      return;
    }

    loadingUrl = true;
    try {
      const diagram = await window['go']['main']['App']['DownloadFromURL'](url);
      if (!diagram || !diagram.xmlPath) return;
      diagrams.update(d => [...d, diagram]);
      activeTabId.set(diagram.id);
      showToast('Diagrama carregado com sucesso!', 'success');
      showUrlInput = false;
      urlValue = '';
    } catch (e) {
      showToast('Erro ao baixar URL: ' + e, 'error');
    } finally {
      loadingUrl = false;
    }
  }

  function cancelUrl() {
    showUrlInput = false;
    urlValue = '';
  }
</script>

<div class="flex items-center gap-2 bg-gray-800 px-4 py-2 border-b border-gray-700 flex-shrink-0 h-12">
  <span class="text-white font-semibold mr-2 text-sm tracking-wide">Draw.io Viewer</span>

  <div class="w-px h-5 bg-gray-600 mx-1"></div>

  <button
    on:click={openLocal}
    disabled={loadingLocal}
    class="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed
           text-white px-3 py-1 rounded text-sm transition-colors flex items-center gap-1"
  >
    {#if loadingLocal}
      <span class="animate-spin text-xs">⟳</span>
    {/if}
    Abrir Local
  </button>

  <button
    on:click={() => (showUrlInput = !showUrlInput)}
    class="bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded text-sm transition-colors"
  >
    Abrir URL
  </button>

  {#if showUrlInput}
    <div class="flex items-center gap-2">
      <input
        bind:value={urlValue}
        on:keydown={(e) => e.key === 'Enter' && openUrl()}
        on:keydown={(e) => e.key === 'Escape' && cancelUrl()}
        placeholder="https://exemplo.com/diagrama.drawio"
        class="bg-gray-700 text-white px-3 py-1 rounded text-sm border border-gray-600 w-96
               outline-none focus:border-blue-500 transition-colors"
        autofocus
      />
      <button
        on:click={openUrl}
        disabled={loadingUrl}
        class="bg-green-700 hover:bg-green-800 disabled:opacity-50 disabled:cursor-not-allowed
               text-white px-3 py-1 rounded text-sm transition-colors flex items-center gap-1"
      >
        {#if loadingUrl}
          <span class="animate-spin text-xs">⟳</span>
        {/if}
        OK
      </button>
      <button
        on:click={cancelUrl}
        class="text-gray-400 hover:text-white px-2 py-1 text-sm transition-colors"
      >
        ✕
      </button>
    </div>
  {/if}
</div>