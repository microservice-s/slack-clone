import axios from 'axios'

var FetchMixin = {
  data () {
    return {
      error: false
    }
  },
  methods: {
    fetchHandler (method, resource, body) {
      console.log(body)
      var init = {method: method,
        baseURL: 'https://api.aethanol.me/v1/',
        url: resource,
        data: body
      }
      // make the AJAX response with axios
      axios(init)
      .then((resp) => {
        return resp
      })
      .catch((error) => {
        this.error = error.response.statusText
      })
    }
//     fetchHandler (method, resource) {
//       var apiUrl = 'https://api.aethanol.me/v1/'
//       console.log('resource: ' + resource)
//       // initialize the method object param
//       var myInit = { method: method }
//       // get the resource url
//       var resourceUrl = apiUrl + resource
//       console.log('url: ' + resourceUrl)
//       // fetch the resource
//       fetch(resourceUrl, myInit)
//         .then(this.status)
//         .then(this.json)
//         // handle data TODO: implement return??
//         .then((data) => {
//           console.log(data)
//         }).catch((err) => {
//           this.$data.error = err
//           console.log(err)
//         })
//     },
//     status: function (response) {
//       if (response.status >= 200 && response.status < 300) {
//         return Promise.resolve(response)
//       } else {
//         return response.text().then(function (errorMessage) {
//           return Promise.reject(new Error(errorMessage))
//         })
//       }
//     },
//     json: function (response) {
//       return response.json()
//     }
  }
}

export default FetchMixin
