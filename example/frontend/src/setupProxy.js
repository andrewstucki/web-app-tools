const proxy = require("http-proxy-middleware");

module.exports = app => {
  const options = {
    target: "http://127.0.0.1:3456"
  };
  app.use(proxy("/oauth", options));
  app.use(proxy("/api", options));
};
