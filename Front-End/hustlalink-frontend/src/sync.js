import { db } from "./db";

const API_URL = "https://your-api.com/users";

export const syncPending = async () => {
  if (!navigator.onLine) return;

  const pending = await db.pending.toArray();

  for (const item of pending) {
    try {
      if (item.type === "ADD_USER") {
        await fetch(API_URL, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(item.payload)
        });
      }

      await db.pending.delete(item.id);
    } catch {
      console.log("Still failing...");
    }
  }
};