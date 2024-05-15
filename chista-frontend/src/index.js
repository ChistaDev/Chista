import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import './index.css';
import { DarkModeProvider } from './contexts/DarkModeContext';
import { ToastMessageProvider } from './contexts/ToastMessageContext';
import { BackendStatusProvider } from './contexts/BackendStatusContext';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <BrowserRouter>
    <DarkModeProvider>
      <ToastMessageProvider>
        <BackendStatusProvider>
          <App />
        </BackendStatusProvider>
      </ToastMessageProvider>
    </DarkModeProvider>
  </BrowserRouter>
);
