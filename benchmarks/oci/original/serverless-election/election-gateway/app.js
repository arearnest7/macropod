const express = require('express')
const app = express()
const port = process.env.PORT
const index = require('./index')
const moment = require('moment');
const exec = require('child_process').execSync;

app.use(express.json());

app.get('/', async (req, res) => {
  var workflow_id = Math.floor(Math.random() * 10000000).toString();
  var workflow_depth = 0;
  var workflow_width = 0;
  var newbody = req.body;
  if ("workflow_id" in req.body) {
    workflow_id = req.body["workflow_id"];
    workflow_depth = req.body["workflow_depth"];
    workflow_width = req.body["workflow_width"];
  } else {
    newbody["workflow_id"] = workflow_id;
    newbody["workflow_depth"] = workflow_depth;
    newbody["workflow_width"] = workflow_width;
  }
  console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "0");
  console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "1");
  var reply = await index.function_handler(newbody);
  console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "2");
  res.send(reply);
})

app.post('/', async (req, res) => {
  var workflow_id = Math.floor(Math.random() * 10000000).toString();
  var workflow_depth = 0;
  var workflow_width = 0;
  var newbody = req.body;
  if ("workflow_id" in req.body) {
    workflow_id = req.body["workflow_id"];
    workflow_depth = req.body["workflow_depth"] + 1;
    workflow_width = req.body["workflow_width"];
  } else {
    newbody["workflow_id"] = workflow_id;
    newbody["workflow_depth"] = workflow_depth;
    newbody["workflow_width"] = workflow_width;
  }
  console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "0");
  console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "1");
  var reply = await index.function_handler(newbody);
  console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + workflow_depth.toString() + "," + workflow_width.toString() + "," + "HTTP" + "," + "2");
  res.send(reply);
})

app.listen(port, () => {
  console.log('function started...')
})
