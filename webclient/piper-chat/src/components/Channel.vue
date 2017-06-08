<template>
  <div class="channel">
    <message v-for="message in messages"
             :message="message">
    </message>
  </div>
</template>

<script>
  import Message from './Message'
  export default {
    data () {
      return {
        messages: [
          {content: 'Hey! Ho!'},
          {content: 'Let\'s go!'}
        ]
      }
    },
    components: {
      Message
    }
  }
</script>
