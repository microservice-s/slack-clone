<template>
  <div class="channel" v-if="!loading">
    <p>{{ $route.params.id }}</p>
    <message v-for="message in messages[$route.params.id]"
             :message="message"
             :key="message.id">
    </message>
    <button v-on:click.prevent="join">Join Channel</button>
    <form v-on:submit.prevent="sendMessage">
      <input v-model="newMessage" id="newMessage" type="text">
      <input type="submit" value="send">
    </form>
    
  </div>
</template>

<script>
  import { mapGetters } from 'vuex'
  import Message from './Message'
  import { createMessage, joinChannel } from '../api'
  export default {
    data () {
      return {
        newMessage: '',
        cID: this.$route.params.id,
        loading: true
      }
    },
    components: {
      Message
    },
    created () {
      this.$store.dispatch('fetchChannelMessages', { self: this, cID: this.cID })
      // create a new websocket connection
      var auth = localStorage.getItem('authorization')
      auth = auth.substr(7)
      var websock = new WebSocket('wss://api.aethanol.me/v1/websocket?auth=' + auth)
      console.log(websock)
      websock.addEventListener('message', (wsevent) => {
        var event = JSON.parse(wsevent.data)
        switch (event.type) {
          case 'new message':
            window.alert('I got the websockets working, but not with redux... ' + event.data.body)
            this.$store.dispatch('newMessage', {message: event.data})
        }
      })
    },
    methods: {
      sendMessage () {
        createMessage(this.cID, this.newMessage, function (err) {
          console.log(err)
        })
      },
      join () {
        joinChannel(this.cID, (err) => {
          console.log(err)
        })
      },
      response () {
        this.loading = false
      }
    },
    computed: {
      ...mapGetters([
        'messages'
      ])
    }
  }
</script>
