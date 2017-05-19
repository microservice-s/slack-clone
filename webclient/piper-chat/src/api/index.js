import axios from 'axios'

export function createMessage ({cID, body}, cb) {
  const newMessage = {
    'channelID': cID,
    'body': body
  }

  var init = {method: 'POST',
    baseURL: 'https://localhost:4000/v1/',
    url: 'messages',
    data: newMessage
  }

  return apiRequest(init, cb)
}

export function getAllMessages (cID, cb) {
  var init = {
    method: 'GET',
    baseURL: 'https://localhost:4000/v1/',
    url: 'channel/' + cID
  }

  return apiRequest(init, cb)
}

function apiRequest (init, cb) {
  axios(init)
      .then((resp) => {
        if (cb) cb(false)
        return resp.data
      })
      .catch((error) => {
        if (cb) cb(true)
        return error.response.statusText
      })
}
