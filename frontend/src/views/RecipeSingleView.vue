<template>
  <div class="p-6">
    <h1 class="text-3xl font-bold text-primary mb-4 text-center">{{ recipe.Title }}</h1>
    <div v-if="recipe.Tags && recipe.Tags.length" class="flex flex-wrap gap-2 justify-center mb-6">
      <span
        v-for="tag in recipe.Tags"
        :key="tag.ID"
        class="bg-secondary text-primary text-xs px-2 py-1 rounded"
      >
        {{ tag.Tag?.Name }}
      </span>
    </div>
    <div class="flex flex-col lg:flex-row gap-6">
        <!-- Left Sidebar (Ingredients + Portions) -->
        <div class="lg:w-1/4 w-full bg-white border border-secondary rounded-xl p-4 shadow">
        <h2 class="text-lg font-semibold text-primary mb-4">Ingredients</h2>

        <div class="mb-4">
            <label class="block text-sm font-medium text-primary mb-1">Portions</label>
            <input
            v-model.number="portions"
            type="number"
            step="1"
            min="1"
            class="w-full border border-secondary rounded px-2 py-1"
            />
        </div>

        <div class="space-y-2">
            <div
                v-for="(ingredient, index) in scaledIngredients"
                :key="index"
                class="flex items-center gap-2"
            >
                <input
                type="checkbox"
                v-model="checkedIngredients"
                :value="index"
                class="accent-accent"
                />
                <span
                :class="[
                    'transition-all duration-200',
                    checkedIngredients.includes(index)
                    ? 'line-through text-secondary opacity-50'
                    : 'text-primary'
                ]"
                >
                {{ ingredient.scaledAmount }} {{ ingredient.Unit }} {{ ingredient.Ingredient.Name }}
                </span>
            </div>
            <router-link
                :to="`/recipes/${recipe.ID}/steps`"
                class="block text-center mt-6 px-4 py-2 bg-accent text-white rounded hover:bg-primary transition"
                >
                Cook Now
            </router-link>
            </div>
        </div>

        <!-- Main Content (Steps) -->
        <div class="flex-1 bg-white border border-secondary rounded-xl p-4 shadow">
        <h1 class="text-2xl font-bold text-primary mb-6">{{ recipe.Title }}</h1>

        <div v-if="recipe.Steps && recipe.Steps.length" class="space-y-6">
            <div
            v-for="(step, index) in recipe.Steps"
            :key="index"
            class="border-b border-gray-200 pb-4"
            >
            <h3 class="font-semibold text-primary">{{ step.Title }}</h3>
            <p class="text-sm text-secondary mt-1">{{ step.Text }}</p>
            </div>
        </div>
        <div v-else class="text-secondary italic">No steps found.</div>
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '../api'

const route = useRoute()
const recipe = ref<any>({ ingredients: [], steps: [] })
const portions = ref(1)

onMounted(async () => {
  try {
    const res = await api.get(`/recipes/${route.params.id}`)
    recipe.value = res.data
  } catch (err) {
    console.error('Failed to load recipe', err)
  }
})

const checkedIngredients = ref<number[]>([])

watch(portions, (val) => {
  if (!val || val < 1) {
    portions.value = 1
  } else {
    portions.value = Math.round(val)
  }
})

function scaleAmount(amount: string): string {
  // support commas as decimal separators
  const num = parseFloat(amount.replace(',', '.'))
  if (isNaN(num)) return amount
  const scaled = num * portions.value
  const result = Number.isInteger(scaled) ? scaled.toString() : scaled.toFixed(2)
  // display using comma for decimals
  return result.replace('.', ',')
}

const scaledIngredients = computed(() =>
  recipe.value.Ingredients?.map((ing: any) => ({
    ...ing,
    scaledAmount: scaleAmount(ing.Amount)
  })) || []
)

</script>

<style scoped>
/* No additional styles */
</style>
