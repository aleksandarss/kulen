<template>
  <div class="overflow-x-auto">
    <table class="min-w-full border text-base text-left rounded-xl overflow-hidden shadow-lg">
      <thead>
        <tr>
          <th class="px-4 py-3 bg-muted text-secondary border-r text-lg">Meal</th>
          <th v-for="day in days" :key="day" class="px-4 py-3 bg-muted text-secondary border-r text-lg">
            {{ day }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="meal in meals" :key="meal" class="border-t">
          <td class="px-4 py-3 font-semibold text-primary border-r bg-muted/10">
            {{ capitalize(meal) }}
          </td>
          <td v-for="day in days" :key="day + meal" class="px-4 py-3 border-r align-top">
            <div class="cursor-pointer hover:bg-accent/20 p-2 rounded transition-colors" @click="() => handleSelect(day, meal)">
              <div class="text-base font-medium text-secondary">
                {{ getEntry(day, meal)?.Recipe?.Title || 'Add' }}
              </div>
              <ul v-if="getEntry(day, meal)?.Extras?.length" class="pl-3 text-sm text-muted mt-1 list-disc">
                <li v-for="extra in getEntry(day, meal).Extras" :key="extra.ID">{{ extra.Name }}</li>
              </ul>
            </div>

            <!-- extras buttons -->
            <div v-if="getEntry(day, meal)" class="mt-2 flex flex-wrap gap-2">
              <button @click.stop="() => handleExtras(getEntry(day, meal))" class="text-sm text-accent underline">
                + Extra
              </button>
              <button
                v-for="extra in getEntry(day, meal).Extras"
                :key="extra.ID"
                @click.stop="() => handleRemoveExtra(extra.ID)"
                class="text-sm text-red-500 hover:underline"
              >
                ❌ {{ extra.Name }}
              </button>
            </div>

            <!-- view recipe link -->
            <div v-if="getEntry(day, meal)?.Recipe" class="mt-2">
              <router-link
                :to="`/recipes/${getEntry(day, meal).Recipe.ID}`"
                class="inline-block px-2 py-1 text-xs font-medium rounded bg-blue-100 text-blue-800 hover:bg-blue-200 transition-colors"
                @click.stop
              >
                🔍 View Recipe
              </router-link>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>


<script setup lang="ts">
import { defineProps } from 'vue'

const props = defineProps<{
  entries: any[]
  onSelect: (day: string, meal: string) => void
  onExtras: (entry: any) => void
  onRemoveExtra: (extraId: number) => void
}>()

const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday']
const meals = ['breakfast', 'lunch', 'dinner']

function getEntry(day: string, meal: string) {
  return props.entries.find((e) => e.Day === day && e.MealType === meal)
}

function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

function handleSelect(day: string, meal: string) {
  props.onSelect(day, meal)
}

function handleExtras(entry: any) {
  props.onExtras(entry)
}

function handleRemoveExtra(extraId: number) {
  props.onRemoveExtra(extraId)
}
</script>

<style scoped>
/* optional styles */
</style>
