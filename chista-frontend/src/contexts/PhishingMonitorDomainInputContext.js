import React, { createContext, useState, useContext } from 'react';

const PhishingMonitorDomainInputContext = createContext();

export const PhishingMonitorDomainInputProvider = ({ children }) => {
  const [phishingMonitorDomainInput, setPhishingMonitorDomainInput] =
    useState('');

  return (
    <PhishingMonitorDomainInputContext.Provider
      value={{ phishingMonitorDomainInput, setPhishingMonitorDomainInput }}
    >
      {children}
    </PhishingMonitorDomainInputContext.Provider>
  );
};

export const usePhishingMonitorDomainInput = () =>
  useContext(PhishingMonitorDomainInputContext);
