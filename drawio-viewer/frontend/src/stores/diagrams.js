import { writable } from 'svelte/store';

export const diagrams = writable([]);
export const activeTabId = writable(null);