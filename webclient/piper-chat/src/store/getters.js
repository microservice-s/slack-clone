export const channels = state => state.channels

export const messages = state => state.messages

export const currentChannel = state => {
  return state.currentChannelID
        ? state.channels[state.currentChannelID]
        : {}
}

