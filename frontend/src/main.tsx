import { QueryClientProvider } from '@tanstack/react-query';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { Toaster } from 'sonner';
import App from './App.tsx';
import { ThemeProvider } from './hooks';
import './i18n';
import './index.css';
import { Devtools } from './query-client/query-client-dev-tools.tsx';
import { queryClient } from './query-client/query-client.ts';

const root = document.getElementById('root');
if (root != null) {
  createRoot(root).render(
    <StrictMode>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider>
          <Toaster />
          <App />
        </ThemeProvider>
        <Devtools />
      </QueryClientProvider>
    </StrictMode>,
  );
} else {
  console.error("couldn't find root element");
}
