const express = require('express')
const app = express()
const port = process.env.PORT
const index = require('./index')
const redis = require('redis');
const moment = require('moment');

var client = redis.createClient();

if ("LOGGING_NAME" in process.env) {
  const client = redis.createClient({url: 'redis://' + process.env.LOGGING_IP, password: process.env.LOGGING_PASSWORD});
}

app.get('/', async (req, res) => {
  if ("LOGGING_NAME" in process.env) {
    await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + "0" + "," + "0" + "," + "0" + "," + "http" + "," + "0" + "\n");
  }
  var reply = await index.function_handler(req);
  if ("LOGGING_NAME" in process.env) {
    await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + "0" + "," + "0" + "," + "0" + "," + "http" + "," + "1" + "\n");
  }
  res.send(reply);
})

app.post('/', async (req, res) => {
  if ("LOGGING_NAME" in process.env) {
    await client.connect();
    await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + "0" + "," + "0" + "," + "0" + "," + "http" + "," + "2" + "\n");
  }
  var reply = await index.function_handler(req);
  if ("LOGGING_NAME" in process.env) {
    await client.append(process.env.LOGGING_NAME, moment().format('MMMM Do YYYY, h:mm:ss a') + "," + "0" + "," + "0" + "," + "0" + "," + "http" + "," + "3" + "\n");
  }
  res.send(reply);
})

app.listen(port, () => {
  console.log('function started...')
})
