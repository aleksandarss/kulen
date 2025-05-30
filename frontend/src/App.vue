<template>
  <div class="min-h-screen bg-gray-100">
    <NavBar />
    <div class="p-6">
      <RouterView />
    </div>
  </div>
</template>

<script setup>
import NavBar from './components/NavBar.vue'
import { RouterView } from 'vue-router'
</script>
