import fs from "node:fs";
import path from "node:path";

import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

function loadRootFrontendEnv() {
  const rootEnvPath = path.resolve(__dirname, "..", ".env.frontend");
  if (!fs.existsSync(rootEnvPath)) {
    return;
  }

  const content = fs.readFileSync(rootEnvPath, "utf8");
  for (const rawLine of content.split(/\r?\n/)) {
    const line = rawLine.trim();
    if (!line || line.startsWith("#")) {
      continue;
    }

    const index = line.indexOf("=");
    if (index <= 0) {
      continue;
    }

    const key = line.slice(0, index).trim();
    let value = line.slice(index + 1).trim();

    if (
      (value.startsWith('"') && value.endsWith('"')) ||
      (value.startsWith("'") && value.endsWith("'"))
    ) {
      value = value.slice(1, -1);
    }

    process.env[key] = value;
  }
}

loadRootFrontendEnv();

const proxyTarget = process.env.VITE_API_BASE_URL || "http://localhost:8080";
const proxyTeamTarget = process.env.VITE_TEAM_API_BASE_URL || "http://localhost:8081";
const proxyAssetTarget = process.env.VITE_ASSET_API_BASE_URL || "http://localhost:8082";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      "/api/v1/auth": {
        target: proxyTarget,
        changeOrigin: true,
      },
      "/api/v1/users": {
        target: proxyTarget,
        changeOrigin: true,
      },
      "/api/v1/teams": {
        target: proxyTeamTarget,
        changeOrigin: true,
      },
      "/api/v1/folders": {
        target: proxyAssetTarget,
        changeOrigin: true,
      },
      "/api/v1/notes": {
        target: proxyAssetTarget,
        changeOrigin: true,
      },
      "/api/v1/shares": {
        target: proxyAssetTarget,
        changeOrigin: true,
      },
    },
  },
});
