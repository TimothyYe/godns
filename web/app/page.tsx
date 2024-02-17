import { Stat } from '@/components/stat';

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-between pt-10 max-w-screen-xl">
      <div className="p-10">
        <Stat />
      </div>
    </main>
  );
}
