import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

import auth from '@/services/auth'
import Signin from '@/components/Signin'
import Join from '@/components/Join'
import Test from '@/components/Test'
import Profile from '@/components/Profile'
import Chat from '@/components/Chat'

function requireAuth (to, from, next) {
  if (!auth.loggedIn()) {
    next({
      path: '/'// ,
      // query: { redirect: to.fullPath }
    })
  } else {
    next()
  }
}

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Signin',
      component: Signin
    },
    {
      path: '/join',
      name: 'Join',
      component: Join
    },
    {
      path: '/test',
      name: 'Test',
      component: Test
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile,
      beforeEnter: requireAuth
    },
    {
      path: '/chat',
      name: 'Chat',
      component: Chat,
      beforeEnter: requireAuth
    }
  ]
})
