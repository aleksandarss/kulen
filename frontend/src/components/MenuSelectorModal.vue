<template>
  <div v-if="show" class="fixed inset-0 z-50 bg-black/40 flex items-center justify-center">
    <div class="bg-white rounded-xl shadow-xl p-6 w-full max-w-md animate-in fade-in zoom-in duration-200">
      <h2 class="text-lg font-bold text-primary mb-4">
        Select recipe for {{ day }} — {{ capitalize(meal) }}
      </h2>

      <div class="space-y-2 max-h-64 overflow-y-auto mb-4">
        <div
          v-for="recipe in recipes"
          :key="recipe.ID"
          @click="selectRecipe(recipe.ID)"
          class="cursor-pointer px-4 py-2 rounded hover:bg-accent/10 border border-transparent hover:border-accent text-sm text-primary"
        >
          {{ recipe.Title }}
        </div>
      </div>

      <div class="flex justify-end gap-2">
        <button @click="close" class="text-sm text-secondary hover:underline">Cancel</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, defineProps, defineEmits } from 'vue'
import api from '../api'

const props = defineProps<{
  show: boolean
  day: string
  meal: string
  entries: Array<{
    id: number
    day: string
    meal_type: string
    recipe_id: number
  }>
}>()

const emit = defineEmits(['close', 'updated'])

const recipes = ref<any[]>([])

watch(() => props.show, async (val) => {
  if (val) {
    try {
      const res = await api.get('/recipes')
      recipes.value = res.data
    } catch (err) {
      console.error('Failed to load recipes:', err)
    }
  }
})

async function selectRecipe(recipeId: number) {
  try {
    const normalizedDay = (props.day || '').toLowerCase().trim()
    const normalizedMeal = (props.meal || '').toLowerCase().trim()

    const existingEntry = (props.entries || []).find((e) => {
      const entryDay = (e.day || e.Day || '').toLowerCase().trim()
      const entryMeal = (e.meal_type || e.MealType || '').toLowerCase().trim()
      return entryDay === normalizedDay && entryMeal === normalizedMeal
    })

    console.debug('existingEntry:', existingEntry)

    if (existingEntry?.ID) {
      console.debug('Deleting existing entry.')
      await api.delete(`/menu/${existingEntry.ID}`)
    }

    await api.post('/menu', {
      user_id: 1,
      day: props.day,
      meal_type: props.meal,
      recipe_id: recipeId,
    })

    emit('updated')
    emit('close')
  } catch (err) {
    console.error('Failed to update menu entry:', err)
  }
}

function close() {
  emit('close')
}

function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1)
}
</script>

<style scoped>
/* Optional: add styles if needed */
</style>
