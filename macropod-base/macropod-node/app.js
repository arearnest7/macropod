const express = require('express')
const app = express()
const func = require('./func')
const moment = require('moment');
const exec = require('child_process').execSync;

app.post('/', async (req, res) => {
  var reply_t = "";
  var code_t = 500;
  var pv_path_t = "";
  var data = req["Data"];
  var workflow_id = req["WorkflowId"];
  var depth = req["Depth"];
  var width = req["Width"];
  var request_type = req["RequestType"];
  var path = req["PvPath"];
  var ctx = {
    Request: "",
    WorkflowId: workflow_id,
    Depth: depth,
    Width: width,
    RequestType: request_type,
    InvokeType: "GRPC",
    IsJson: false
  };
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "request_start");
  ctx.Request = data;
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "function_start");
  [reply_t, code_t] = await func.FunctionHandler(ctx);
  await console.log(moment(exec("date -u '+%F %H:%M:%S.%6N %Z'").toString(),"YYYY-MM-DD HH:mm:ss.SSSSSS z").format("YYYY-MM-DD HH:mm:ss.SSSSSS UTC") + "," + workflow_id + "," + depth.toString() + "," + width.toString() + "," + request_type + "," + "request_end");
  res.send(reply_t);
})

app.listen(process.env.FUNC_PORT, () => {})
