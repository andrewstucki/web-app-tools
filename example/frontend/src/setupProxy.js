const path = require('path').resolve(process.cwd(), '..', '.env');
require('dotenv').config({
  path,
});
const proxy = require('http-proxy-middleware');

module.exports = (app) => {
  const target = process.env.HOST_PORT || '127.0.0.1:3456';
  const options = {
    target: 'http://' + target,
  };
  app.use(proxy('/oauth', options));
  app.use(proxy('/api', options));
};
