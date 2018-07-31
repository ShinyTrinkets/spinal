const { watch } = require('chokidar')

module.exports = function trigger (filesOrFolders, action) {
  const trigger = watch(filesOrFolders)
  trigger.on('change', action)
  return trigger
}
