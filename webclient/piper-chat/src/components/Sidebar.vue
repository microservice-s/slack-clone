<template>
    <div class="side-bar" v-if="!loading">
      <h2>Channels</h2>
      <ul>
        <li v-for="channel in channels"
          v-bind:key="channel.id" >
          <router-link :to="/channel/ + channel.id">{{channel.name}}</router-link>
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
          loading: true,
          error: false,
          user: ''
        }
      },
      created () {
        this.$store.dispatch('fetchChannels', { self: this })
      },
      computed: {
        ...mapGetters([
          'channels'
        ])
      },
      methods: {
        filterChannels () {
          this.loading = false
          // return this.$store.state.channels
        }
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
