<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import Menu from './components/Menu.vue'
</script>

<template>
  <div class="m-2">
    <Menu class="mb-2" />
    <RouterView />
  </div>
</template>

<style scoped>
</style>
