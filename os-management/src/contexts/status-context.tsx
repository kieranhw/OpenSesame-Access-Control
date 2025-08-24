"use client";
import React, { createContext, useContext, useEffect, useState, useRef, ReactNode } from "react";
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

  const etagRef = useRef<number>(0);

  useEffect(() => {
    let cancelled = false;

    const poll = async () => {
      while (!cancelled) {
        try {
          setLoading(true);
          const controller = new AbortController();
          const timer = setTimeout(() => controller.abort(), 30_000);
          const url = `/management/status?timeout=35&etag=${etagRef.current}`;
          const res = await fetch(url, { signal: controller.signal });
          clearTimeout(timer);

          if (res.status === 304) {
            // No changes, just loop again
          }

          if (res.status === 200) {
            const data: StatusResponse = await res.json();

            setSystemName(data.system_name);
            setEntryDevices(data.entry_devices || []);
            setDiscoveredDevices(data.discovered_devices || []);
            // setAccessDevices(data.access_devices || []);

            etagRef.current = data.etag;
            console.log("Response received:", data);
            setError(null);
          } else {
            throw new Error(`Unexpected status: ${res.status}`);
          }
        } catch (err: any) {
          if (err.name === "AbortError") {
            // client aborted before server response
          } else {
            console.warn("Polling error: ", err?.message || err);
            setError("Hub unreachable");

            // retry after 5 sec
            await new Promise((resolve) => setTimeout(resolve, 5000));
          }
        } finally {
          setLoading(false);
        }
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
