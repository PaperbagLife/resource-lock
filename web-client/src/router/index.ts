import { createRouter, createWebHistory } from 'vue-router'
import GlobalView from '@/components/GlobalView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: GlobalView
    },
  ]
})

export default router
