<script setup lang="ts">
import {
  IconLayoutDashboard,
  IconDatabase,
  IconArrowsExchange,
  IconTimeline,
  IconFlask,
} from "@tabler/icons-vue";

const menuItems = [
  { name: "Dashboard", path: "/", icon: IconLayoutDashboard },
  { name: "Resources", path: "/resources", icon: IconDatabase },
  { name: "Allocations", path: "/allocations", icon: IconArrowsExchange },
  { name: "Timeline", path: "/timeline", icon: IconTimeline },
  { name: "Playground", path: "/playground", icon: IconFlask },
];
</script>

<template>
  <aside
    class="w-64 bg-surface border-r border-border h-screen flex flex-col fixed left-0 top-0"
  >
    <div class="p-6 border-b border-border">
      <h1
        class="text-xl font-bold tracking-tight text-primary flex items-center gap-2"
      >
        <span class="w-2 h-6 bg-accent rounded-sm"></span>
        ALLOCRA
        <span
          class="text-[10px] text-muted font-mono bg-border px-1 rounded-sm ml-auto"
          >V0.1.0</span
        >
      </h1>
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
</template>
