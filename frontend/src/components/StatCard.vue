<script setup lang="ts">
import { type FunctionalComponent } from "vue";

defineProps<{
  label: string;
  value: string | number;
  trend?: number;
  icon?: FunctionalComponent;
}>();
</script>

<template>
  <div
    class="bg-surface border border-border p-5 rounded-sm flex flex-col justify-between group hover:border-accent/50 transition-colors"
  >
    <div class="flex items-center justify-between">
      <span class="text-xs font-mono text-muted uppercase tracking-wider">{{
        label
      }}</span>
      <component
        v-if="icon"
        :is="icon"
        :size="16"
        class="text-muted group-hover:text-accent transition-colors"
      />
    </div>
    <div class="mt-4 flex items-end justify-between">
      <div class="text-2xl font-bold text-primary tracking-tight font-mono">
        {{ value }}
      </div>
      <div
        v-if="trend !== undefined"
        :class="[trend >= 0 ? 'text-green-500' : 'text-red-500']"
        class="text-[10px] font-mono border border-current/20 px-1.5 py-0.5 rounded-sm bg-current/5"
      >
        {{ trend >= 0 ? "+" : "" }}{{ trend }}%
      </div>
    </div>
  </div>
</template>
