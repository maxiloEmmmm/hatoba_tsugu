const path = require('path')
const webpack = require("webpack");
const AddAssetHtmlPlugin = require('add-asset-html-webpack-plugin')
function resolve(dir) {
  return path.join(__dirname, dir)
}
let publicPath = process.env.NODE_ENV === 'production' ?  "/ui" : ""
// let publicPath = ""
module.exports = {
  //子网站使用
  // baseUrl: '/xx/',
  publicPath,
  css: {
    loaderOptions: {
      less: {
        javascriptEnabled: true,
      },
    },
  },
  pages: {
    index: {
      // page 的入口
      entry: 'src/main.js',
      template: 'public/index.html',
      filename: 'index.html',
      title: 'Index Page',
    }
  },
  chainWebpack: config => {
    config.plugin("DllReference")
        .use(
          new webpack.DllReferencePlugin({
            context: process.cwd(),
            manifest: require("./build/vendor/vendor-manifest.json")
          })
        )
      config.plugin("assetsAdd")
        .use(
          new AddAssetHtmlPlugin({
            filepath: path.resolve(__dirname, './build/vendor/*.js'),
            // dll 引用路径
            publicPath: publicPath + '/vendor',
            // dll最终输出的目录
            outputPath: './vendor'
          })
        )
      return config;
  },
  configureWebpack: {
    resolve: {
      extensions: ['.js', '.vue', '.json'],
      alias: {
        'vue$': 'vue/dist/vue.js',
        "@": resolve("./src")
        // "@ant-design/icons/lib/dist$": path.resolve(__dirname, "./build/icons.js")
      }
    }
  }
}