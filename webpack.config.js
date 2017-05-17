const webpack = require('webpack');
const path = require('path');
const url = require('url');

require('dotenv').config();

const hostStr = (process.env.WEBPACK_DEV_SERVER_HOST !== undefined) ?
  process.env.WEBPACK_DEV_SERVER_HOST : 'http://localhost:8080';
const port = (hostStr !== undefined && url.parse(hostStr).port) ?
  url.parse(hostStr).port : '8080';
const hostname = (hostStr !== undefined && url.parse(hostStr).hostname) ?
  url.parse(hostStr).hostname : '8080';

module.exports = {
  entry: './src/js/app.js',
  output: {
    path: path.resolve(__dirname, 'public/assets/js'),
    publicPath: 'http://' + path.join(`${hostname}:${port}`, 'assets/js'),
    filename: 'common.js',
  },
  module: {
    rules: [
      {
        'test': /\.(js|jsx)$/,
        include: path.join(__dirname, 'src', 'js'),
        'use': [
          'react-hot-loader',
          'babel-loader',
        ],
      },
      {
        test: /\.scss$/,
        use: [{
            loader: "style-loader" // creates style nodes from JS strings
        }, {
            loader: "css-loader" // translates CSS into CommonJS
        }, {
            loader: "sass-loader" // compiles Sass to CSS
        }]
      },
    ],
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
  ],
  devServer: {
    hot: true, // this enables hot reload
    //hotOnly: true, // do not reload browser if hot reload failed
    inline: true, // use inline method for hmr
    headers: {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, PATCH, OPTIONS",
      "Access-Control-Allow-Headers": "X-Requested-With, content-type, Authorization"
    },
    host: hostname,
    port: port,
    contentBase: path.join(__dirname, "public"),
    watchOptions: {
      poll: false,
    }
  }
};
