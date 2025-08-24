"use client";
import React, { createContext, useContext, useEffect, useState, ReactNode } from "react";
import { StatusResponse, EntryDevice, DiscoveredDevice, AccessDevice } from "../types/status";

interface StatusContextType {
  systemName?: string;
  entryDevices: EntryDevice[];
  discoveredDevices: DiscoveredDevice[];
  accessDevices: AccessDevice[];
  loading: boolean;
  error: string | null;
}

const StatusContext = createContext<StatusContextType | undefined>(undefined);

export const StatusProvider = ({ children }: { children: ReactNode }) => {
  const [systemName, setSystemName] = useState<string | undefined>(undefined);
  const [entryDevices, setEntryDevices] = useState<EntryDevice[]>([]);
  const [discoveredDevices, setDiscoveredDevices] = useState<DiscoveredDevice[]>([]);
  const [accessDevices, setAccessDevices] = useState<AccessDevice[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStatus = async () => {
    try {
      setLoading(true);
      const res = await fetch("/management/status");
      if (!res.ok) {
        throw new Error(`Failed to fetch status: ${res.status}`);
      }
      const data: StatusResponse = await res.json();
      console.log("Fetched status:", data);

      // Split into individual state slices
      setSystemName(data.system_name);
      setEntryDevices(data.entry_devices || []);
      setDiscoveredDevices(data.discovered_devices || []);
    //   setAccessDevices(data.access_devices || []);
      setError(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  // Fetch immediately and then every 10s
  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchStatus, 10000);
    return () => clearInterval(interval);
  }, []);

  return (
    <StatusContext.Provider
      value={{
        systemName,
        entryDevices,
        discoveredDevices,
        accessDevices,
        loading,
        error,
      }}
    >
      {children}
    </StatusContext.Provider>
  );
};

export const useStatusContext = (): StatusContextType => {
  const ctx = useContext(StatusContext);
  if (!ctx) {
    throw new Error("useStatus must be used within a StatusProvider");
  }
  return ctx;
};
