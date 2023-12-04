export default async function Home() {
  const options = {
    method: "GET",
    headers: {
      "Magiclip-Signature":
        "YgaFcCSV3TxVVyz25r8PRztnQnfDXHV8tsaezeDeijsKXa8C1SmbMI2qRgSx6aobeiqKBouHY7TCBrNpxQDQDA==",
      "Magiclip-Message": "some data",
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
