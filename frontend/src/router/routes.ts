import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/pages/auth/AuthPage.vue'),
  },
  {
    path: '/register',
    redirect: '/',
  },
  {
    path: '/app',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: '/app/matches' },
      { path: 'matches', component: () => import('@/pages/matches/MatchListPage.vue') },
      { path: 'bracket', component: () => import('@/pages/bracket/BracketPage.vue') }, // New
      { path: 'groups', component: () => import('@/pages/groups/GroupsPage.vue') },
      { path: 'groups/:id', component: () => import('@/pages/groups/GroupDetailPage.vue') },
      { path: 'profile', component: () => import('@/pages/profile/ProfilePage.vue') }, // New
      { path: 'info', component: () => import('@/pages/info/InfoPage.vue') },
      { path: 'predictions-breakdown', component: () => import('@/pages/predictions/PredictionsBreakdownPage.vue') },
      {
        path: 'admin/results',
        component: () => import('@/pages/admin/AdminResultsPage.vue'),
        meta: { requiresAdmin: true },
      },
      {
        path: 'admin/matches',
        component: () => import('@/pages/admin/AdminMatchesPage.vue'),
        meta: { requiresAdmin: true },
      },
    ],
  },
  {
    path: '/:catchAll(.*)*',
    component: () => import('@/pages/ErrorNotFound.vue'),
  },
];

export default routes;
