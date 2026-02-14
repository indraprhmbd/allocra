<script setup lang="ts">
import { ref } from "vue";
import Sidebar from "./Sidebar.vue";
import Topbar from "./Topbar.vue";
import MobileNav from "./MobileNav.vue";
import BottomTabBar from "./BottomTabBar.vue";

const mobileMenuOpen = ref(false);
</script>

<template>
  <div class="min-h-screen bg-background">
    <!-- Mobile Navigation -->
    <MobileNav @toggle-menu="mobileMenuOpen = !mobileMenuOpen" />

    <!-- Sidebar (Desktop + Mobile Drawer) -->
    <Sidebar :mobile-open="mobileMenuOpen" @close="mobileMenuOpen = false" />

    <Topbar />

    <main
      class="md:ml-64 pt-16 md:pt-16 pb-20 md:pb-8 p-4 md:p-8 min-h-screen transition-all duration-300"
    >
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Mobile Bottom Tab Bar -->
    <BottomTabBar />
  </div>
</template>
