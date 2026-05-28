<template>
  <q-layout view="lHh Lpr lFf" class="bg-grey-1">
    <!-- Desktop Side Navigation Drawer -->
    <q-drawer
      v-if="$q.platform.is.desktop"
      v-model="drawer"
      show-if-above
      bordered
      :width="240"
      class="bg-white"
    >
      <q-scroll-area class="fit">
        <div
          class="q-px-md q-py-lg text-weight-bold text-primary text-h6 text-center border-bottom"
        >
          {{ appName }}
        </div>
        <q-list class="q-py-md">
          <q-item
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            clickable
            v-ripple
            class="q-mx-sm q-my-xs rounded-borders"
            active-class="bg-blue-1 text-primary text-weight-bold"
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
        <q-toolbar-title class="text-weight-bold">{{ appName }}</q-toolbar-title>
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
      class="bg-white border-top shadow-up-3"
      style="bottom: 0px"
    >
      <q-tabs
        dense
        align="justify"
        class="text-grey-7 text-primary"
        active-color="primary"
        indicator-color="transparent"
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
import { useAuthStore } from 'src/stores/auth';

const drawer = ref(false);
const authStore = useAuthStore();
const router = useRouter();

const appName = process.env.VITE_APP_NAME;
const navLinks = [
  { to: '/app/matches', icon: 'sports_soccer', label: 'Partidos' },
  { to: '/app/groups', icon: 'group', label: 'Grupos' },
];

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
