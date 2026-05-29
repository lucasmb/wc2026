<template>
  <q-layout view="lHh Lpr lFf" :class="$q.dark.isActive ? 'bg-grey-11' : 'bg-grey-1'">
    <!-- Desktop Side Navigation Drawer -->
    <q-drawer
      v-if="$q.platform.is.desktop"
      v-model="drawer"
      show-if-above
      bordered
      :dark="$q.dark.isActive"
      :class="$q.dark.isActive ? 'bg-grey-10' : 'bg-white'"
    >
      <q-scroll-area class="fit">
        <div
          class="q-px-md q-py-lg text-weight-bold text-h6 text-center border-bottom"
          :class="$q.dark.isActive ? 'text-white' : 'text-primary'"
        >
          MUNDIAL 2026
        </div>
        <q-list class="q-py-md">
          <q-item
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            clickable
            v-ripple
            class="q-mx-sm q-my-xs rounded-borders"
            :active-class="
              $q.dark.isActive
                ? 'bg-blue-10 text-white text-weight-bold'
                : 'bg-blue-1 text-primary text-weight-bold'
            "
          >
            <q-item-section avatar>
              <q-icon :name="link.icon" />
            </q-item-section>
            <q-item-section>{{ link.label }}</q-item-section>
          </q-item>
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <!-- App Header Bar -->
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <q-btn
          v-if="$q.platform.is.desktop"
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="drawer = !drawer"
        />

        <!-- Emblem Vectorial Oficial de la Copa del Mundo de Oro -->
        <svg
          viewBox="0 0 24 24"
          width="28"
          height="28"
          fill="#e5b83e"
          xmlns="http://www.w3.org/2000/svg"
          class="q-mr-sm"
        >
          <path
            d="M18,2H6A1,1,0,0,0,5,3V7a6,6,0,0,0,4.82,5.89,3.91,3.91,0,0,0,1.18,1.93V18H9a1,1,0,0,0,0,2h6a1,1,0,0,0,0,-2H13V14.82a3.91,3.91,0,0,0,1.18,-1.93A6,6,0,0,0,19,7V3A1,1,0,0,0,18,2ZM7,7V4H17V7a4,4,0,0,1-8,0ZM12,14a2,2,0,1,1,2,-2A2,2,0,0,1,12,14Z"
          />
        </svg>

        <q-toolbar-title class="text-weight-bold">Prode Mundial 2026</q-toolbar-title>

        <!-- Theme Switcher Button -->
        <q-btn
          flat
          round
          dense
          class="q-mr-sm"
          :icon="$q.dark.isActive ? 'dark_mode' : 'light_mode'"
          :color="$q.dark.isActive ? 'secondary' : 'white'"
          @click="toggleDarkMode"
        />

        <q-btn flat round icon="logout" @click="handleLogout" />
      </q-toolbar>
    </q-header>

    <!-- Router View Viewport Container -->
    <q-page-container>
      <router-view />
    </q-page-container>

    <!-- Mobile Screen Bottom Toolbar Navigation Tabs -->
    <q-footer
      v-if="$q.platform.is.mobile"
      class="border-top shadow-up-3"
      :class="$q.dark.isActive ? 'bg-grey-10 text-white' : 'bg-white text-grey-7'"
      style="bottom: 0px"
    >
      <q-tabs
        dense
        align="justify"
        indicator-color="transparent"
        :dark="$q.dark.isActive"
        active-color="primary"
      >
        <q-route-tab
          v-for="link in navLinks"
          :key="link.to"
          :to="link.to"
          :icon="link.icon"
          :label="link.label"
          exact
          content-class="text-caption"
        />
      </q-tabs>
    </q-footer>
  </q-layout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useQuasar } from 'quasar';

const drawer = ref(false);
const authStore = useAuthStore();
const router = useRouter();
const $q = useQuasar();

const navLinks = [
  { to: '/app/matches', icon: 'sports_soccer', label: 'Partidos' },
  { to: '/app/groups', icon: 'group', label: 'Grupos' },
  { to: '/app/bracket', icon: 'account_tree', label: 'Llaves' },
  { to: '/app/profile', icon: 'person', label: 'Perfil' },
];

function toggleDarkMode() {
  $q.dark.set(!$q.dark.isActive);
}

function handleLogout() {
  authStore.logout();
  void router.push('/');
}
</script>

<style scoped>
.rounded-borders {
  border-radius: 8px;
}
.border-bottom {
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}
</style>
