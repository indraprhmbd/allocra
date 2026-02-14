<script setup lang="ts">
import { IconTerminal2, IconWifi } from "@tabler/icons-vue";
</script>

<template>
  <header
    class="h-16 border-b border-border bg-surface flex items-center px-8 fixed top-0 right-0 left-64 z-10"
  >
    <div class="flex items-center gap-4 text-xs font-mono text-muted">
      <div
        class="flex items-center gap-1.5 border border-border px-2 py-1 rounded-sm"
      >
        <IconWifi :size="14" class="text-accent" />
        NETWORK: NODE-01
      </div>
      <div
        class="flex items-center gap-1.5 border border-border px-2 py-1 rounded-sm"
      >
        <IconTerminal2 :size="14" />
        REPL_READY
      </div>
    </div>

    <div class="ml-auto flex items-center gap-6">
      <button
        class="text-xs font-mono text-accent hover:underline decoration-accent underline-offset-4"
      >
        FORCE_SYNC
      </button>
      <div
        class="w-8 h-8 rounded-sm bg-border border border-border flex items-center justify-center text-[10px] text-primary font-mono"
      >
        AD
      </div>
    </div>
  </header>
</template>
