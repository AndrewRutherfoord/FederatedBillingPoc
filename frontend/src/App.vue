<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import Menu from './components/Menu.vue'
</script>

<template>
  <div class="m-2">
    <Menu />
    <RouterView />
  </div>
</template>

<style scoped>
</style>
