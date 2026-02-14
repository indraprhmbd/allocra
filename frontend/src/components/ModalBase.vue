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
      class="fixed inset-0 z-50 flex items-center justify-center p-4"
    >
      <!-- Backdrop -->
      <div
        class="absolute inset-0 bg-background/80 backdrop-blur-sm"
        @click="emit('close')"
      ></div>

      <!-- Modal Content -->
      <div
        class="relative w-full max-w-lg bg-surface border border-border rounded-sm shadow-2xl flex flex-col max-h-[90vh]"
      >
        <div
          class="flex items-center justify-between p-5 border-b border-border"
        >
          <h3 class="text-sm font-bold uppercase tracking-widest text-primary">
            {{ title }}
          </h3>
          <button
            @click="emit('close')"
            class="text-muted hover:text-primary transition-colors"
          >
            <IconX :size="18" />
          </button>
        </div>

        <div class="flex-1 p-6 overflow-y-auto">
          <slot></slot>
        </div>

        <div
          v-if="$slots.actions"
          class="p-4 border-t border-border bg-background/30 flex justify-end gap-3"
        >
          <slot name="actions"></slot>
        </div>
      </div>
    </div>
  </Transition>
</template>
