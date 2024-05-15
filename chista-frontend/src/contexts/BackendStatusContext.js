import React, { createContext, useState, useContext } from 'react';

const BackendStatusContext = createContext();

export const BackendStatusProvider = ({ children }) => {
  const [backendStatus, setBackendStatus] = useState(true);

  return (
    <BackendStatusContext.Provider value={{ backendStatus, setBackendStatus }}>
      {children}
    </BackendStatusContext.Provider>
  );
};

export const useBackendStatus = () => useContext(BackendStatusContext);
