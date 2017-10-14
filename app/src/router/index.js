import Vue from 'vue'
import Router from 'vue-router'
import Index from '@/components/Index'
import StrategySubmission from '@/components/StrategySubmission'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Bastille',
      component: Index
    },
    {
      path: '/leaderboard',
      name: 'Leaderboard',
      component: null
    },
    {
      path: '/submit',
      name: 'Strategy Submission',
      component: StrategySubmission
    }
  ]
})
