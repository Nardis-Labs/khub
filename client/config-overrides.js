/*eslint-disable */
const { override, addWebpackModuleRule } = require('customize-cra');

module.exports = {
  webpack: override(
    addWebpackModuleRule({
      test: [/\.css$/, /\.less$/],
      use: ['style-loader', 'css-loader', 'postcss-loader', { loader: 'less-loader' }]
    })
  )
};
