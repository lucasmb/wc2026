import { boot } from 'quasar/wrappers';
import PocketBase from 'pocketbase';

// Determine default backend connection endpoint dynamically
const PB_URL = process.env.VITE_BACKEND_URL;
const pb = new PocketBase(PB_URL);

export default boot(({ app }) => {
  app.config.globalProperties.$pb = pb;
});

export { pb, PB_URL };
