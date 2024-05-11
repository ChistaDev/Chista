import React, { createContext, useState, useContext } from 'react';

const PhishingScanExcludeInputContext = createContext();

export const PhishingScanExcludeInputProvider = ({ children }) => {
  const [phishingScanExcludeInput, setPhishingScanExcludeInput] = useState('');

  return (
    <PhishingScanExcludeInputContext.Provider
      value={{ phishingScanExcludeInput, setPhishingScanExcludeInput }}
    >
      {children}
    </PhishingScanExcludeInputContext.Provider>
  );
};

export const usePhishingScanExcludeInput = () =>
  useContext(PhishingScanExcludeInputContext);
