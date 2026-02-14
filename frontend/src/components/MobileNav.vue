<script setup lang="ts">
import { IconMenu2 } from "@tabler/icons-vue";

defineEmits(["toggle-menu"]);
</script>

<template>
  <header
    class="h-16 bg-surface border-b border-border fixed top-0 left-0 right-0 z-20 md:hidden flex items-center justify-between px-4"
  >
    <!-- Hamburger -->
    <button @click="$emit('toggle-menu')" class="text-muted hover:text-primary">
      <IconMenu2 :size="24" />
    </button>

    <!-- Logo -->
    <div class="flex items-center gap-2">
      <img
        src="/assets/header-logo.webp"
        alt="Allocra"
        class="h-6 w-auto brightness-0 invert"
      />
    </div>

    <!-- Status indicator -->
    <div
      class="w-2 h-2 rounded-full bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.5)]"
    ></div>
  </header>
</template>
