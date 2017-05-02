<template>
  <div id="app">
    <img src="./assets/logo.png">
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: 'app'
  // methods: {
  //   fetch: function (resource, method) {
  //     // initialize the method object param
  //     var myInit = { method: method }
  //     // get the resource url
  //     var resourceUrl = this.apiUrl + resource
  //     // fetch the resource
  //     fetch(resourceUrl, myInit)
  //       .then(this.status)
  //       .then(this.json)
  //       // handle data TODO: implement return??
  //       .then(function (data) {
  //         console.log('data: ' + data)
  //       }).catch(function (err) {
  //         console.log(err)
  //       })
  //   },
  //   status: function (response) {
  //     if (response.status >= 200 && response.status < 300) {
  //       return Promise.resolve(response)
  //     } else {
  //       return response.text().then(function (errorMessage) {
  //         console.log(errorMessage)
  //         return Promise.reject(new Error(errorMessage))
  //       })
  //     }
  //   },
  //   json: function (response) {
  //     return response.json()
  //   }
  // }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
