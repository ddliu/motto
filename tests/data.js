console.log("load data.js");
var helper = require('helper.js');

helper.echo('load helper.js in data.js');
module.exports = [
    {name: "cat", weight: 3},
    {name: "dog", weight: 15},
    {name: "snake", weight: 2},
    {name: "rat", weight: 1},
    {name: "lion", weight: 300},
];