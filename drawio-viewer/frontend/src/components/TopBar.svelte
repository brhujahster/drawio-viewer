<script>
  import { diagrams, activeTabId } from '../stores/diagrams.js';

  let showUrlInput = false;
  let urlValue = '';

  async function openLocal() {
    try {
      const diagram = await window['go']['main']['App']['OpenLocalFile']();
      if (!diagram || !diagram.xmlPath) return;
      diagrams.update(d => [...d, diagram]);
      activeTabId.set(diagram.id);
    } catch (e) {
      alert('Erro ao abrir arquivo: ' + e);
    }
  }

  async function openUrl() {
    const url = urlValue.trim();
    if (!url) return;
    try {
      const diagram = await window['go']['main']['App']['DownloadFromURL'](url);
      if (!diagram || !diagram.xmlPath) return;
      diagrams.update(d => [...d, diagram]);
      activeTabId.set(diagram.id);
      showUrlInput = false;
      urlValue = '';
    } catch (e) {
      alert('Erro ao baixar URL: ' + e);
    }
  }

  function cancelUrl() {
    showUrlInput = false;
    urlValue = '';
  }
</script>

<div class="flex items-center gap-2 bg-gray-800 px-4 py-2 border-b border-gray-700 flex-shrink-0">
  <span class="text-white font-semibold mr-2 text-sm">Draw.io Viewer</span>

  <button
    on:click={openLocal}
    class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-sm transition-colors"
  >
    Abrir Local
  </button>

  <button
    on:click={() => (showUrlInput = !showUrlInput)}
    class="bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded text-sm transition-colors"
  >
    Abrir URL
  </button>

  {#if showUrlInput}
    <input
      bind:value={urlValue}
      on:keydown={(e) => e.key === 'Enter' && openUrl()}
      placeholder="https://exemplo.com/diagrama.drawio"
      class="bg-gray-700 text-white px-3 py-1 rounded text-sm border border-gray-600 w-96 outline-none focus:border-blue-500"
      autofocus
    />
    <button
      on:click={openUrl}
      class="bg-green-700 hover:bg-green-800 text-white px-3 py-1 rounded text-sm transition-colors"
    >
      OK
    </button>
    <button
      on:click={cancelUrl}
      class="text-gray-400 hover:text-white px-2 py-1 text-sm transition-colors"
    >
      ✕
    </button>
  {/if}
</div>