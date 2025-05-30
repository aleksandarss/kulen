import { createRouter, createWebHistory } from 'vue-router'
import RecipeList from '../views/RecipeList.vue'
import WeeklyMenu from '../views/WeeklyMenu.vue'
import ShoppingList from '../views/ShoppingList.vue'
import RecipeStep from '@/views/RecipeStep.vue'
import RecipeSingleView from '@/views/RecipeSingleView.vue'

const routes = [
  { path: '/', redirect: '/recipes' },
  { path: '/recipes', component: RecipeList },
  { path: '/menu', component: WeeklyMenu },
  { path: '/shopping-list', component: ShoppingList },
  { path: '/recipes/:id/steps', component: RecipeStep},
  { path: '/recipes/:id', name: 'RecipeView', component: RecipeSingleView }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
