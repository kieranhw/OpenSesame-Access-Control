"use client";
import Image from "next/image";

function Header() {
  return (
    <header className="border-divider bg-card flex-none border-b px-4 py-2">
      <div className="flex h-8 items-center gap-3">
        <Image
          src="/sesame.png"
          alt="OpenSesame Logo"
          width={30}
          height={30}
          draggable="false"
        />
        <h1 className="text-foreground font-semibold">
          OpenSesame Access Control
        </h1>
      </div>
    </header>
  );
}

export default Header;
