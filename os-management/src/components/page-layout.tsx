export async function PageLayout({ children }: { children: React.ReactNode }) {
  return <div className="mx-auto flex w-full max-w-6xl flex-1 flex-col gap-4 p-4">{children}</div>;
}
