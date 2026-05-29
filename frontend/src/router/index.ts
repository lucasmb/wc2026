import { route } from 'quasar/wrappers';
import {
  createMemoryHistory,
  createRouter,
  createWebHashHistory,
  createWebHistory,
} from 'vue-router';
import routes from './routes';
import { useAuthStore } from '@/stores/auth';

export default route(function (/* { store, ssrContext } */) {
  const createHistory = process.env.SERVER
    ? createMemoryHistory
    : process.env.VUE_ROUTER_MODE === 'history'
      ? createWebHashHistory
      : createWebHistory;

  const Router = createRouter({
    scrollBehavior: () => ({ left: 0, top: 0 }),
    routes,
    history: createHistory(process.env.VUE_ROUTER_BASE),
  });

  // Robust recursive matched navigation guard checking parent layouts
  Router.beforeEach((to, _from, next) => {
    const authStore = useAuthStore();
    const isAuthRequired = to.matched.some((record) => record.meta.requiresAuth);

    if (isAuthRequired && !authStore.isLoggedIn) {
      next({ path: '/' });
    } else {
      next();
    }
  });

  return Router;
});
