const path = require("path");

const root = __dirname;
const feDir = path.join(root, "fe");
const beDir = path.join(root, "be");

/** @type {import("pm2").StartOptions[]} */
const apps = [
  {
    name: "hani-be",
    cwd: beDir,
    script: path.join(beDir, "bin/api"),
    interpreter: "none",
    instances: 1,
    exec_mode: "fork",
    autorestart: true,
    max_memory_restart: "1G",
    env: {
      NODE_ENV: "production",
    },
    env_production: {
      NODE_ENV: "production",
    },
  },
  {
    name: "hani-fe",
    cwd: feDir,
    script: path.join(feDir, "node_modules/next/dist/bin/next"),
    args: "start",
    interpreter: "node",
    instances: 1,
    exec_mode: "fork",
    autorestart: true,
    max_memory_restart: "512M",
    env: {
      NODE_ENV: "production",
      PORT: 3005,
    },
    env_production: {
      NODE_ENV: "production",
      PORT: 3005,
    },
  },
];

module.exports = { apps };
