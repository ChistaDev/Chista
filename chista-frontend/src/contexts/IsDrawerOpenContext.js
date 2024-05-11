import React, { createContext, useState, useContext } from 'react';

const IsDrawerOpenContext = createContext();

export const IsDrawerOpenProvider = ({ children }) => {
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  return (
    <IsDrawerOpenContext.Provider value={{ isDrawerOpen, setIsDrawerOpen }}>
      {children}
    </IsDrawerOpenContext.Provider>
  );
};

export const useIsDrawerOpen = () => useContext(IsDrawerOpenContext);
