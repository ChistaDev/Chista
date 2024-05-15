import React, { createContext, useState, useContext } from 'react';

const ToastMessageContext = createContext();

export const ToastMessageProvider = ({ children }) => {
  const [openToast, setOpenToast] = useState(false);
  const [severity, setSeverity] = useState('');
  const [toastContent, setToastContent] = useState('');

  return (
    <ToastMessageContext.Provider
      value={{
        openToast,
        setOpenToast,
        severity,
        setSeverity,
        toastContent,
        setToastContent,
      }}
    >
      {children}
    </ToastMessageContext.Provider>
  );
};

export const useToastMessage = () => useContext(ToastMessageContext);
