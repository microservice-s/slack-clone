import axios from 'axios'

export default {
  apiRequest (init, cb) {
    // make the AJAX response with axios
    return axios(init)
      .then((resp) => {
        if (cb) cb(false)
        return resp.data
      })
      .catch((error) => {
        if (cb) cb(true)
        return error.response.statusText
      })
  }
}
