import axios from 'axios'

export default {
  apiRequest (init, cb) {
    // var init = {method: method,
    //   baseURL: 'https://api.aethanol.me/v1/',
    //   url: resource,
    //   data: body
    // }
    // make the AJAX response with axios
    return axios(init)
      .then((resp) => {
        if (cb) cb(true)
        return resp.data
      })
      .catch((error) => {
        if (cb) cb(true)
        return error.response.statusText
      })
  }
}
