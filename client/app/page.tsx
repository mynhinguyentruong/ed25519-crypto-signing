import { createSignature } from "@/utils/signing";

export default async function Home() {
  const message = "input a message";

  const sig = createSignature(message);

  const options = {
    method: "GET",
    headers: {
      "Magiclip-Signature": sig,
      "Magiclip-Message": message,
    },
  };

  const res = await fetch("http://localhost:4242/ping", options);
  const data = await res.json();

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      {JSON.stringify(data)}
    </main>
  );
}
