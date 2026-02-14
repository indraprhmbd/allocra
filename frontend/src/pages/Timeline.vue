<script setup lang="ts">
import { ref } from "vue";
import TimelineGrid from "../components/TimelineGrid.vue";
import { IconFilter, IconCalendarEvent } from "@tabler/icons-vue";

const selectedDate = ref<string>(new Date().toISOString().split("T")[0] || "");
</script>

<template>
  <div class="space-y-6 md:space-y-8 h-full flex flex-col">
    <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-primary tracking-tight">
          Temporal Grid
        </h2>
        <p class="text-muted text-sm mt-1">
          Holistic visualization of engine scheduling and temporal constraints.
        </p>
      </div>

      <div class="flex items-center gap-3 self-start md:self-auto">
        <div
          class="bg-surface border border-border rounded-sm px-3 py-1.5 flex items-center gap-2 text-xs font-mono text-muted relative"
        >
          <IconCalendarEvent :size="14" />
          <input
            v-model="selectedDate"
            type="date"
            class="bg-transparent border-none outline-none text-muted focus:text-primary transition-colors cursor-pointer"
          />
        </div>
        <button
          class="border border-border p-2 rounded-sm text-muted hover:text-primary transition-colors"
        >
          <IconFilter :size="18" />
        </button>
      </div>
    </div>

    <div
      class="flex-1 bg-surface border border-border rounded-sm overflow-hidden min-h-[500px] flex flex-col"
    >
      <div class="flex-1 overflow-x-auto overflow-y-hidden">
        <div class="min-w-[800px] h-full">
          <TimelineGrid :date="selectedDate" />
        </div>
      </div>
    </div>

    <p class="text-xs text-muted text-center md:hidden pb-4">
      Swipe horizontally to view full timeline
    </p>
  </div>
</template>
