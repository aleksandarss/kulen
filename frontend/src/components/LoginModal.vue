<template>
  <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
    <div class="bg-white rounded-xl shadow-lg p-6 max-w-sm w-full">
      <h2 class="text-lg font-bold mb-4 text-center">Login</h2>

      <div class="space-y-4">
        <input
          type="email"
          v-model="email"
          placeholder="Email"
          class="w-full p-2 border rounded"
        />
        <input
          type="password"
          v-model="password"
          placeholder="Password"
          class="w-full p-2 border rounded"
        />
        <button
          @click="login"
          class="w-full bg-primary text-white p-2 rounded hover:bg-opacity-90"
        >
          Login
        </button>
        <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
      </div>

      <button @click="$emit('close')" class="mt-4 text-sm text-gray-500 hover:underline w-full text-center">
        Cancel
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import api from '../api'

defineProps<{ show: boolean }>()
const emit = defineEmits(['close', 'logged-in'])

const email = ref('')
const password = ref('')
const error = ref('')

async function login() {
  error.value = ''
  try {
    const res = await api.post('/login', {
      email: email.value,
      password: password.value,
    })

    localStorage.setItem('token', res.data.access_token)
    localStorage.setItem('refresh_token', res.data.refresh_token)
    emit('logged-in')
  } catch (err: any) {
    error.value = err?.response?.data?.error || 'Login failed'
  }
}
</script>

<style scoped>
/* Basic styling included; customize as needed */
</style>
