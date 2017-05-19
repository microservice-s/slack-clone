/*
 * action types
 */
// channels
export const NEW_CHANNEL = 'NEW_CHANNEL'
export const UPDATE_CHANNEL = 'UPDATE_CHANNEL'
export const DELETE_CHANNEL = 'DELETE_CHANNEL'

export const USER_JOINED_CHANNEL = 'USER_JOINED_CHANNEL'
export const USER_LEFT_CHANNEL = 'USER_LEFT_CHANNEL'

// users
export const NEW_USER = 'NEW_USER'

// messages
export const NEW_MESSAGE = 'NEW_MESSAGE'
export const UPDATE_MESSAGE = 'UPDATE_MESSAGE'
export const DELETE_MESSAGE = 'DELETE_MESSAGE'


export function newChannel (channel) {
  return {
    type: NEW_CHANNEL,
    channel
  }
}

export function updateChannel (channel) {
  return {
    type: UPDATE_CHANNEL,
    channel
  }
}

export function deleteChannel (cID) {
  return {
    type: DELETE_CHANNEL,
    cID
  }
}

export function userJoined (cID, user) {
  return {
    type: USER_JOINED_CHANNEL,
    cID,
    user
  }
}

export function userLeft (cID, user) {
  return {
    type: USER_LEFT_CHANNEL,
    cID,
    user
  }
}


