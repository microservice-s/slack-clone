import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import * as getters from './getters'
// import * as actions from './actions'
// import mutations from './mutations'

Vue.use(Vuex)

const state = {
  currentChannelID: 0,
  messages: {
  }
}

export default new Vuex.Store({
  state,
  getters,
  mutations: {
    FETCH_CHANNELS (state, channels) {
      state.channels = channels
    },
    FETCH_CHANNEL_MESSAGES (state, messages) {
      state.messages[messages[0].channelID] = messages
    },
    SET_USER (state, user) {
      state.user = user
    },
    NEW_MESSAGE (state, message) {
      state.messages[message.channelID].push(message)
    }
  },
  actions: {
    newMessage ({ commit }, { message }) {
      commit('NEW_MESSAGE', message)
    },
    setUser ({ commit }, { user }) {
      commit('SET_USER', user)
      console.log(user)
    },
    fetchChannels ({ commit }, { self }) {
      const auth = localStorage.getItem('authorization')
      console.log(auth)
      var init = {
        method: 'GET',
        baseURL: 'https://api.aethanol.me/v1/',
        url: 'channels',
        headers: {
          Authorization: auth
        }
      }
      axios(init)
        .then((resp) => {
          commit('FETCH_CHANNELS', resp.data)
          self.filterChannels()
        })
        .catch(error => {
          console.log(error.statusText)
        })
    },
    fetchChannelMessages ({ commit }, { self, cID }) {
      const auth = localStorage.getItem('authorization')
      var init = {
        method: 'GET',
        baseURL: 'https://api.aethanol.me/v1/',
        url: 'channels/' + cID,
        headers: {
          Authorization: auth
        }
      }
      axios(init)
        .then((resp) => {
          console.log(cID + ' in the then')
          commit('FETCH_CHANNEL_MESSAGES', resp.data, cID)
          self.response()
        })
        .catch(error => {
          console.log(error.statusText)
        })
    }
  }
})
