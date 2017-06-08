<template>
    <div class="profile">
        <div v-if="error">Error loading user</div>
        <img :src="user.photoURL" alt="profile picture">
        <div>{{user.userName}}</div>
        <input v-model="user.firstName"></input> <br>
        <input v-model="user.lastName"></input> <br>
        <input v-on:click.prevent="save" type="submit" value="save">
        <div>{{user.email}}</div>
        <router-link to="/signout">Sign Out</router-link>
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
          user: ''
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
          this.user = data
        })
      },
      methods: {
        save: function () {
          var init = {method: 'PATCH',
            baseURL: 'https://api.aethanol.me/v1/',
            url: 'users/me',
            headers: {'Authorization': localStorage.authorization},
            data: {
              firstName: this.user.firstName,
              lastName: this.user.lastName
            }
          }
          ajax.apiRequest(init, err => {
            if (err) {
              this.error = true
            }
          }).then(data => {
          })
        }
      }
    }
</script>

<style scoped></style>
