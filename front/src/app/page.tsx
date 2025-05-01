import Image from "next/image";
import JsonGenerator from '@/components/JsonGenerator';

export default function Home() {
  return (
    <main className="min-h-screen">
      <JsonGenerator />
    </main>
  );
}
