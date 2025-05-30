<template>
  <nav class="bg-primary text-white px-4 py-3 flex gap-3 text-sm font-medium shadow-md">
    <RouterLink
      to="/recipes"
      class="px-3 py-1.5 rounded bg-accent hover:bg-secondary text-white transition-colors"
    >
      Recipes
    </RouterLink>
    <RouterLink
      to="/menu"
      class="px-3 py-1.5 rounded bg-accent hover:bg-secondary text-white transition-colors"
    >
      Weekly Menu
    </RouterLink>
    <RouterLink
      to="/shopping-list"
      class="px-3 py-1.5 rounded bg-accent hover:bg-secondary text-white transition-colors"
    >
      Shopping List
    </RouterLink>
  </nav>
</template>

<script setup>
import { RouterLink } from 'vue-router'
</script>
