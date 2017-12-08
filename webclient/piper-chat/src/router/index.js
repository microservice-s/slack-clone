import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

import auth from '@/services/auth'
import Signin from '@/components/Signin'
import Join from '@/components/Join'
import Test from '@/components/Test'
import Profile from '@/components/Profile'
import Chat from '@/components/Chat'
import NotFound from '@/components/NotFound'
import Channel from '@/components/Channel'

function requireAuth (to, from, next) {
  if (!auth.signedIn()) {
    next({
      path: '/'// ,
      // query: { redirect: to.fullPath }
    })
  } else {
    next()
  }
}

function authed (to, from, next) {
  if (auth.signedIn()) {
    next({
      path: '/chat'// ,
      // query: { redirect: to.fullPath }
    })
  } else {
    next()
  }
}

export default new Router({
  mode: 'history',
  routes: [
    { path: '*', component: NotFound },
    {
      path: '/',
      name: 'Signin',
      component: Signin,
      beforeEnter: authed
    },
    {
      path: '/join',
      name: 'Join',
      component: Join,
      beforeEnter: authed
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
    },
    {
      path: '/channel/:id',
      name: 'Channel',
      component: Channel,
      beforeEnter: requireAuth
    },
    { path: '/signout',
      beforeEnter (to, from, next) {
        auth.signOut()
        next('/')
      }
    }
  ]
})
