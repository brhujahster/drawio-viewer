<script>
  import { tick } from 'svelte';
  import { showToast } from '../stores/toast.js';

  export let xmlPath = '';

  let container;
  let viewer = null;
  let pages = [];
  let currentPage = 0;
  let error = null;
  let loadedPath = '';
  let loading = false;

  $: if (container && xmlPath && xmlPath !== loadedPath) {
    loadDiagram(xmlPath);
  }

  async function loadDiagram(path) {
    loadedPath = path;
    error = null;
    viewer = null;
    currentPage = 0;
    pages = [];
    loading = true;

    try {
      const xml = await window['go']['main']['App']['ReadFile'](path);

      if (!xml || xml.trim() === '') {
        throw new Error('O arquivo está vazio.');
      }

      const parseError = getXmlParseError(xml);
      if (parseError) {
        throw new Error('XML malformado: ' + parseError);
      }

      await tick();
      renderDiagram(xml);
    } catch (e) {
      const msg = String(e).replace(/^Error:\s*/, '');
      error = msg;
      showToast(msg, 'error');
    } finally {
      loading = false;
    }
  }

  function getXmlParseError(xml) {
    const parser = new DOMParser();
    const doc = parser.parseFromString(xml, 'text/xml');
    const errNode = doc.querySelector('parsererror');
    if (errNode) {
      const text = errNode.textContent?.split('\n')[0] ?? 'erro desconhecido';
      return text;
    }
    return null;
  }

  function renderDiagram(xml) {
    if (!container) return;
    container.innerHTML = '';

    const parser = new DOMParser();
    const doc = parser.parseFromString(xml, 'text/xml');
    const diagramEls = doc.querySelectorAll('diagram');

    pages =
      diagramEls.length > 0
        ? Array.from(diagramEls).map((d, i) => d.getAttribute('name') || `Página ${i + 1}`)
        : ['Página 1'];

    const config = JSON.stringify({
      highlight: '#0000ff',
      nav: true,
      resize: true,
      page: currentPage,
      xml,
    });

    const div = document.createElement('div');
    div.className = 'mxgraph';
    div.style.cssText = 'width:100%;height:100%;';
    div.setAttribute('data-mxgraph', config);
    container.appendChild(div);

    if (window.GraphViewer) {
      viewer = window.GraphViewer.processElement(div);
    } else {
      showToast('GraphViewer não disponível. Verifique se viewer.min.js foi carregado.', 'error');
    }
  }

  function zoomIn() {
    viewer?.graph?.zoomIn();
  }

  function zoomOut() {
    viewer?.graph?.zoomOut();
  }

  function zoomReset() {
    viewer?.graph?.fit();
  }

  function changePage(idx) {
    currentPage = idx;
    loadedPath = '';
    loadDiagram(xmlPath);
  }
</script>

<div class="flex flex-col h-full">
  <div class="flex items-center gap-1 bg-gray-800 px-3 py-1 border-b border-gray-700 flex-shrink-0 h-9">
    <button
      on:click={zoomOut}
      class="bg-gray-700 hover:bg-gray-600 text-white w-7 h-6 rounded font-bold text-base
             flex items-center justify-center transition-colors"
      title="Diminuir zoom"
    >−</button>

    <button
      on:click={zoomReset}
      class="bg-gray-700 hover:bg-gray-600 text-white px-2 h-6 rounded text-xs transition-colors"
      title="Ajustar à tela"
    >Fit</button>

    <button
      on:click={zoomIn}
      class="bg-gray-700 hover:bg-gray-600 text-white w-7 h-6 rounded font-bold text-base
             flex items-center justify-center transition-colors"
      title="Aumentar zoom"
    >+</button>

    {#if pages.length > 1}
      <div class="w-px h-5 bg-gray-600 mx-2"></div>
      <span class="text-gray-400 text-xs">Página:</span>
      <select
        value={currentPage}
        on:change={(e) => changePage(+e.target.value)}
        class="bg-gray-700 text-white text-xs rounded px-2 h-6 border border-gray-600 outline-none cursor-pointer"
      >
        {#each pages as page, i}
          <option value={i}>{page}</option>
        {/each}
      </select>
    {/if}

    {#if loading}
      <span class="text-gray-400 text-xs ml-auto animate-pulse">Carregando...</span>
    {/if}
  </div>

  {#if error}
    <div class="flex-1 flex flex-col items-center justify-center gap-3 text-red-400 text-sm px-8 text-center">
      <span class="text-3xl">⚠</span>
      <p>{error}</p>
    </div>
  {:else}
    <div bind:this={container} class="flex-1 bg-white overflow-hidden"></div>
  {/if}
</div>