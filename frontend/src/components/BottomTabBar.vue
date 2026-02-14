<script setup lang="ts">
import {
  IconLayoutDashboard,
  IconDatabase,
  IconArrowsExchange,
  IconTimeline,
  IconFlask,
} from "@tabler/icons-vue";

const menuItems = [
  { name: "Dash", path: "/", icon: IconLayoutDashboard },
  { name: "Nodes", path: "/resources", icon: IconDatabase },
  { name: "Alloc", path: "/allocations", icon: IconArrowsExchange },
  { name: "Time", path: "/timeline", icon: IconTimeline },
  { name: "Play", path: "/playground", icon: IconFlask },
];
</script>

<template>
  <nav
    class="fixed bottom-0 left-0 right-0 h-16 bg-surface border-t border-border z-20 md:hidden pb-safe"
  >
    <div class="flex items-center justify-around h-full">
      <router-link
        v-for="item in menuItems"
        :key="item.path"
        :to="item.path"
        class="flex flex-col items-center gap-1 px-2 py-1 min-w-[60px] transition-colors"
        :class="[
          $route.path === item.path
            ? 'text-accent'
            : 'text-muted hover:text-primary',
        ]"
      >
        <component
          :is="item.icon"
          :size="20"
          :stroke-width="$route.path === item.path ? 2.5 : 1.5"
        />
        <span class="text-[10px] font-medium">{{ item.name }}</span>
      </router-link>
    </div>
  </nav>
</template>

<style scoped>
.pb-safe {
  padding-bottom: env(safe-area-inset-bottom);
}
</style>
