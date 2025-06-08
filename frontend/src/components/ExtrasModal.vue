<template>
  <div v-if="show && entry" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-white rounded-xl p-6 max-w-md w-full shadow-lg">
      <h2 class="text-lg font-semibold text-primary mb-4">
        Extras for {{ entry?.Recipe?.Title || 'Meal' }}
      </h2>

      <ul class="mb-4">
        <li v-for="extra in entry?.Extras || []" :key="extra.ID" class="flex justify-between items-center mb-2">
          <span>{{ extra.Name }}</span>
          <button
            @click="remove(extra.ID)"
            class="text-red-500 hover:underline text-sm"
          >
            Remove
          </button>
        </li>
      </ul>

      <div class="flex gap-2 mb-4">
        <input
          v-model="newExtra"
          type="text"
          placeholder="New extra"
          class="border rounded px-2 py-1 flex-1"
        />
        <button
          @click="add"
          class="bg-accent text-white px-3 py-1 rounded hover:bg-accent/90"
        >
          Add
        </button>
      </div>

      <div class="text-right">
        <button @click="$emit('close')" class="text-secondary hover:underline text-sm">Close</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import api from '../api'

const props = defineProps<{
  show: boolean
  entry: any | null
}>()

const emits = defineEmits(['close', 'updated'])

const newExtra = ref('')

async function add() {
  if (!props.entry || !newExtra.value.trim()) return

  try {
    const token = localStorage.getItem('token')
    if (!token) return

    const payload = {
      recipe_id: props.entry.RecipeID,
      day: props.entry.Day,
      meal_type: props.entry.MealType,
      extras: [...(props.entry.Extras?.map(e => e.Name) || []), newExtra.value],
    }

    await api.post('/menu', payload, {
      headers: { Authorization: `Bearer ${token}` },
    })

    newExtra.value = ''
    emits('updated')
    emits('close')
  } catch (err) {
    console.error('Failed to add extra:', err)
  }
}


async function remove(extraId: number) {
  if (!props.entry) return
  try {
    const token = localStorage.getItem('token')
    await api.delete(`/menu/${props.entry.ID}/extras/${extraId}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    emits('updated')
  } catch (err) {
    console.error('Failed to remove extra:', err)
  }
}
</script>
