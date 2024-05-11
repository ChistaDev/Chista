import React, { createContext, useState, useContext } from 'react';

const DataLeakMonitorDomainInputContext = createContext();

export const DataLeakMonitorDomainInputProvider = ({ children }) => {
  const [dataLeakMonitorDomainInput, setDataLeakMonitorDomainInput] =
    useState('');

  return (
    <DataLeakMonitorDomainInputContext.Provider
      value={{ dataLeakMonitorDomainInput, setDataLeakMonitorDomainInput }}
    >
      {children}
    </DataLeakMonitorDomainInputContext.Provider>
  );
};

export const useDataLeakMonitorDomainInput = () =>
  useContext(DataLeakMonitorDomainInputContext);
