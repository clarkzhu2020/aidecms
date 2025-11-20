import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import MainLayout from '../layout/MainLayout.vue'

const routes: Array<RouteRecordRaw> = [
    {
        path: '/login',
        name: 'Login',
        component: () => import('../views/Login.vue'),
        meta: { requiresAuth: false }
    },
    {
        path: '/',
        component: MainLayout,
        redirect: '/dashboard',
        children: [
            {
                path: 'dashboard',
                name: 'Dashboard',
                component: () => import('../views/Dashboard.vue'),
                meta: { title: 'Dashboard', icon: 'Odometer' }
            },
            {
                path: 'users',
                name: 'Users',
                component: () => import('../views/Users.vue'), // To be created
                meta: { title: 'Users', icon: 'User' }
            },
            {
                path: 'posts',
                name: 'Posts',
                component: () => import('../views/Posts.vue'), // To be created
                meta: { title: 'Posts', icon: 'Document' }
            },
            {
                path: 'web3',
                name: 'Web3',
                meta: { title: 'Web3', icon: 'Wallet' },
                children: [
                    {
                        path: 'wallet',
                        name: 'Web3Wallet',
                        component: () => import('../views/web3/Wallet.vue'),
                        meta: { title: 'Wallet' }
                    },
                    {
                        path: 'transactions',
                        name: 'Web3Transactions',
                        component: () => import('../views/web3/Transactions.vue'),
                        meta: { title: 'Transactions' }
                    }
                ]
            }
        ]
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, _from, next) => {
    const isAuthenticated = localStorage.getItem('token')
    if (to.path !== '/login' && !isAuthenticated) {
        next('/login')
    } else {
        next()
    }
})

export default router
