const express = require('express')
const app = express()
const port = process.env.PORT
const index = require('./index')
const moment = require('moment');

app.get('/', async (req, res) => {
  console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + "0" + "," + "0" + "," + "0" + "," + "GET" + "," + "0" + "\n");
  var reply = await index.function_handler(req);
  console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + "0" + "," + "0" + "," + "0" + "," + "GET" + "," + "1" + "\n");
  res.send(reply);
})

app.post('/', async (req, res) => {
  console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "2" + "\n");
  var reply = await index.function_handler(req);
  console.log(moment().format('MMMM Do YYYY, h:mm:sss a') + "," + "0" + "," + "0" + "," + "0" + "," + "POST" + "," + "3" + "\n");
  res.send(reply);
})

app.listen(port, () => {
  console.log('function started...')
})
