import axios from 'axios'

const api = 'https://api.aethanol.me/v1/'

// reused function to do all the api interaction
function apiRequest (init, cb) {
  return axios(init)
      .then((resp) => {
        if (cb) cb(false)
        return resp.data
      })
      .catch((error) => {
        if (cb) cb(true)
        console.log(error)
        return error.response.statusText
      })
}

export function createMessage (cID, body, cb) {
  const newMessage = {
    'channelID': cID,
    'body': body
  }
  const auth = localStorage.getItem('authorization')
  console.log(auth)
  var init = {method: 'POST',
    baseURL: api,
    url: 'messages',
    data: newMessage,
    headers: {
      Authorization: auth
    }
  }

  return apiRequest(init, cb)
}

export function signIn (email, pass, self, cb) {
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
        // console.log(resp)
        self.$store.dispatch('setUser', { user: resp.data })
        if (cb) cb(true)
        this.onChange(true)
      })
      .catch((error) => {
        console.log(error)
        if (cb) cb(false) // make the callback fail
        this.onChange(false)
      })
}

export function getAllChannels (cb) {
  const auth = localStorage.getItem('authorization')
  console.log(auth)
  var init = {
    method: 'GET',
    baseURL: api,
    url: 'channels',
    headers: {
      Authorization: auth
    }
  }
  return apiRequest(init, cb)
}

export function joinChannel (cID, cb) {
  const auth = localStorage.getItem('authorization')
  var init = {method: 'LINK',
    baseURL: api,
    url: 'channels/' + cID,
    headers: {
      Authorization: auth
    }
  }
  console.log(init)

  return apiRequest(init, cb)
}

export function getAllMessages (cID, cb) {
  var init = {
    method: 'GET',
    baseURL: api,
    url: 'channels/' + cID
  }

  return apiRequest(init, cb)
}
