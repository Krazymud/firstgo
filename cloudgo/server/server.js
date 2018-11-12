'use strict';
let du_user = 0;
let du_num = 0;
let du_phone = 0;
let du_email = 0;
let http = require("http");
let fs = require('fs');
let querystring = require('querystring');
let store = {};
let read = function() {
    fs.readFile('json/store.json', function (err, data) {
        if (err) {
            console.log(err);
        }
        else {
            if(data.toString() !== '')
                store = JSON.parse(data);
        }
    })
}
read();
function checkDuplicate(data, json){
    if(JSON.stringify(data) === "{}")
        return 1;
    du_user = du_num = du_phone = du_email = 0;
    for(let key in data){
        if(key === json.username)
            du_user = 1;
        if(data[key].number === json.number)
            du_num = 1;
        if(data[key].phone === json.phone)
            du_phone = 1;
        if(data[key].email === json.email)
            du_email = 1;
    }
    return (du_user === 0 &&  du_num === 0 &&
        du_phone === 0 && du_email === 0);
}
let server = http.createServer(function(request, response) {
    console.log(request.method + ': ' + request.url);
    let body = "";
    request.on('data', function (chunk) {
        body += chunk;  //一定要使用+=，如果body=chunk，因为请求favicon.ico，body会等于{}
        let json = querystring.parse(body);
        if(checkDuplicate(store, json)) {
            store[json.username] = {};
            store[json.username].number = json.number;
            store[json.username].phone = json.phone;
            store[json.username].email = json.email;
            console.log(JSON.stringify(json));
            fs.writeFile('json/store.json', JSON.stringify(store), function (err) {
                if (err) {
                    console.log(err);
                }
                else {
                    console.log('ok');
                }
            })
            response.writeHead(200, {"Content-Type": "text/html"});
            response.end();
        }
        else{
            response.writeHead(409, {"Content-Type": "text/html"})
            if(du_user === 1){
                response.write("用户名重复！");
            }
            if(du_num === 1){
                response.write("学号重复！");
            }
            if(du_phone === 1){
                response.write("电话号码重复！");
            }
            if(du_email === 1){
                response.write("邮箱重复！");
            }
            response.end();
        }
    });
    let path;
    if(request.url.indexOf('?') === -1) {
        if (request.url === '/')
            path = process.cwd() + '/html/index.html';
        else
            path = process.cwd() + request.url;
        console.log(path);
        fs.readFile(path, function (err, data) {
            if (err) {
                console.log(err);
            }
            else {
                let type = path.substr(path.lastIndexOf('.') + 1);
                console.log(type);
                if(type === 'jpg')
                	response.writeHead(200, {"Content-Type": "image/jpeg"});
                else if(type === 'ico')
                	response.writeHead(200, {"Content-Type": "image/icon"});
                else
                	response.writeHead(200, {"Content-Type": "text/" + type});
                if (type === "text")
                	response.write(data.toString());
                else response.write(data, "binary");
                response.end();
            }
        })
    }
    else{
        if(isinstore(request.url)) {
            let key = request.url.substr(request.url.indexOf('=') + 1);
            response.writeHead(200, {"Content-Type": "text/html"});
            response.write("<!DOCTYPE html>" +
                "<html> " +
                "<head><meta charset=\"utf-8\">" +
                "<script type=\"text/javascript\" src=\"node_modules/jquery/dist/jquery.min.js\"></script>" +
                "<link href='css/user.css' type='text/css' rel='stylesheet'/>" +
                "<link rel = 'Shortcut Icon' href='../favicon.ico>'" +
                "</head> " +
                "<body> "  +
                "<div id=\"frame\"> " +
                "<h2>用户详情</h2> " +
                "<div id=\"content\">" +
                "<p>用户名：" + key + "</p>" + "<p>学号：" + store[key].number + "</p>" +
                "<p>电话：" + store[key].phone + "</p>" + "<p>邮箱：" + store[key].email + "</p>" +
                "</div> " +
                "<button id='return'>返回</button>" +
                "</div> " +
                "<script src=\"../js/content.js\"></script>" +
                "</body> " +
                "</html>");
            response.end();
        }
    }
});

function isinstore(str){
    for(let key in store){
        if(key === str.substr(str.indexOf('=') + 1))
            return 1;
    }
}

server.listen(8000);
console.log("Server is listening");