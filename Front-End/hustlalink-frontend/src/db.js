import Dexie from "dexie";

export const db = new Dexie("HustlaLinkDB");

db.version(1).stores({
  users: "id, firstName, secondName, category, phone",
  pending: "++id, type, payload"
});