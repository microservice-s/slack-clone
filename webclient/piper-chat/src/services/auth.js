import axios from 'axios'

export default {
  signIn (email, pass, cb) {
    var user = {
      email: email,
      password: pass
    }
    var init = {method: 'POST',
      baseURL: 'https://api.aethanol.me/v1/',
      url: 'sessions',
      data: user
    }
    // make the AJAX response with axios
    axios(init)
      .then((resp) => {
        // store the auth header
        localStorage.authorization = resp.headers.authorization
        if (cb) cb(true)
        this.onChange(true)
      })
      .catch((error) => {
        console.log(error)
        if (cb) cb(false) // make the callback fail
        this.onChange(false)
      })
  },
  join (user, cb) {
    var init = {method: 'POST',
      baseURL: 'https://api.aethanol.me/v1/',
      url: 'users',
      data: user
    }
    // make the AJAX response with axios
    return axios(init)
      .then((resp) => {
        localStorage.authorization = resp.headers.authorization
        if (cb) cb(true)
        this.onChange(true)
      })
      .catch((error) => {
        if (cb) cb(false) // make the callback fail
        this.onChange(false)
        return error.response.data
      })
  },
  getToken () {
    return localStorage.authorization
  },
  signOut () {
    var init = {method: 'DELETE',
      baseURL: 'https://api.aethanol.me/v1/',
      url: 'sessions/mine',
      headers: {'Authorization': localStorage.authorization}
    }
    // make the AJAX response with axios
    axios(init)
      .then((resp) => {
        this.onChange(true)
      })
      .catch((error) => {
        console.log(error)
        this.onChange(false)
      })
    this.delToken()
    // if (cb) cb()
    // this.onChange(false)
  },
  delToken () {
    delete localStorage.authorization
  },
  signedIn () {
    if (!localStorage.authorization) {
      return false
    }
    var init = {method: 'GET',
      baseURL: 'https://api.aethanol.me/v1/',
      url: 'users/me',
      headers: {'Authorization': localStorage.authorization}
    }
    // make the AJAX response with axios
    return axios(init)
      .then((resp) => {
        // store the auth header
        return true
        // this.onChange(true)
      })
      .catch((error) => {
        console.log(error)
        this.delToken()
        return false
        // this.onChange(false)
      })
  },

  onChange () {}
}
