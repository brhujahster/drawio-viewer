<script>
  import TopBar from './components/TopBar.svelte';
  import TabBar from './components/TabBar.svelte';
  import DiagramCanvas from './components/DiagramCanvas.svelte';
  import { diagrams, activeTabId } from './stores/diagrams.js';

  $: activeDiagram = $diagrams.find((d) => d.id === $activeTabId) ?? null;
</script>

<div class="flex flex-col h-screen bg-gray-900 text-gray-100 overflow-hidden">
  <TopBar />
  <TabBar />

  {#if activeDiagram}
    <div class="flex-1 overflow-hidden">
      <DiagramCanvas xmlPath={activeDiagram.xmlPath} />
    </div>
  {:else}
    <div class="flex-1 flex items-center justify-center text-gray-500 text-sm">
      <p>Abra um arquivo .drawio para começar</p>
    </div>
  {/if}
</div>