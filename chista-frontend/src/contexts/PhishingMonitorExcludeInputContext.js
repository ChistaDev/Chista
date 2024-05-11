import React, { createContext, useState, useContext } from 'react';

const PhishingMonitorExcludeInputContext = createContext();

export const PhishingMonitorExcludeInputProvider = ({ children }) => {
  const [phishingMonitorExcludeInput, setPhishingMonitorExcludeInput] =
    useState('');

  return (
    <PhishingMonitorExcludeInputContext.Provider
      value={{ phishingMonitorExcludeInput, setPhishingMonitorExcludeInput }}
    >
      {children}
    </PhishingMonitorExcludeInputContext.Provider>
  );
};

export const usePhishingMonitorExcludeInput = () =>
  useContext(PhishingMonitorExcludeInputContext);
