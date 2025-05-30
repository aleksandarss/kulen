<template>
  <div class="relative" @pointerdown.self="handlePointerDown">
    <input
      v-model="input"
      @input="fetchSuggestions"
      @focus="showSuggestions = true"
      @blur="handleBlur"
      type="text"
      class="w-full border border-secondary rounded px-2 py-1"
      placeholder="Name"
    />

    <ul
      v-if="showSuggestions && suggestions.length > 0"
      class="absolute z-50 bg-white border border-secondary rounded shadow w-full mt-1 max-h-48 overflow-y-auto"
    >
      <li
        v-for="(suggestion, index) in suggestions"
        :key="index"
        class="px-3 py-2 text-sm cursor-pointer hover:bg-accent hover:text-white"
        @pointerdown.prevent="selectSuggestion(suggestion.Name)"
      >
        {{ suggestion.Name }}
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, defineEmits, defineProps } from 'vue'
import api from '../api'

const props = defineProps<{ modelValue: string }>()
const emit = defineEmits(['update:modelValue'])

const input = ref(props.modelValue)
const suggestions = ref<{ Name: string }[]>([])
const showSuggestions = ref(false)

watch(() => props.modelValue, (val) => {
  input.value = val
})

watch(input, (val) => {
  emit('update:modelValue', val)
})

let debounceTimeout: ReturnType<typeof setTimeout>

function fetchSuggestions() {
  clearTimeout(debounceTimeout)
  const q = input.value.trim()
  if (q.length < 2) {
    suggestions.value = []
    showSuggestions.value = false
    return
  }

  debounceTimeout = setTimeout(async () => {
    try {
      const res = await api.get('/ingredients', {
        params: { query: q }
      })
      suggestions.value = res.data
      showSuggestions.value = true
    } catch (err) {
      console.error('Error fetching suggestions:', err)
      suggestions.value = []
      showSuggestions.value = false
    }
  }, 200)
}

function selectSuggestion(name: string) {
  input.value = name
  showSuggestions.value = false
}

function handleBlur() {
  setTimeout(() => {
    showSuggestions.value = false
  }, 150)
}

function handlePointerDown() {
  // Prevent blur from firing before suggestion is selected
}
</script>
