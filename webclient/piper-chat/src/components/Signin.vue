<template>
    <div class="signin">
        <img src="../assets/logo.png">
        <h1>Piper Chat</h1>
        <h2>Sign In</h2>
        <div v-if="error" class="error">Bad login information</div>
        <form v-on:submit.prevent="signIn" >
            <input v-model="email" type="text" name="email" placeholder="email"><br/>
            <input v-model="password" type="password" name="password" placeholder="password"><br/>
            <input type="submit" value="submit">
        </form>
        <router-link to="/join">Join Piper Chat</router-link>
    </div>
</template>

<script>
    import { signIn } from '../api'
    export default {
      name: 'signin',
      data () {
        return {
          email: '',
          password: '',
          error: false
        }
      },
      methods: {
        signIn: function () {
          signIn(this.email, this.password, this, signedIn => {
            if (!signedIn) {
              this.error = true
            } else {
              this.$router.replace(this.$route.query.redirect || '/chat')
            }
          })
        }
      }
    }
</script>

<style scoped>
</style>
