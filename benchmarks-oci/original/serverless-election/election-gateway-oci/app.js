const express = require('express')
const app = express()
const port = process.env.PORT
const index = require('./index')

app.get('/', (req, res) => {
  res.send(index.function_handler(req))
})

app.listen(port, () => {
  console.log('function started...')
})
