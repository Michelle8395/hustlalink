import { db } from "./db";

const API_URL = "http://localhost:9000/users";

// Check server health
const isServerUp = async () => {
  try {
    const res = await fetch(API_URL, { method: "HEAD" });
    return res.ok;
  } catch {
    return false;
  }
};

// fetch users
export const getUsers = async () => {
  const online = navigator.onLine;
  const serverUp = online && await isServerUp();

  if (serverUp) {
    try {
      const res = await fetch(API_URL);
      const data = await res.json();

      await db.users.clear();
      await db.users.bulkAdd(data);

      return data;
    } catch {
        console.log('server offline fetching from indexeddb')
      return await db.users.toArray();
    }
  } else {
    return await db.users.toArray();
  }
};

// post if server online
export const addUser = async (user) => {
  const online = navigator.onLine;

  if (online) {
    try {
      const res = await fetch(API_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(user)
      });

      if (!res.ok) throw new Error();

      return await res.json();
    } catch {
      await db.pending.add({ type: "ADD_USER", payload: user });
    }
  } else {
    await db.pending.add({ type: "ADD_USER", payload: user });
  }
};