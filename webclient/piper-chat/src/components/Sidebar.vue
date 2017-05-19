<template>
    <div class="side-bar">
      <h2>Channels</h2>
      <!--<router-link v-for"channel in channels">
        {{channel.title}}
      </router-link>-->
      <ul>
        <li v-for="channel in channels">
          {{channel.title}}
        </li>
      </ul>
      <router-link to="/signout">Sign Out</router-link>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
      name: 'side-bar',
      data () {
        return {
          loading: false,
          error: false,
          user: ''
        }
      },
      computed: {
        ...mapGetters([
          'channels'
        ])
      }
    }
</script>

<style scoped>
  .side-bar {
    height: 100%; /* 100% Full-height */
    position: fixed; /* Stay in place */
    z-index: 1; /* Stay on top */
    top: 0;
    left: 0;
    overflow-x: hidden; /* Disable horizontal scroll */
    padding-top: 20px; /* Place content 60px from the top */
    transition: 0.5s; /* 0.5 second transition effect to slide in the sidenav */
  }
</style>
