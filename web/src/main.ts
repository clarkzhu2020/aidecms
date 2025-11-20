import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import Home from './views/Home.vue'
import Features from './views/Features.vue'
import Docs from './views/Docs.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/features', component: Features },
  { path: '/docs', component: Docs },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

createApp(App)
  .use(router)
  .mount('#app')
