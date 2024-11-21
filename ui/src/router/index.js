import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/upgrade',
      name: 'upgrade',
      component: () => import('@/components/Upgrade')
    }
  ]
})
