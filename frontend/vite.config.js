import process from 'node:process'
import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(() => {
    const backendURL = process.env.BACKEND_URL || ''

    return {
        plugins: [react()],
        define: {
            'import.meta.env.BACKEND_URL': JSON.stringify(backendURL),
        },
        test: {
            environment: 'jsdom',
            setupFiles: './src/test/setup.js',
            coverage: {
                provider: 'v8',
                reporter: ['text', 'html', 'json'],
            },
        },
    }
})
