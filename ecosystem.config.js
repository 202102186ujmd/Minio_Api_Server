module.exports = {
  apps: [
    {
      name: "minio-api",
      script: "bin/minio-api",
      exec_mode: "fork",
      instances: 1,
      watch: false,
      autorestart: true,
      max_memory_restart: "512M",
      env: {
        APP_ENV: "production",
        APP_PORT: "8080",
        LOG_LEVEL: "info",
        LOG_FORMAT: "json"
      }
    }
  ]
};
