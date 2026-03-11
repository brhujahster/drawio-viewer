<script>
  import { diagrams, activeTabId } from '../stores/diagrams.js';

  function selectTab(id) {
    activeTabId.set(id);
  }

  function closeTab(diagram) {
    let currentList;
    const unsub = diagrams.subscribe(d => (currentList = d));
    unsub();

    const idx = currentList.findIndex(d => d.id === diagram.id);
    const newList = currentList.filter(d => d.id !== diagram.id);

    diagrams.set(newList);

    activeTabId.update(current => {
      if (current !== diagram.id) return current;
      if (newList.length === 0) return null;
      return newList[Math.min(idx, newList.length - 1)].id;
    });

    if (diagram.isTemp) {
      window['go']['main']['App']['DeleteTempFile'](diagram.xmlPath).catch(() => {});
    }
  }
</script>

{#if $diagrams.length > 0}
  <div class="flex bg-gray-900 border-b border-gray-700 overflow-x-auto flex-shrink-0">
    {#each $diagrams as diagram (diagram.id)}
      <div
        class="flex items-center gap-1 px-4 py-2 text-sm cursor-pointer border-r border-gray-700 whitespace-nowrap select-none transition-colors
          {$activeTabId === diagram.id
            ? 'bg-gray-700 text-white border-t-2 border-t-blue-500'
            : 'text-gray-400 hover:bg-gray-800 hover:text-gray-200'}"
        on:click={() => selectTab(diagram.id)}
      >
        <span>{diagram.name}</span>
        <button
          on:click|stopPropagation={() => closeTab(diagram)}
          class="ml-2 text-gray-500 hover:text-white text-xs w-4 h-4 flex items-center justify-center rounded hover:bg-gray-600 transition-colors"
        >
          ✕
        </button>
      </div>
    {/each}
  </div>
{/if}