"use client";
import React, { createContext, useContext, useEffect, useState, useRef, ReactNode, useCallback } from "react";
import api from "@/lib/api/api";
import { DeviceFactory } from "@/domain/device/device-factory";
import { EntryDevice } from "@/domain/device/entry-device";
import { EntryDeviceData } from "@/domain/common/device";

interface StatusContextType {
  systemName?: string;
  entryDevices: EntryDevice[];
  // accessDevices: AccessDevice[];
  loading: boolean;
  bumpState: () => void;
  error: string | null;
}

const StatusContext = createContext<StatusContextType | undefined>(undefined);

const LOG_PREFIX = "[StatusContext]";

export const StatusProvider = ({ children }: { children: ReactNode }) => {
  const [systemName, setSystemName] = useState<string | undefined>(undefined);
  const [entryDevices, setEntryDevices] = useState<EntryDevice[]>([]);
  // const [accessDevices, setAccessDevices] = useState<AccessDevice[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const etagRef = useRef<number>(0);
  const cancelledRef = useRef<boolean>(false);

  function bumpState() {
    setEntryDevices((prev) => [...prev]);
  }

  const parseEntryDevices = useCallback((incoming: EntryDeviceData[] | undefined, prev: EntryDevice[]) => {
    const byId = new Map(prev.map((d) => [d.id, d]));
    const next: EntryDevice[] = [];

    for (const dto of incoming ?? []) {
      const existing = byId.get(dto.id);
      if (existing) {
        existing.data = dto;
        next.push(existing);
      } else {
        next.push(DeviceFactory.create(dto) as EntryDevice);
      }
    }

    return next;
  }, []);

  const poll = useCallback(async () => {
    while (!cancelledRef.current) {
      setLoading(true);

      const { data, error } = await api.status.LONG_POLL(35, etagRef.current);

      if (data) {
        setSystemName(data.systemName);
        setEntryDevices((prev) => parseEntryDevices(data.entryDevices, prev));
        // setAccessDevices(...)
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
  }, [parseEntryDevices]);

  useEffect(() => {
    cancelledRef.current = false;
    // fire and forget
    void poll();

    return () => {
      cancelledRef.current = true;
    };
  }, [poll]);

  return (
    <StatusContext.Provider
      value={{
        systemName,
        entryDevices,
        // accessDevices,
        loading,
        error,
        bumpState,
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
