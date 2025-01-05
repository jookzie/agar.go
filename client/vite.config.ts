import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

export default defineConfig(({ command, mode }) => {
  if (command === 'serve') {
    // Development server configuration
    return {
	  server: {
		port: 3000,
		host: '0.0.0.0',
		watch: {
		  usePolling: true
		}
	  }
    };
  } else if (command === 'build') {
    // Build configuration
    return {
	  server: {
		port: 80,
		host: '0.0.0.0',
		watch: {
		  usePolling: true
		}
	  }
    };
  }
});
