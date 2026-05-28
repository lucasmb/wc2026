import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/pages/auth/AuthPage.vue'), // Unified View
  },
  {
    path: '/register',
    redirect: '/', // Simple redirect back to unified view
  },
  {
    path: '/app',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/app/matches' },
      { path: 'matches', component: () => import('@/pages/matches/MatchListPage.vue') },
      { path: 'groups', component: () => import('@/pages/groups/GroupsPage.vue') },
      { path: 'groups/:id', component: () => import('@/pages/groups/GroupDetailPage.vue') },
    ],
  },
  {
    path: '/:catchAll(.*)*',
    component: () => import('@/pages/ErrorNotFound.vue'),
  },
];

export default routes;
