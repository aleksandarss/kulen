import { ref } from 'vue'

const isLoggedIn = ref(!!localStorage.getItem('token'))

function logout() {
  localStorage.removeItem('token')
  localStorage.removeItem('refresh_token')
  isLoggedIn.value = false
}

export function useAuth() {
  return {
    isLoggedIn,
    logout,
  }
}