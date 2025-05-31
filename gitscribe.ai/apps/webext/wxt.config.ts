import { defineConfig } from 'wxt';
import path from 'path';

// See https://wxt.dev/api/config.html
export default defineConfig({
  modules: ['@wxt-dev/module-react'],
  webExt: {
    startUrls: ['https://meet.google.com'],
    chromiumArgs: ['--user-data-dir=./.wxt/chrome-data'],
  },
  manifest: {
    name: 'GitScribe',
    description: 'A browser extension to create documentation during meetings',
    version: '0.0.1',
    manifest_version: 3,
    permissions: ['activeTab'],
  },
  vite: () => ({
    css: {
      postcss: {
        plugins: [
          require('@tailwindcss/postcss')
        ]
      }
    },
    resolve: {
      alias: {
        '@workspace/ui': path.resolve(__dirname, '../../packages/ui/src'),
        '@workspace/ui/globals.css': path.resolve(__dirname, '../../packages/ui/src/styles/globals.css'),
        '@': path.resolve(__dirname, './src')
      }
    }
  })
});