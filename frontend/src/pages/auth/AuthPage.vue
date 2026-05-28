<template>
  <!-- Standard full viewport scrolling container instead of q-page -->
  <div class="fullscreen scroll row items-center justify-center q-pa-md bg-gradient">
    <div class="col-12 col-sm-8 col-md-4">
      <q-card class="rounded-borders q-pa-md shadow-2 bg-white">
        <!-- Brand Title -->
        <q-card-section class="text-center q-pb-none">
          <div class="text-h5 text-weight-bold text-primary">MUNDIAL 2026</div>
          <div class="text-subtitle2 text-grey-6 q-mt-xs">Prode Prediction Platform</div>
        </q-card-section>

        <!-- Auth Mode Toggle Tabs -->
        <q-card-section class="q-pb-none">
          <q-tabs
            v-model="tabMode"
            dense
            class="text-grey"
            active-color="primary"
            indicator-color="primary"
            align="justify"
            narrow-indicator
          >
            <q-tab name="login" label="Login" />
            <q-tab name="register" label="Register" />
          </q-tabs>
          <q-separator class="q-mt-sm" />
        </q-card-section>

        <!-- Email/Password Forms -->
        <q-card-section class="q-pt-md">
          <q-tab-panels v-model="tabMode" animated class="bg-transparent">
            <!-- Login Panel -->
            <q-tab-panel name="login" class="q-pa-none">
              <q-form @submit.prevent="handleLogin" class="q-gutter-y-md">
                <q-input
                  v-model="email"
                  type="email"
                  label="Email"
                  outlined
                  dense
                  :rules="[(val) => !!val || 'Email is required']"
                  hide-bottom-space
                />
                <q-input
                  v-model="password"
                  type="password"
                  label="Password"
                  outlined
                  dense
                  :rules="[(val) => !!val || 'Password is required']"
                  hide-bottom-space
                />
                <q-btn
                  label="Log In"
                  type="submit"
                  color="primary"
                  class="full-width q-py-sm"
                  unelevated
                  :loading="loading"
                />
              </q-form>
            </q-tab-panel>

            <!-- Registration Panel -->
            <q-tab-panel name="register" class="q-pa-none">
              <q-form @submit.prevent="handleRegister" class="q-gutter-y-md">
                <q-input
                  v-model="username"
                  type="text"
                  label="Display Username"
                  outlined
                  dense
                  :rules="[(val) => !!val || 'Username is required']"
                  hide-bottom-space
                />
                <q-input
                  v-model="email"
                  type="email"
                  label="Email"
                  outlined
                  dense
                  :rules="[(val) => !!val || 'Email is required']"
                  hide-bottom-space
                />
                <q-input
                  v-model="password"
                  type="password"
                  label="Password"
                  outlined
                  dense
                  :rules="[
                    (val) => !!val || 'Password is required',
                    (val) => val.length >= 8 || 'Must be at least 8 characters',
                  ]"
                  hide-bottom-space
                />
                <q-btn
                  label="Create Account"
                  type="submit"
                  color="primary"
                  class="full-width q-py-sm"
                  unelevated
                  :loading="loading"
                />
              </q-form>
            </q-tab-panel>
          </q-tab-panels>
        </q-card-section>

        <!-- OAuth SSO Divider -->
        <q-card-section class="q-py-none row items-center justify-center q-gutter-x-sm">
          <q-separator class="col" />
          <span class="text-caption text-grey-5">OR SIGN IN WITH</span>
          <q-separator class="col" />
        </q-card-section>

        <!-- Social login providers -->
        <q-card-section class="q-pt-md">
          <q-btn
            outline
            color="grey-8"
            class="full-width q-py-sm bg-grey-1"
            unelevated
            @click="handleGoogleSignIn"
          >
            <div class="row items-center justify-center q-gutter-x-sm">
              <q-icon name="login" size="xs" color="red" />
              <span class="text-weight-bold">Google Auth Connection</span>
            </div>
          </q-btn>
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { useAuthStore } from '@/stores/auth';

console.log('PROCESS: ', process.env);
console.log(process.env.VITE_APP_NAME);
console.log(process.env.VITE_BACKEND_URL);

const router = useRouter();
const $q = useQuasar();
const authStore = useAuthStore();

const tabMode = ref<'login' | 'register'>('login');
const email = ref('');
const password = ref('');
const username = ref('');
const loading = ref(false);

async function handleLogin() {
  loading.value = true;
  try {
    await authStore.loginEmail(email.value, password.value);
    $q.notify({ type: 'positive', message: 'Welcome back!' });
    void router.push('/app/matches');
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Invalid credentials';
    $q.notify({ type: 'negative', message: msg });
  } finally {
    loading.value = false;
  }
}

async function handleRegister() {
  loading.value = true;
  try {
    await authStore.register(email.value, password.value, username.value);
    $q.notify({ type: 'positive', message: 'Account registered successfully!' });
    void router.push('/app/matches');
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Registration failed';
    $q.notify({ type: 'negative', message: msg });
  } finally {
    loading.value = false;
  }
}

async function handleGoogleSignIn() {
  try {
    await authStore.loginGoogle();
    $q.notify({ type: 'positive', message: 'Google Authentication complete!' });
    void router.push('/app/matches');
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Google OAuth failed';
    $q.notify({ type: 'negative', message: msg });
  }
}
</script>

<style scoped>
.bg-gradient {
  background: linear-gradient(135deg, #1a2d5a 0%, #0d1b2a 100%);
}
.rounded-borders {
  border-radius: 12px;
}
</style>
