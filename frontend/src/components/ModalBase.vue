<script setup lang="ts">
import { onMounted, onUnmounted } from "vue";
import { IconX } from "@tabler/icons-vue";

const props = defineProps<{
  show: boolean;
  title: string;
}>();

const emit = defineEmits(["close"]);

const handleEscape = (e: KeyboardEvent) => {
  if (e.key === "Escape" && props.show) {
    emit("close");
  }
};

onMounted(() => window.addEventListener("keydown", handleEscape));
onUnmounted(() => window.removeEventListener("keydown", handleEscape));
</script>

<template>
  <Transition name="fade">
    <div
      v-if="show"
      class="fixed inset-0 z-50 flex items-end md:items-center justify-center sm:p-4"
    >
      <!-- Backdrop -->
      <div
        class="absolute inset-0 bg-background/80 backdrop-blur-sm"
        @click="emit('close')"
      ></div>

      <!-- Modal Content -->
      <div
        class="relative w-full md:max-w-lg bg-surface border border-border md:rounded-sm shadow-2xl flex flex-col h-full md:h-auto md:max-h-[90vh] animate-slide-up md:animate-zoom"
      >
        <div
          class="flex items-center justify-between p-4 md:p-5 border-b border-border"
        >
          <h3 class="text-sm font-bold uppercase tracking-widest text-primary">
            {{ title }}
          </h3>
          <button
            @click="emit('close')"
            class="text-muted hover:text-primary transition-colors p-1"
          >
            <IconX :size="20" />
          </button>
        </div>

        <div class="flex-1 p-4 md:p-6 overflow-y-auto">
          <slot></slot>
        </div>

        <div
          v-if="$slots.actions"
          class="p-4 border-t border-border bg-background/30 flex justify-end gap-3 sticky bottom-0"
        >
          <slot name="actions"></slot>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.animate-slide-up {
  animation: slideUp 0.3s ease-out;
}
.md\:animate-zoom {
  animation: zoom 0.2s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}

@keyframes zoom {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
