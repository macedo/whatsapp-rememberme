const esbuild = require("esbuild");
const { stimulusPlugin } = require("esbuild-plugin-stimulus");

esbuild.build({
  entryPoints: [
    "./web/src/application.js"
  ],
  bundle: true,
  outfile: "./web/static/application.bundle.js",
  plugins: [stimulusPlugin()],
}).catch(() => process.exit(1));
