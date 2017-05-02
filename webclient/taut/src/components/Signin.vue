<template>
    <div class="signin">
        <h1>Sign In</h1>
        <!--<div class="loading" v-if="loading">
            Loading...
        </div>-->
        <div v-if="error" class="error">Bad login information</div>
        <form v-on:submit.prevent="signIn" >
            <input v-model="user.userName" type="text" name="username" placeholder="username"><br/>
            <input v-model="user.password" type="password" name="password" placeholder="password"><br/>
            <input type="submit" value="submit">
        </form>
        <router-link to="/join">Join Taut</router-link>
    </div>
</template>

<script>
    import FetchMixin from './mixins'
    export default {
      name: 'signin',
      data () {
        return {
          user: {
            userName: '',
            password: ''
          }
        }
      },
      mixins: [FetchMixin],
      methods: {
        signIn: function () {
          var user = {
            userName: this.user.userName,
            password: this.user.password
          }
          this.fetchHandler('POST', 'sessions', user)
        }
      }
    }
</script>

<style scoped>
</style>
