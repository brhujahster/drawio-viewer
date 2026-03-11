<script>
  import { tick } from 'svelte';

  export let xmlPath = '';

  let container;
  let viewer = null;
  let pages = [];
  let currentPage = 0;
  let error = null;
  let loadedPath = '';

  $: if (container && xmlPath && xmlPath !== loadedPath) {
    loadDiagram(xmlPath);
  }

  async function loadDiagram(path) {
    loadedPath = path;
    error = null;
    viewer = null;
    currentPage = 0;
    pages = [];

    try {
      const xml = await window['go']['main']['App']['ReadFile'](path);
      await tick();
      renderDiagram(xml);
    } catch (e) {
      error = 'Erro ao ler arquivo: ' + e;
    }
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
  <div class="flex items-center gap-1 bg-gray-800 px-3 py-1 border-b border-gray-700 flex-shrink-0">
    <button
      on:click={zoomOut}
      class="bg-gray-700 hover:bg-gray-600 text-white w-7 h-7 rounded text-lg font-bold flex items-center justify-center transition-colors"
      title="Zoom -"
    >−</button>

    <button
      on:click={zoomReset}
      class="bg-gray-700 hover:bg-gray-600 text-white px-2 h-7 rounded text-xs transition-colors"
      title="Ajustar à tela"
    >Fit</button>

    <button
      on:click={zoomIn}
      class="bg-gray-700 hover:bg-gray-600 text-white w-7 h-7 rounded text-lg font-bold flex items-center justify-center transition-colors"
      title="Zoom +"
    >+</button>

    {#if pages.length > 1}
      <span class="text-gray-400 text-sm ml-4">Página:</span>
      <select
        value={currentPage}
        on:change={(e) => changePage(+e.target.value)}
        class="bg-gray-700 text-white text-sm rounded px-2 h-7 border border-gray-600 outline-none cursor-pointer"
      >
        {#each pages as page, i}
          <option value={i}>{page}</option>
        {/each}
      </select>
    {/if}
  </div>

  {#if error}
    <div class="flex-1 flex items-center justify-center text-red-400 text-sm">{error}</div>
  {:else}
    <div bind:this={container} class="flex-1 bg-white overflow-hidden"></div>
  {/if}
</div>