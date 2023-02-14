const HtmlWebpackPlugin = require("html-webpack-plugin");
const { resolve } = require("path");

module.exports = {
  mode: "development",
  devtool: "inline-source-map",

  entry: "./src/main.tsx",
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: "ts-loader",
        exclude: /node_modules/,
      },
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      },
      {
        test: /\.svg$/,
        loader: "svg-inline-loader",
      },
    ],
  },
  plugins: [new HtmlWebpackPlugin()],

  resolve: {
    extensions: [".tsx", ".ts", ".js"],
    alias: {
      "@": resolve(__dirname, "src"),
    },
    fallback: { fs: false },
  },

  output: {
    filename: "bundle.js",
    path: resolve(__dirname, "dist"),
  },
};
