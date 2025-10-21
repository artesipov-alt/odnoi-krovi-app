import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';
import fs from "fs";
import crypto from "crypto";
import checker from 'vite-plugin-checker';
import svgr from 'vite-plugin-svgr';

const appDirectory = fs.realpathSync(process.cwd());
const resolveApp = (relativePath) => path.resolve(appDirectory, relativePath);

export const alias = {
    pages: resolveApp('src/pages'),
    context: resolveApp('src/context'),
    services: resolveApp('src/services'),
    components: resolveApp('src/components'),
};

export default defineConfig({
    root: './',
    build: {
        outDir: 'build',
    },
    resolve: {
        alias,
    },
    publicDir: 'public',
    define: { global: 'window' },
    plugins: [
        checker({
            typescript: {
                tsconfigPath: './tsconfig.build.json',
            },
            overlay: {
                panelStyle: 'top: 0px; max-height: unset; height: unset;',
            },
        }),
        svgr({ include: '**/*.svg?react' }),
        react(),
    ],
    css: {
        modules: {
            generateScopedName: (name, filename, css) => {
                const hash = crypto.createHash('sha1').update(css).digest('hex').toString().substring(0, 5);
                const fileName = filename.replace(/.*\/(\w+).module.(less|css).*$/, '$1');

                return `${fileName}_${name}_${hash}`;
            },
        },
    },
    server: {
        port: 5173,
        host: '0.0.0.0',
        proxy: {
            '/api': {
                target: 'http://localhost:3000', // Локальный бэкенд
                changeOrigin: true,
                secure: false
            }
        }
    }
});
