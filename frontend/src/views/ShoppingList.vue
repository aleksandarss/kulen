<template>
  <div class="p-4 max-w-xl mx-auto">
    <h1 class="text-xl font-bold text-primary mb-6">Shopping List</h1>

    <div v-if="items.length === 0" class="text-secondary">No items found.</div>

    <ul class="space-y-2">
      <li
        v-for="(item, index) in items"
        :key="index"
        class="flex items-center gap-3 bg-white rounded-lg px-4 py-2 shadow-sm"
      >
        <input
          type="checkbox"
          v-model="item.checked"
          class="accent-accent w-4 h-4"
        />
        <span :class="{ 'line-through text-gray-400': item.checked }">
          {{ item.Name }} — {{ item.Amount }} {{ item.Unit }}
        </span>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import api from '../api'

type ShoppingItem = {
  name: string
  unit: string
  amount: number
  checked?: boolean
}

const items = ref<ShoppingItem[]>([])

onMounted(async () => {
  try {
    const res = await api.get('/shopping-list', { params: { user_id: 1 } })
    items.value = res.data.map((item: ShoppingItem) => ({
      ...item,
      checked: false
    }))
  } catch (err) {
    console.error('Failed to load shopping list:', err)
  }
})
</script>
