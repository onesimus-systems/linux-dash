/* jshint esversion: 6 */
var express = require('express');
var app     = require('express')();
var server  = require('http').Server(app);
var path    = require('path');
var spawn   = require('child_process').spawn;
var fs      = require('fs');
var ws      = require('websocket').server;

var args = process.argv.slice(2);
var port = 80;

if (args.length > 0) {
    let portStr = args[0];
    if (!isNaN(Number(portStr))) {
        port = Number(portStr);
    }
}

server.listen(port);
console.log(`Linux Dash Server Started On Port ${port}!`);

app.use(express.static(path.resolve(__dirname + '/../')));

app.get('/', function (req, res) {
    res.sendFile(path.resolve(__dirname + '/../index.html'));
});

app.get('/websocket', function (req, res) {
    res.send({
        websocket_support: true,
    });
});

wsServer = new ws({
    httpServer: server
});

var nixJsonAPIScript = __dirname + '/linux_json_api.sh';

function getPluginData(pluginName, callback) {
    var command = spawn(nixJsonAPIScript, [ pluginName, '' ]);
    var output  = [];

    command.stdout.on('data', function(chunk) {
        output.push(chunk.toString());
    });

    command.on('close', function (code) {
        callback(code, output);
    });
}

wsServer.on('request', function(request) {
    var wsClient = request.accept('', request.origin);
    wsClient.on('message', function(wsReq) {
        var moduleName = wsReq.utf8Data;
        getPluginData(moduleName, function(code, output) {
            if (code === 0) {
                var wsResponse = {
                    moduleName: moduleName,
                    output: output.join('')
                };
                wsClient.sendUTF(JSON.stringify(wsResponse));
            }
        });
    });
});

app.get('/server/', function (req, res) {
    getPluginData(req.query.module, function(code, output) {
        if (code === 0) {
            res.send(output.toString());
        } else {
            res.sendStatus(500);
        }
    });
});
