<template>
    <div class="join">
         <img src="../assets/logo.png">
        <h1>Join Twat</h1>
         <div v-if="error" class="error">{{errorMessage}}</div>
        <form v-on:submit.prevent="join">
            <input v-model="user.email" type="text" name="email" placeholder="email"><br/>
            <input v-model="user.userName" type="text" name="userName" placeholder="username"><br/>
            <input v-model="user.firstName" type="text" name="firstName" placeholder="first name"><br/>
            <input v-model="user.lastName" type="text" name="lastName" placeholder="last name"><br/>
            <input v-model="user.password" type="password" name="password" placeholder="Enter your password"><br/>
            <input v-model="user.passwordConf" type="password" name="passwordConf" placeholder="Repeat your password"><br/>
            <input type="submit" value="submit">
        </form>
    </div>
</template>

<script>
    import auth from '../services/auth'
    export default {
      name: 'join',
      data () {
        return {
          user: {
            email: '',
            userName: '',
            firstName: '',
            lastName: '',
            password: '',
            passwordConf: ''
          },
          error: false,
          errorMessage: ''
        }
      },
      methods: {
        join: function () {
          var user = {
            email: this.user.email,
            userName: this.user.userName,
            firstName: this.user.firstName,
            lastName: this.user.lastName,
            password: this.user.password,
            passwordConf: this.user.passwordConf
          }
          auth.join(user, joined => {
            if (!joined) {
              this.error = true
            } else {
              this.$router.replace(this.$route.query.redirect || '/profile')
            }
          }).then(data => {
            this.errorMessage = data
          })
        }
      }
    }
</script>

<style scoped>
</style>
