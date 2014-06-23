console.log("load index.js");
var data = require('./data.js');
var sort = require('./helper.js').sort;

return sort(data)[0].name;