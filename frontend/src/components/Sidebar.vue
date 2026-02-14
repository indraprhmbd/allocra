<script setup lang="ts">
import {
  IconLayoutDashboard,
  IconDatabase,
  IconArrowsExchange,
  IconTimeline,
  IconFlask,
  IconX,
} from "@tabler/icons-vue";

defineProps<{
  mobileOpen?: boolean;
}>();

const emit = defineEmits(["close"]);

const menuItems = [
  { name: "Dashboard", path: "/", icon: IconLayoutDashboard },
  { name: "Resources", path: "/resources", icon: IconDatabase },
  { name: "Allocations", path: "/allocations", icon: IconArrowsExchange },
  { name: "Timeline", path: "/timeline", icon: IconTimeline },
  { name: "Playground", path: "/playground", icon: IconFlask },
];
</script>

<template>
  <div>
    <!-- Mobile Overlay -->
    <div
      v-if="mobileOpen"
      @click="emit('close')"
      class="fixed inset-0 bg-black/50 z-40 md:hidden backdrop-blur-sm"
    ></div>

    <!-- Sidebar -->
    <aside
      class="w-64 bg-surface border-r border-border h-screen flex flex-col fixed left-0 top-0 z-50 transition-transform duration-300 md:translate-x-0"
      :class="[
        mobileOpen ? 'translate-x-0' : '-translate-x-full mb:translate-x-0',
        'md:block',
      ]"
    >
      <div class="p-6 border-b border-border flex items-center justify-between">
        <div class="flex items-center gap-3">
          <img
            src="/assets/header-logo.webp"
            alt="Allocra"
            class="h-8 w-auto brightness-0 invert"
          />
          <span
            class="text-[10px] text-muted font-mono bg-border px-1 rounded-sm ml-auto"
            >V0.1.0</span
          >
        </div>

        <!-- Mobile Close Button -->
        <button
          @click="emit('close')"
          class="text-muted hover:text-primary md:hidden"
        >
          <IconX :size="20" />
        </button>
      </div>

      <nav class="flex-1 p-4 space-y-1">
        <router-link
          v-for="item in menuItems"
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-3 py-2 text-sm font-medium transition-colors rounded-sm group relative"
          :class="[
            $route.path === item.path
              ? 'bg-accent/10 text-accent'
              : 'text-muted hover:bg-border hover:text-primary',
          ]"
          @click="emit('close')"
        >
          <component :is="item.icon" :size="18" />
          {{ item.name }}
          <div
            v-if="$route.path === item.path"
            class="absolute left-0 top-0 bottom-0 w-1 bg-accent rounded-sm"
          ></div>
        </router-link>
      </nav>
      <div class="p-4 border-t border-border mt-auto">
        <div
          class="flex items-center gap-3 px-3 py-2 text-xs text-muted font-mono uppercase tracking-wider"
        >
          System Status:
          <span
            class="text-green-500 underline underline-offset-4 decoration-current"
            >NOMINAL</span
          >
        </div>
      </div>
    </aside>
  </div>
</template>
