<script setup lang="ts">
import Sidebar from "./Sidebar.vue";
import Topbar from "./Topbar.vue";
</script>

<template>
  <div class="min-h-screen bg-background">
    <Sidebar />
    <Topbar />
    <main class="ml-64 pt-16 p-8 min-h-screen">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>
