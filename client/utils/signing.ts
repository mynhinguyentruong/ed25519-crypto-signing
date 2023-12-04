import fs from "fs";
import sshpk from "sshpk";

export function createSignature(data: string = "some data") {
  // data is a message to be signed
  const keyPriv = fs.readFileSync("../id_ed25519");

  const key = sshpk.parsePrivateKey(keyPriv, "pem");

  /* Sign some data with the key */
  let s = key.createSign("sha512");
  s.update(data);
  const signature = s.sign();
  console.log({ signature });

  // get signature part and encode to base64 string
  const base64String = signature.part.sig.data.toString("base64");
  console.log({ base64String });

  return base64String;
}
