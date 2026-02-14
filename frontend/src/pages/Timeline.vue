<script setup lang="ts">
import TimelineGrid from "../components/TimelineGrid.vue";
import { IconFilter, IconCalendarEvent } from "@tabler/icons-vue";
</script>

<template>
  <div class="space-y-8 h-full flex flex-col">
    <div class="flex items-end justify-between">
      <div>
        <h2 class="text-2xl font-bold text-primary tracking-tight">
          Temporal Grid
        </h2>
        <p class="text-muted text-sm mt-1">
          Holistic visualization of engine scheduling and temporal constraints.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <div
          class="bg-surface border border-border rounded-sm px-3 py-1.5 flex items-center gap-2 text-xs font-mono text-muted"
        >
          <IconCalendarEvent :size="14" />
          2026-02-13
        </div>
        <button
          class="border border-border p-2 rounded-sm text-muted hover:text-primary transition-colors"
        >
          <IconFilter :size="18" />
        </button>
      </div>
    </div>

    <div
      class="flex-1 bg-surface border border-border rounded-sm overflow-hidden min-h-[500px]"
    >
      <TimelineGrid />
    </div>
  </div>
</template>
