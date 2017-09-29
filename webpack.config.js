var webpack = require('webpack');

module.exports = {
  entry: './main.js',
  plugins: [
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': '"production"'
    }),
  ],
  output: {
    filename: 'dist/bundle.js',
    libraryTarget: 'var',
    library: 'mw',
  }

}
