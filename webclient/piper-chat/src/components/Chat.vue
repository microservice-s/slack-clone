<template>
    <div class="chat">
        <side-bar></side-bar>
    </div>
</template>

<script>
    import Sidebar from './Sidebar'
    export default {
      name: 'chat',
      components: {
        'side-bar': Sidebar
      },
      data () {
        return {
          loading: false,
          error: false,
          user: ''
        }
      }
    }
</script>

<style scoped></style>
