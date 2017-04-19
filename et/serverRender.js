var webPage = require('webpage');
var page = webPage.create();
var system = require('system')
page.open(system.args[1], function (status) {
  var content = page.content;
  console.log(content);
  phantom.exit();
});
