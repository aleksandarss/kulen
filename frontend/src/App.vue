<template>
  <LoginModal :open="loginOpen" @close="loginOpen = false" @logged-in="fetchUserData()" />

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

import { ref, onMounted } from 'vue'
import LoginModal from '@/components/LoginModal.vue' // adjust path if needed

const loginOpen = ref(false)

function fetchUserData() {
  // fetch menu etc. after login
}

onMounted(() => {
  const token = localStorage.getItem('token')
  if (!token) loginOpen.value = true
})

</script>
