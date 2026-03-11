<script>
  import { toasts } from '../stores/toast.js';

  const icons = {
    error: '✕',
    success: '✓',
    info: 'ℹ',
  };

  const colors = {
    error: 'bg-red-700 border-red-500',
    success: 'bg-green-700 border-green-500',
    info: 'bg-blue-700 border-blue-500',
  };

  function dismiss(id) {
    toasts.update(t => t.filter(x => x.id !== id));
  }
</script>

{#if $toasts.length > 0}
  <div class="fixed bottom-4 right-4 flex flex-col gap-2 z-50 pointer-events-none">
    {#each $toasts as toast (toast.id)}
      <div
        class="flex items-start gap-3 px-4 py-3 rounded-lg shadow-xl text-sm max-w-sm border pointer-events-auto
          {colors[toast.type] ?? colors.info} text-white"
      >
        <span class="font-bold mt-0.5 flex-shrink-0">{icons[toast.type] ?? icons.info}</span>
        <span class="flex-1 leading-snug">{toast.message}</span>
        <button
          on:click={() => dismiss(toast.id)}
          class="text-white opacity-60 hover:opacity-100 flex-shrink-0 transition-opacity"
        >✕</button>
      </div>
    {/each}
  </div>
{/if}