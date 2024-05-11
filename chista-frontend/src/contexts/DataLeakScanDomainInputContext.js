import React, { createContext, useState, useContext } from 'react';

const DataLeakScanDomainInputContext = createContext();

export const DataLeakScanDomainInputProvider = ({ children }) => {
  const [dataLeakScanDomainInput, setDataLeakScanDomainInput] = useState('');

  return (
    <DataLeakScanDomainInputContext.Provider
      value={{ dataLeakScanDomainInput, setDataLeakScanDomainInput }}
    >
      {children}
    </DataLeakScanDomainInputContext.Provider>
  );
};

export const useDataLeakScanDomainInput = () =>
  useContext(DataLeakScanDomainInputContext);
