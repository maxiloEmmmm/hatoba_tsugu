const path = require("path");
const webpack = require("webpack");
// const {BundleAnalyzerPlugin} = require('webpack-bundle-analyzer')

module.exports = {
  entry: {
    vendor: [
      "vue-router",
      "vuex",
      "axios",
      "vue/dist/vue.js",
      "vee-validate",
      "moment",
      "maxilo-vue",
      "lodash",
      "vue-i18n",
      "codemirror",
      "js-beautify",
      "core-js",
    ]
  },
  output: {
    path: path.join(__dirname, "./vendor"),
    filename: "[name].dll.js",
    library: "[name]_[hash]"
  },
  plugins: [
    new webpack.DllPlugin({
      path: path.join(__dirname, "./vendor", "[name]-manifest.json"),
      name: "[name]_[hash]",
      context: process.cwd()
    }),
    new webpack.ContextReplacementPlugin(/moment[/\\]locale$/, /zh-cn$/),
    // new BundleAnalyzerPlugin()
  ]
};