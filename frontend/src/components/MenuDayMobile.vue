<template>
  <div class="p-4">
    <div class="flex justify-between items-center mb-4">
      <button @click="prevDay" class="text-accent">&larr;</button>
      <h2 class="text-lg font-semibold text-primary">{{ currentDay }}</h2>
      <button @click="nextDay" class="text-accent">&rarr;</button>
    </div>

    <div class="space-y-3">
      <div
        v-for="meal in meals"
        :key="meal"
        class="bg-white rounded-xl shadow px-4 py-3"
        @click="onSelect(currentDay, meal)"
      >
        <div class="text-sm text-secondary mb-1">
          {{ capitalize(meal) }}
        </div>
        <div class="text-md text-primary font-medium">
          {{ getEntry(currentDay, meal)?.Recipe?.Title || 'Tap to add' }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, defineProps } from 'vue'

const props = defineProps<{
  entries: any[]
  currentDayIndex: number
  onDayChange: (newIndex: number) => void
  onSelect: (day: string, meal: string) => void
}>()

const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday']
const meals = ['breakfast', 'lunch', 'dinner']

const currentDay = computed(() => days[props.currentDayIndex])

function getEntry(day: string, meal: string) {
  return props.entries.find((e) => e.Day === day && e.MealType === meal)
}

function capitalize(str: string) {
  return str.charAt(0).toUpperCase() + str.slice(1)
}

function prevDay() {
  if (props.currentDayIndex > 0) {
    props.onDayChange(props.currentDayIndex - 1)
  }
}

function nextDay() {
  if (props.currentDayIndex < days.length - 1) {
    props.onDayChange(props.currentDayIndex + 1)
  }
}
</script>
