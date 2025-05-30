<template>
  <div class="h-screen w-full flex flex-col items-center justify-between p-6 bg-background relative">
    <!-- Back Button -->
    <router-link
      :to="`/recipes/${route.params.id}`"
      class="absolute top-4 left-4 text-accent text-xl hover:text-primary transition"
      aria-label="Back to recipe"
    >
      &larr;
    </router-link>

    <div class="text-sm text-secondary mt-2">
      Step {{ currentIndex + 1 }} of {{ steps.length }}
    </div>

    <div class="flex-1 flex flex-col justify-center items-center w-full max-w-xl overflow-y-auto px-4 py-6">
      <h2 class="text-xl font-bold text-primary text-center mb-4">
        {{ steps[currentIndex]?.Title }}
      </h2>
      <p class="text-base text-primary text-center whitespace-pre-line break-words">
        {{ steps[currentIndex]?.Text }}
      </p>
    </div>

    <div class="flex justify-between items-center gap-4 w-full max-w-md px-2 mt-4">
      <button
        @click="prevStep"
        :disabled="currentIndex === 0"
        class="text-sm text-accent disabled:opacity-30"
      >
        &larr; Back
      </button>
      <button
        @click="nextStep"
        :disabled="currentIndex >= steps.length - 1"
        class="text-sm text-accent disabled:opacity-30"
      >
        Next &rarr;
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '../api'

const route = useRoute()
const steps = ref<{ Title: string; Text: string }[]>([])
const currentIndex = ref(0)

async function loadSteps() {
  try {
    const res = await api.get(`/recipes/${route.params.id}`)
    steps.value = res.data.Steps || []
  } catch (err) {
    console.error('Failed to load recipe steps:', err)
  }
}

function nextStep() {
  if (currentIndex.value < steps.value.length - 1) currentIndex.value++
}

function prevStep() {
  if (currentIndex.value > 0) currentIndex.value--
}

onMounted(loadSteps)
</script>
