import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/components/Index'
import StrategySubmission from '@/components/StrategySubmission'
import Dashboard from '@/components/Dashboard'
import store from '@/store/index'

Vue.use(Router)

const router = new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'Bastille',
      component: Index,
      beforeEnter: (to, from, next) => {
        if (store.state.authLoggedIn) {
          next('/dashboard')
        } else {
          next()
        }
      }
      // component: Vue.component('dynamic-index', {
      //   functional: true,
      //   render: (createElement, context) => {
      //     return store.state.authLoggedIn
      //       ? createElement(Dashboard, context.data, context.children) : createElement(Index, context.data, context.children)
      //   }
      // })
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard,
      meta: { requiresAuth: true }
    },
    {
      path: '/submit',
      name: 'Strategy Submission',
      component: StrategySubmission,
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (store.state.loading) {
    store.watch(
      (state) => state.loading,
      (value) => {
        if (value === false) {
          next()
        }
      }
    )
  } else {
    next()
  }
})

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!store.state.authLoggedIn) {
      next('/')
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router
