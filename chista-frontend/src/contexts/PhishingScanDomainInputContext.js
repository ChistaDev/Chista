import React, { createContext, useState, useContext } from 'react';

const PhishingScanDomainInputContext = createContext();

export const PhishingScanDomainInputProvider = ({ children }) => {
  const [phishingScanDomainInput, setPhishingScanDomainInput] = useState('');

  return (
    <PhishingScanDomainInputContext.Provider
      value={{ phishingScanDomainInput, setPhishingScanDomainInput }}
    >
      {children}
    </PhishingScanDomainInputContext.Provider>
  );
};

export const usePhishingScanDomainInput = () =>
  useContext(PhishingScanDomainInputContext);
