
interface DeviceStatusCellProps {
  isOnline: boolean;
}

export function DeviceStatusCell({ isOnline }: DeviceStatusCellProps) {
  return (
    <div className="flex items-center gap-2">
      <span className="relative flex size-3">
        <span
          className={`absolute inline-flex h-full w-full rounded-full ${
            isOnline ? "bg-green-500" : "bg-destructive animate-ping"
          } opacity-75`}
        ></span>
        <span
          className={`relative inline-flex size-3 rounded-full ${isOnline ? "bg-green-500" : "bg-destructive"}`}
        ></span>
      </span>
      <p>{isOnline ? "Online" : "Offline"}</p>
    </div>
  );
}
