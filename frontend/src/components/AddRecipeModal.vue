<template>
  <transition name="fade-scale">
    <div
      v-if="visible"
      class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50"
      @click.self="emit('close')"
    >
      <div
        class="bg-white w-full max-w-lg max-h-screen overflow-y-auto p-6 rounded-lg shadow-lg relative"
        style="scrollbar-gutter: stable"
      >
        <h2 class="text-xl font-semibold text-primary mb-4">Add New Recipe</h2>

        <form @submit.prevent="submit">
          <!-- Title -->
          <div class="mb-4">
            <label class="block text-sm text-primary mb-1">Title</label>
            <input
              v-model="title"
              type="text"
              class="w-full border border-secondary rounded px-3 py-2"
              required
            />
          </div>

          <!-- Instructions -->
          <div class="mb-4">
            <label class="block text-sm text-primary mb-1">Instructions</label>
            <textarea
              v-model="instructions"
              class="w-full border border-secondary rounded px-3 py-2"
              rows="4"
              required
            ></textarea>
          </div>

          <!-- Ingredients -->
          <div class="mb-4">
            <label class="block text-sm text-primary mb-2">Ingredients</label>
            <div
              v-for="(ingredient, index) in ingredients"
              :key="index"
              class="grid grid-cols-3 gap-2 mb-2"
            >
              <IngredientInput v-model="ingredient.name" />
              <input
                v-model="ingredient.amount"
                type="text"
                placeholder="Amount"
                class="border border-secondary rounded px-2 py-1"
                required
              />
              <input
                v-model="ingredient.unit"
                type="text"
                placeholder="Unit"
                class="border border-secondary rounded px-2 py-1"
                required
              />
            </div>
            <button
              type="button"
              @click="addIngredient"
              class="text-sm text-accent hover:underline"
            >
              + Add Ingredient
            </button>
          </div>

          <!-- Tags -->
          <div class="mb-4">
            <label class="block text-sm text-primary mb-2">Tags</label>
            <div v-if="allTags.length" class="flex flex-wrap gap-3">
              <label
                v-for="tag in allTags"
                :key="tag.id"
                class="flex items-center gap-2 text-sm text-primary"
              >
                <input
                  type="checkbox"
                  :value="tag.name"
                  v-model="selectedTags"
                  class="accent-accent"
                />
                {{ tag.name }}
              </label>
            </div>
            <div v-else class="text-sm text-secondary">No tags found.</div>
          </div>

          <!-- Steps -->
          <div class="mb-4">
            <label class="block text-sm text-primary mb-2">Steps</label>
            <div
              v-for="(step, index) in steps"
              :key="index"
              class="mb-4 border border-secondary rounded p-3"
            >
              <input
                v-model="step.title"
                type="text"
                :placeholder="`Step ${index + 1} Title`"
                class="w-full mb-2 border border-secondary rounded px-3 py-2 text-sm"
                required
              />
              <textarea
                v-model="step.text"
                :placeholder="`Step ${index + 1} Instructions`"
                rows="3"
                class="w-full border border-secondary rounded px-3 py-2 text-sm"
                required
              ></textarea>
            </div>

            <button
              type="button"
              @click="addStep"
              class="text-sm text-accent hover:underline"
              :disabled="steps.length >= 10"
            >
              + Add Step
            </button>
            <div v-if="steps.length >= 10" class="text-xs text-secondary mt-1">
              Max 10 steps allowed
            </div>
          </div>

          <!-- Buttons -->
          <div class="flex justify-end gap-3 mt-6">
            <button type="button" @click="emit('close')" class="px-4 py-2 rounded bg-secondary text-primary">
              Cancel
            </button>
            <button type="submit" class="px-4 py-2 rounded bg-accent text-white hover:bg-primary">
              Save
            </button>
          </div>
        </form>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, watch, defineEmits, defineProps, onMounted } from 'vue'
import api from '../api'
import IngredientInput from './IngredientInput.vue'

const emit = defineEmits(['close', 'submit'])
const props = defineProps<{ visible: boolean }>()

const title = ref('')
const instructions = ref('')
const ingredients = ref([{ name: '', amount: '', unit: '' }])
const selectedTags = ref<string[]>([])
const allTags = ref<{ id: number; name: string }[]>([])
const steps = ref([{ title: '', text: '' }])

watch(props, () => {
  if (!props.visible) return
  title.value = ''
  instructions.value = ''
  ingredients.value = [{ name: '', amount: '', unit: '' }]
  selectedTags.value = []
  steps.value = [{ title: '', text: '' }]
})

onMounted(async () => {
  try {
    const res = await api.get('/tags')
    allTags.value = res.data
  } catch (err) {
    console.error('Failed to load tags:', err)
  }
})

function addIngredient() {
  ingredients.value.push({ name: '', amount: '', unit: '' })
}

function addStep() {
  if (steps.value.length < 10) {
    steps.value.push({ title: '', text: '' })
  }
}

async function submit() {
  try {
    await api.post('/recipes', {
      title: title.value,
      instructions: instructions.value,
      ingredients: ingredients.value.filter(i => i.name && i.amount && i.unit),
      tags: selectedTags.value,
      steps: steps.value.filter(s => s.title.trim() && s.text.trim())
    })
    emit('submit')
  } catch (err) {
    console.error('Failed to create recipe:', err)
  } finally {
    emit('close')
  }
}
</script>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.25s ease;
}
.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
