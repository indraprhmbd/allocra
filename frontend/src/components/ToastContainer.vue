<script setup lang="ts">
import { useToast } from "../composables/useToast";
import {
  IconCheck,
  IconX,
  IconInfoCircle,
  IconAlertTriangle,
} from "@tabler/icons-vue";

const { toasts, removeToast } = useToast();

const icons = {
  success: IconCheck,
  error: IconX,
  info: IconInfoCircle,
  warning: IconAlertTriangle,
};

const variants = {
  success: "border-green-500/50 bg-green-500/5 text-green-400",
  error: "border-red-500/50 bg-red-500/5 text-red-400",
  info: "border-accent/50 bg-accent/5 text-accent",
  warning: "border-yellow-500/50 bg-yellow-500/5 text-yellow-400",
};
</script>

<template>
  <div
    class="fixed bottom-6 right-6 z-[100] flex flex-col gap-3 pointer-events-none"
  >
    <TransitionGroup name="toast">
      <div
        v-for="t in toasts"
        :key="t.id"
        :class="[
          'pointer-events-auto flex items-center gap-3 px-4 py-3 rounded-sm border backdrop-blur-md shadow-lg min-w-[240px]',
          variants[t.variant],
        ]"
      >
        <component :is="icons[t.variant]" :size="18" class="shrink-0" />
        <span class="text-xs font-medium uppercase tracking-wider">{{
          t.message
        }}</span>
        <button
          @click="removeToast(t.id)"
          class="ml-auto p-1 hover:bg-white/10 rounded-full transition-colors"
        >
          <IconX :size="14" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}
.toast-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
