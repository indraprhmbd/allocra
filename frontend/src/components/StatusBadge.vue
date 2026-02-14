<script setup lang="ts">
const props = defineProps<{
  status: string;
  variant?: "default" | "success" | "warning" | "error" | "info";
}>();

const variants = {
  default: "text-muted border-border bg-border/10",
  success: "text-green-500 border-green-500/20 bg-green-500/5",
  warning: "text-yellow-500 border-yellow-500/20 bg-yellow-500/5",
  error: "text-red-500 border-red-500/20 bg-red-500/5",
  info: "text-accent border-accent/20 bg-accent/5",
};
</script>

<template>
  <span
    class="text-[10px] font-mono px-1.5 py-0.5 border rounded-sm uppercase tracking-wider inline-flex items-center"
    :class="variants[variant || 'default']"
  >
    <slot>{{ status }}</slot>
  </span>
</template>
