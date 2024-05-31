import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  assetsInclude: [
      'src/assets/topolith.wasm',
      'src/assets/wasm_exec.js',
  ],
})
