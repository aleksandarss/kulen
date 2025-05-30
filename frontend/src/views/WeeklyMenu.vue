<template>
  <div class="p-4 max-w-5xl mx-auto">
    <h1 class="text-xl font-bold text-primary mb-6">Weekly Menu</h1>

    <!-- Mobile -->
    <MenuDayMobile
      v-if="isMobile"
      :entries="entries"
      :currentDayIndex="currentDayIndex"
      :onDayChange="(i) => (currentDayIndex = i)"
      :onSelect="openSelector"
      :onRemove="removeEntry"
      :onNavigate="navigateToRecipe"
    />

    <!-- Desktop -->
    <MenuGrid
      v-else
      :entries="entries"
      :onSelect="openSelector"
      :onRemove="removeEntry"
      :onNavigate="navigateToRecipe"
    />

    <MenuSelectorModal
      :show="modalOpen"
      :day="selectedDay"
      :meal="selectedMeal"
      :entries="entries"
      @close="modalOpen = false"
      @updated="loadMenu"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import MenuGrid from '../components/MenuGrid.vue'
import MenuDayMobile from '../components/MenuDayMobile.vue'
import MenuSelectorModal from '../components/MenuSelectorModal.vue'
import api from '../api'

const entries = ref<any[]>([])
const currentDayIndex = ref(0)
const isMobile = ref(window.innerWidth < 768)
const modalOpen = ref(false)
const selectedDay = ref('')
const selectedMeal = ref('')
const router = useRouter()

function openSelector(day: string, meal: string) {
  selectedDay.value = day
  selectedMeal.value = meal
  modalOpen.value = true
}

async function loadMenu() {
  const token = localStorage.getItem('token')
  if (!token) {
    entries.value = [] // or you can just return
    return
  }

  try {
    const res = await api.get('/menu')
    entries.value = res.data
  } catch (err) {
    console.error('Failed to load menu:', err)
  }
}


async function removeEntry(entryId: number) {
  try {
    await api.delete(`/menu/${entryId}`)
    await loadMenu()
  } catch (err) {
    console.error('Failed to remove entry:', err)
  }
}

function navigateToRecipe(recipeId: number) {
  router.push(`/recipes/${recipeId}`)
}

onMounted(loadMenu)
</script>

<style scoped>
/****** You can style remove buttons or recipe links here if needed *******/
</style>
