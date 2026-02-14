<script setup>
import Sidebar from "./components/Sidebar.vue";
import Header from "./components/Header.vue";
import ToastContainer from "./components/ToastContainer.vue";
</script>

<template>
  <ToastContainer />
  <router-view />
</template>
