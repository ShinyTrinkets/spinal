const timer = require('./timer')
const watcher = require('./watcher')

module.exports = function trigger (name, props, action) {
  if (name === 'timer') {
    return timer(props, action)
  }
  if (name === 'watcher') {
    return watcher(props, action)
  }
  throw new Error(`Unknown trigger type: ${name}`)
}
