export interface response {
  erros: string;
  evaluated: string;
}

export const parse = async (source: string): Promise<response> => {
  const res = await fetch("http://localhost:1323/parse", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      source: source
    })
  })

  if (!res.ok) {
    const err = await res.json();
    throw new Error(err)
  }

  return await res.json();
}