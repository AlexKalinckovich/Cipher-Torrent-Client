import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import { fileURLToPath, URL } from 'node:url';

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, '../containerization', '');
    const frontendPort = Number(env.FRONTEND_PORT ?? '5173');

    return {
        plugins: [react()],
        resolve: {
            alias: {
                '@': fileURLToPath(new URL('./src', import.meta.url))
            }
        },
        server: {
            port: frontendPort,
            strictPort: true,
        },
    };
});