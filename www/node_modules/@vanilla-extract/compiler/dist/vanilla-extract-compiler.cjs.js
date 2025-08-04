'use strict';

if (process.env.NODE_ENV === "production") {
  module.exports = require("./vanilla-extract-compiler.cjs.prod.js");
} else {
  module.exports = require("./vanilla-extract-compiler.cjs.dev.js");
}
