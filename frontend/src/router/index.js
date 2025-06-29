import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Me from '../views/Me.vue'
import Admin from '../views/Admin.vue'
import Callback from '../views/Callback.vue'
import Deposit from '../views/Deposit.vue'
import Withdraw from '../views/Withdraw.vue'

const routes = [
  { path: '/', component: Home, meta: { requiresAuth: true } },
  { path: '/login', component: Login },
  { path: '/me', component: Me, meta: { requiresAuth: true } },
  { path: '/admin', component: Admin, meta: { requiresAuth: true, admin: true } },
  { path: '/deposit', component: Deposit, meta: { requiresAuth: true } },
  { path: '/withdraw', component: Withdraw, meta: { requiresAuth: true } },
  { path: '/callback', component: Callback },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const roles = JSON.parse(localStorage.getItem('roles') || '[]')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.meta.admin && !roles.includes('ROLE_ADMIN')) {
    next('/')
  } else if ((to.path === '/login' || to.path === '/callback') && token) {
    next('/')
  } else {
    next()
  }
})

export default router