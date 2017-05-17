const webpack = require('webpack');
const path = require('path');

module.exports = {
  entry: './src/js/app.js',
  output: {
    path: path.resolve(__dirname, 'public/assets/js'),
    publicPath: 'http://localhost:8080/assets/js',
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
    host: "localhost",
    port: 8080,
    contentBase: path.join(__dirname, "public"),
    watchOptions: {
      poll: false,
    }
  }
};
