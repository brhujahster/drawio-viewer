<script>
  import { tick } from 'svelte';
  import { showToast } from '../stores/toast.js';

  export let xmlPath = '';
  export let xmlContent = '';

  let container;
  let pages = [];
  let currentPage = 0;
  let error = null;
  let loadedKey = '';
  let loading = false;
  let zoomLevel = 1;
  let panX = 0;
  let panY = 0;
  let isPanning = false;
  let startX = 0;
  let startY = 0;

  $: diagramKey = xmlContent ? `content:${xmlContent.length}:${xmlContent.slice(0, 40)}` : `path:${xmlPath}`;

  $: if (container && diagramKey && diagramKey !== loadedKey) {
    loadDiagram();
  }

  async function loadDiagram() {
    loadedKey = diagramKey;
    error = null;
    currentPage = 0;
    pages = [];
    loading = true;
    zoomLevel = 1;
    panX = 0;
    panY = 0;

    try {
      let xml;

      if (xmlContent) {
        xml = xmlContent;
      } else if (xmlPath) {
        xml = await window['go']['main']['App']['ReadFile'](xmlPath);
      } else {
        return;
      }

      if (!xml || xml.trim() === '') throw new Error('Arquivo vazio.');

      const parseError = getXmlParseError(xml);
      if (parseError) throw new Error('XML malformado: ' + parseError);

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
    return errNode ? (errNode.textContent?.split('\n')[0] ?? 'erro desconhecido') : null;
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

    if (!window.GraphViewer) {
      showToast('GraphViewer não disponível. Verifique se viewer.min.js foi carregado.', 'error');
      return;
    }

    const w = container.clientWidth;
    const h = container.clientHeight;

    const config = JSON.stringify({
      highlight: '#0000ff',
      nav: false,
      resize: true,
      page: currentPage,
      xml,
    });

    const wrapper = document.createElement('div');
    wrapper.setAttribute('draggable', 'false');

    wrapper.style.cssText = `position:absolute;top:0;left:0;width:${w}px;height:${h}px;transform-origin:top left;`;

    const div = document.createElement('div');
    div.className = 'mxgraph';
    div.style.cssText = `width:${w}px;height:${h}px;pointer-events:none;`;
    div.setAttribute('draggable', 'false');
    div.setAttribute('data-mxgraph', config);

    wrapper.appendChild(div);
    container.appendChild(wrapper);

    if (typeof GraphViewer.createViewerForElement === 'function') {
      GraphViewer.createViewerForElement(div);
    } else {
      GraphViewer.processElements();
    }

    applyTransform();
  }

  function getWrapper() {
    return container?.firstElementChild ?? null;
  }

  function applyTransform() {
    const el = getWrapper();
    if (!el) return;
    el.style.transform = `translate(${panX}px, ${panY}px) scale(${zoomLevel})`;
  }

  function zoomIn() {
    zoomLevel = Math.min(+(zoomLevel + 0.2).toFixed(2), 5);
    applyTransform();
  }

  function zoomOut() {
    zoomLevel = Math.max(+(zoomLevel - 0.2).toFixed(2), 0.2);
    applyTransform();
  }

  function zoomReset() {
    zoomLevel = 1;
    panX = 0;
    panY = 0;
    applyTransform();
  }

  function onWheel(e) {
    e.preventDefault();
    const delta = e.deltaY > 0 ? -0.1 : 0.1;
    zoomLevel = Math.min(Math.max(+(zoomLevel + delta).toFixed(2), 0.2), 5);
    applyTransform();
  }

    function onMouseDown(e) {
    if (e.button !== 0) return;
    e.preventDefault();
    isPanning = true;
    startX = e.clientX - panX;
    startY = e.clientY - panY;
    container.style.cursor = 'grabbing';
  }

  function onMouseMove(e) {
    if (!isPanning) return;
    panX = e.clientX - startX;
    panY = e.clientY - startY;
    applyTransform();
  }

  function onMouseUp() {
    isPanning = false;
    container.style.cursor = 'grab';
  }

  function changePage(idx) {
    currentPage = idx;
    loadedKey = '';
    loadDiagram();
  }
</script>

<div class="flex flex-col h-full">
  <div class="flex items-center gap-1 bg-gray-800 px-3 py-1 border-b border-gray-700 flex-shrink-0 h-9">
    <button on:click={zoomOut}
      class="bg-gray-700 hover:bg-gray-600 text-white w-7 h-6 rounded font-bold text-base flex items-center justify-center transition-colors"
      title="Diminuir zoom">−</button>
    <button on:click={zoomReset}
      class="bg-gray-700 hover:bg-gray-600 text-white px-2 h-6 rounded text-xs transition-colors"
      title="Ajustar à tela">Fit</button>
    <button on:click={zoomIn}
      class="bg-gray-700 hover:bg-gray-600 text-white w-7 h-6 rounded font-bold text-base flex items-center justify-center transition-colors"
      title="Aumentar zoom">+</button>

    {#if pages.length > 1}
      <div class="w-px h-5 bg-gray-600 mx-2"></div>
      <span class="text-gray-400 text-xs">Página:</span>
      <select value={currentPage} on:change={(e) => changePage(+e.target.value)}
        class="bg-gray-700 text-white text-xs rounded px-2 h-6 border border-gray-600 outline-none cursor-pointer">
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
       <div
      bind:this={container}
      class="flex-1 bg-white overflow-hidden relative"
      style="cursor: grab; user-select: none;"
      draggable="false"
      on:mousedown={onMouseDown}
      on:mousemove={onMouseMove}
      on:mouseup={onMouseUp}
      on:mouseleave={onMouseUp}
      on:wheel={onWheel}
    ></div>
  {/if}
</div>