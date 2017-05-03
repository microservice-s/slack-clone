<template>
    <div class="profile">
        <router-link to="/signout">Sign Out</router-link>
        <div>{{user.userName}}</div>
        <div>{{user.firstName}}</div>
        <div>{{user.lastName}}</div>
        <div>{{user.email}}</div>
    </div>
</template>

<script>
    import ajax from '../services/ajax'
    export default {
      name: 'profile',
      data () {
        return {
          loading: false,
          error: false,
          user: null
        }
      },
      created: function () {
        var init = {method: 'GET',
          baseURL: 'https://api.aethanol.me/v1/',
          url: 'users/me',
          headers: {'Authorization': localStorage.authorization}
        }
        ajax.apiRequest(init, err => {
          if (err) {
            this.error = true
          }
        }).then(data => {
          console.log(data)
          this.user = data
        })
      }
    }
</script>

<style scoped></style>
