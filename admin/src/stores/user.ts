import { defineStore } from 'pinia'
import { ref } from 'vue'
// import { useRouter } from 'vue-router'


export const useUserStore = defineStore('user', () => {
    const token = ref(localStorage.getItem('token') || '')
    const user = ref(JSON.parse(localStorage.getItem('user') || '{}'))
    // const router = useRouter()

    function login(userData: any, authToken: string) {
        token.value = authToken
        user.value = userData
        localStorage.setItem('token', authToken)
        localStorage.setItem('user', JSON.stringify(userData))
    }

    function logout() {
        token.value = ''
        user.value = {}
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        window.location.href = '/login'
    }

    return { token, user, login, logout }
})
