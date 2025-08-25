"use client";
import React, { createContext, useContext, useEffect, useState, useRef, ReactNode } from "react";
import {
  EntryDevice,
  DiscoveredDevice,
  // AccessDevice,
} from "../types/status";
import api from "@/lib/api/api";

interface StatusContextType {
  systemName?: string;
  entryDevices: EntryDevice[];
  discoveredDevices: DiscoveredDevice[];
  // accessDevices: AccessDevice[];
  loading: boolean;
  error: string | null;
}

const StatusContext = createContext<StatusContextType | undefined>(undefined);

const LOG_PREFIX = "[StatusContext]";

export const StatusProvider = ({ children }: { children: ReactNode }) => {
  const [systemName, setSystemName] = useState<string | undefined>(undefined);
  const [entryDevices, setEntryDevices] = useState<EntryDevice[]>([]);
  const [discoveredDevices, setDiscoveredDevices] = useState<DiscoveredDevice[]>([]);
  // const [accessDevices, setAccessDevices] = useState<AccessDevice[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const etagRef = useRef<number>(0);

  useEffect(() => {
    let cancelled = false;

    const poll = async () => {
      while (!cancelled) {
        setLoading(true);
        const { data, error } = await api.status.LONG_POLL(35, etagRef.current);

        if (data) {
          setSystemName(data.systemName);
          setEntryDevices(data.entryDevices || []);
          setDiscoveredDevices(data.discoveredDevices || []);
          // setAccessDevices(data.accessDevices || []);
          etagRef.current = data.etag;
          setError(null);
          console.log(LOG_PREFIX, "System Status Updated", data);
        }

        if (error && error.status !== 304) {
          console.warn(LOG_PREFIX, "Polling error:", error.message, "status:", error.status);
          setError("Hub unreachable");
          // backoff before retry
          await new Promise((resolve) => setTimeout(resolve, 5000));
        }

        setLoading(false);
      }
    };

    poll();

    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <StatusContext.Provider
      value={{
        systemName,
        entryDevices,
        discoveredDevices,
        // accessDevices,
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
