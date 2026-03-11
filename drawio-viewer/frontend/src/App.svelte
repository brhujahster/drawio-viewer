<script>
  import TopBar from './components/TopBar.svelte';
  import TabBar from './components/TabBar.svelte';
  import DiagramCanvas from './components/DiagramCanvas.svelte';
  import Toast from './components/Toast.svelte';
  import { diagrams, activeTabId } from './stores/diagrams.js';

  $: activeDiagram = $diagrams.find((d) => d.id === $activeTabId) ?? null;
</script>

<div class="flex flex-col h-screen bg-gray-900 text-gray-100 overflow-hidden">
  <TopBar />
  <TabBar />

  <main class="flex-1 overflow-hidden">
    {#if activeDiagram}
      <DiagramCanvas xmlPath={activeDiagram.xmlPath} />
    {:else}
      <div class="h-full flex flex-col items-center justify-center gap-4 text-gray-600 select-none">
        <svg xmlns="http://www.w3.org/2000/svg" class="w-16 h-16 opacity-30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1"
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p class="text-sm">Abra um arquivo <code class="bg-gray-800 px-1 py-0.5 rounded text-gray-400">.drawio</code> para começar</p>
      </div>
    {/if}
  </main>

  <Toast />
</div>