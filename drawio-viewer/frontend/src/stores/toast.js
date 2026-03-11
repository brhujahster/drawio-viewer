import { writable } from 'svelte/store';

export const toasts = writable([]);

export function showToast(message, type = 'error', duration = 4000) {
  const id = Date.now() + Math.random();
  toasts.update(t => [...t, { id, message, type }]);
  setTimeout(() => {
    toasts.update(t => t.filter(x => x.id !== id));
  }, duration);
}