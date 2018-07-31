const cron = require('node-cron')

module.exports = function trigger (expression, action) {
  return cron.schedule(expression, action.bind(this, 'timer'))
}
