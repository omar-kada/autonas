import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import App from './App.tsx';
import './i18n';
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from './queryClient/queryClient';
import { Devtools } from './queryClient/queryClientDevTools';

const root = document.getElementById('root');
if (root != null) {
  createRoot(root).render(
    <StrictMode>
      <QueryClientProvider client={queryClient}>
        <App />
        <Devtools />
      </QueryClientProvider>
    </StrictMode>,
  );
} else {
  console.error("couldn't find root element");
}
