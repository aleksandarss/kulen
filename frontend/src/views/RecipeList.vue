<template>
  <AddRecipeModal
    :visible="showModal"
    @close="showModal = false"
    @submit="handleCreate"
  />

  <EditRecipeModal
    v-if="editingRecipe"
    :recipe="editingRecipe"
    @close="editingRecipe = null"
    @updated="loadRecipes"
  />

  <div class="bg-background min-h-screen p-4">
    <div class="flex justify-between items-center mb-6 pb-2 border-b border-secondary">
      <h1 class="text-xl font-semibold text-primary">Recipes</h1>
      <button
        class="bg-accent text-white px-4 py-2 rounded hover:bg-primary transition-colors"
        @click="showModal = true"
      >
        + Add Recipe
      </button>
    </div>

    <div v-if="loading" class="text-secondary">Loading recipes...</div>
    <div v-else-if="recipes.length === 0" class="text-secondary">No recipes found.</div>

    <div class="grid gap-6 md:grid-cols-2">
      <div
        v-for="recipe in recipes"
        :key="recipe.ID"
        class="bg-white rounded-lg shadow p-4 border border-secondary"
      >
        <router-link
          :to="`/recipes/${recipe.ID}`"
          class="block hover:bg-gray-50 p-1 -m-1 rounded transition"
        >
          <h2 class="text-lg font-semibold text-primary mb-2 underline">
            {{ recipe.Title }}
          </h2>
        </router-link>

        <p class="text-sm text-secondary mb-4">{{ recipe.Instructions }}</p>

        <div class="mt-3">
          <h3 class="text-sm font-semibold text-primary mb-1">Ingredients:</h3>
          <ul class="flex flex-wrap gap-2 text-sm text-secondary">
            <li
              v-for="ingredient in recipe.Ingredients"
              :key="ingredient.ID"
              class="bg-accent/10 text-accent px-2 py-1 rounded-md"
            >
              {{ ingredient.Ingredient?.Name }} ({{ ingredient.Amount }} {{ ingredient.Unit }})
            </li>
          </ul>
        </div>

        <div v-if="recipe.Tags?.length" class="mt-3">
          <span
            v-for="tag in recipe.Tags"
            :key="tag.ID"
            class="inline-block bg-secondary text-primary text-xs px-2 py-1 rounded mr-2"
          >
            {{ tag.Tag?.Name }}
          </span>
        </div>

        <button
          class="text-sm text-accent hover:underline mt-2"
          @click="startEdit(recipe)"
        >
          Edit
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '../api'
import AddRecipeModal from '../components/AddRecipeModal.vue'
import EditRecipeModal from '../components/EditRecipeModal.vue'

const recipes = ref([])
const loading = ref(true)
const showModal = ref(false)

onMounted(async () => {
  await fetchRecipes()
})

async function fetchRecipes() {
  try {
    const res = await api.get('/recipes')
    recipes.value = res.data
  } catch (err) {
    console.error('Error loading recipes:', err)
  } finally {
    loading.value = false
  }
}

async function handleCreate() {
  await fetchRecipes()
}

const editingRecipe = ref(null)

function startEdit(recipe: any) {
  editingRecipe.value = recipe
}
</script>
