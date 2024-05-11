import React, { createContext, useState, useContext } from 'react';

const DarkModeContext = createContext();

export const DarkModeProvider = ({ children }) => {
  const [mode, setMode] = useState(true);

  return (
    <DarkModeContext.Provider value={{ mode, setMode }}>
      {children}
    </DarkModeContext.Provider>
  );
};

export const useDarkMode = () => useContext(DarkModeContext);
