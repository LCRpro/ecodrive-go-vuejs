import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import Me from '../views/Me.vue'
import Admin from '../views/Admin.vue'
import Callback from '../views/Callback.vue'
import Deposit from '../views/Deposit.vue'
import Withdraw from '../views/Withdraw.vue'
import Course from '../views/Course.vue'
import Driver from '../views/Driver.vue'
import Mentions from '../views/Mentions.vue'
import Conditions from '../views/Conditions.vue'

const routes = [
  { path: '/', component: Home },
  { path: '/login', component: Login },
  { path: '/me', component: Me, meta: { requiresAuth: true } },
  { path: '/admin', component: Admin, meta: { requiresAuth: true, admin: true } },
  { path: '/deposit', component: Deposit, meta: { requiresAuth: true } },
  { path: '/withdraw', component: Withdraw, meta: { requiresAuth: true } },
  { path: '/course', component: Course, meta: { requiresAuth: true } },
  { path: '/driver', component: Driver, meta: { requiresAuth: true } },
  { path: '/mentions', component: Mentions },
  { path: '/conditions', component: Conditions },
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
  } else if (to.meta.passager && !roles.includes('ROLE_PASSAGER')) {
    next('/')
  } else if (to.meta.driver && !roles.includes('ROLE_DRIVER')) {
    next('/')
  } else if ((to.path === '/login' || to.path === '/callback') && token) {
    next('/')
  } else {
    next()
  }
})

export default router