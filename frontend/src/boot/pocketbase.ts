import { boot } from 'quasar/wrappers';
import PocketBase from 'pocketbase';

// Determine default backend connection endpoint dynamically
const PB_URL = process.env.VITE_BACKEND_URL;
const pb = new PocketBase(PB_URL);

export default boot(({ app }) => {
  app.config.globalProperties.$pb = pb;
});

export function getFileUrl(recordId: string, filename: string): string {
  if (!filename) return '';
  return `${PB_URL}/api/files/_pb_users_auth_/${recordId}/${filename}`;
}

export { pb, PB_URL };
