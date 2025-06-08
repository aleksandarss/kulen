<template>
  <div class="overflow-x-auto">
    <table class="min-w-full border text-sm text-left rounded-xl overflow-hidden shadow">
      <thead>
        <tr>
          <th class="px-3 py-2 bg-muted text-secondary border-r">Meal</th>
          <th v-for="day in days" :key="day" class="px-3 py-2 bg-muted text-secondary border-r">
            {{ day }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="meal in meals" :key="meal" class="border-t">
          <td class="px-3 py-2 font-semibold text-primary border-r">
            {{ capitalize(meal) }}
          </td>
          <td v-for="day in days" :key="day + meal" class="px-3 py-2 border-r">
            <div class="cursor-pointer hover:bg-accent/10" @click="() => handleSelect(day, meal)">
              <div class="text-sm text-secondary">
                {{ getEntry(day, meal)?.Recipe?.Title || 'Add' }}
              </div>
              <ul v-if="getEntry(day, meal)?.Extras?.length" class="pl-2 text-xs text-muted">
                <li v-for="extra in getEntry(day, meal).Extras" :key="extra.ID">- {{ extra.Name }}</li>
              </ul>
            </div>
            <div v-if="getEntry(day, meal)" class="mt-1">
              <button @click.stop="() => handleExtras(getEntry(day, meal))" class="text-xs text-accent mr-2">+ Extra</button>
              <button
                v-for="extra in getEntry(day, meal).Extras"
                :key="extra.ID"
                @click.stop="() => handleRemoveExtra(extra.ID)"
                class="text-xs text-red-500 mr-1"
              >
                ❌ {{ extra.Name }}
              </button>
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
