const { defineConfig } = require('npm-check-updates')

module.exports = defineConfig({
  target: 'newest',
  cooldown:  7,
})
