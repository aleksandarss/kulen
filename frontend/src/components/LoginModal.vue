<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
  >
    <div class="bg-white p-6 rounded-xl shadow-lg w-80 space-y-4">
      <h2 class="text-xl font-bold text-gray-800">Login</h2>
      <div class="space-y-2">
        <input
          v-model="email"
          type="email"
          placeholder="Email"
          class="w-full px-3 py-2 border rounded"
        />
        <input
          v-model="password"
          type="password"
          placeholder="Password"
          class="w-full px-3 py-2 border rounded"
        />
      </div>
      <div class="text-red-500 text-sm" v-if="error">{{ error }}</div>
      <button
        @click="handleLogin"
        class="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700"
      >
        Login
      </button>
      <button
        @click="onClose"
        class="w-full text-sm text-gray-500 mt-2 hover:underline"
      >
        Cancel
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from '../api'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits(['close', 'logged-in'])

const email = ref('')
const password = ref('')
const error = ref('')

function onClose() {
  emit('close')
}

async function handleLogin() {
  try {
    const res = await axios.post('/login', {
      email: email.value,
      password: password.value,
    })
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('refresh', res.data.refresh)
    emit('logged-in')
    emit('close')
  } catch {
    error.value = 'Login failed'
  }
}
</script>
