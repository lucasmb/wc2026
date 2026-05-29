<template>
  <!-- Contenedor general que abarca todo el viewport con scroll fluido -->
  <div class="fullscreen scroll row items-center justify-center q-pa-md bg-gradient">
    <div class="col-12 col-sm-8 col-md-4">
      <q-card class="rounded-borders q-pa-md shadow-2 bg-white text-grey-9">
        <!-- Título de la Marca -->
        <q-card-section class="text-center q-pb-none">
          <div class="text-h5 text-weight-bolder text-primary">PRODE MUNDIAL 2026</div>
          <div class="text-subtitle2 text-grey-6 q-mt-xs text-weight-medium">
            Plataforma de Predicciones
          </div>
        </q-card-section>

        <!-- Selector de Modo de Pestañas -->
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
            <q-tab name="login" label="Iniciar Sesión" />
            <q-tab name="register" label="Crear Cuenta" />
          </q-tabs>
          <q-separator class="q-mt-sm" />
        </q-card-section>

        <!-- Formularios de Correo y Contraseña -->
        <q-card-section class="q-pt-md">
          <q-tab-panels v-model="tabMode" animated class="bg-transparent">
            <!-- Panel de Inicio de Sesión -->
            <q-tab-panel name="login" class="q-pa-none">
              <q-form @submit.prevent="handleLogin" class="q-gutter-y-md">
                <q-input
                  v-model="email"
                  type="email"
                  label="Correo Electrónico"
                  outlined
                  dense
                  :rules="[(val) => !!val || 'El correo electrónico es requerido']"
                  hide-bottom-space
                />
                <q-input
                  v-model="password"
                  type="password"
                  label="Contraseña"
                  outlined
                  dense
                  :rules="[(val) => !!val || 'La contraseña es requerida']"
                  hide-bottom-space
                />
                <q-btn
                  label="Ingresar"
                  type="submit"
                  color="primary"
                  class="full-width q-py-sm text-weight-bold"
                  unelevated
                  :loading="loading"
                />
              </q-form>
            </q-tab-panel>

            <!-- Panel de Registro -->
            <q-tab-panel name="register" class="q-pa-none">
              <q-form @submit.prevent="handleRegister" class="q-gutter-y-md">
                <q-input
                  v-model="username"
                  type="text"
                  label="Nombre de Usuario"
                  outlined
                  dense
                  :rules="[(val) => !!val || 'El nombre de usuario es requerido']"
                  hide-bottom-space
                />
                <q-input
                  v-model="email"
                  type="email"
                  label="Correo Electrónico"
                  outlined
                  dense
                  :rules="[(val) => !!val || 'El correo electrónico es requerido']"
                  hide-bottom-space
                />
                <q-input
                  v-model="password"
                  type="password"
                  label="Contraseña"
                  outlined
                  dense
                  :rules="[
                    (val) => !!val || 'La contraseña es requerida',
                    (val) => val.length >= 8 || 'Debe contener al menos 8 caracteres',
                  ]"
                  hide-bottom-space
                />
                <q-btn
                  label="Registrar Cuenta"
                  type="submit"
                  color="primary"
                  class="full-width q-py-sm text-weight-bold"
                  unelevated
                  :loading="loading"
                />
              </q-form>
            </q-tab-panel>
          </q-tab-panels>
        </q-card-section>

        <!-- Divisor Social SSO -->
        <q-card-section class="q-py-none row items-center justify-center q-gutter-x-sm">
          <q-separator class="col" />
          <span class="text-caption text-grey-5 text-weight-bold">O CONECTAR CON</span>
          <q-separator class="col" />
        </q-card-section>

        <!-- Botón de Conexión de Google de Alta Fidelidad -->
        <q-card-section class="q-pt-md">
          <q-btn
            outline
            color="grey-4"
            class="full-width q-py-sm bg-grey-1 shadow-1 google-btn-transition"
            unelevated
            @click="handleGoogleSignIn"
          >
            <div class="row items-center justify-center q-gutter-x-sm text-grey-8">
              <!-- Emblema Vectorial Oficial de Google de 4 Colores -->
              <svg viewBox="0 0 24 24" width="18" height="18" xmlns="http://www.w3.org/2000/svg">
                <path
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                  fill="#4285F4"
                />
                <path
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                  fill="#34A853"
                />
                <path
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.06H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.94l2.85-2.22.81-.63z"
                  fill="#FBBC05"
                />
                <path
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.06l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                  fill="#EA4335"
                />
              </svg>
              <span class="text-weight-bold">Iniciar Sesión con Google</span>
            </div>
          </q-btn>
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { Notify } from 'quasar';
import { pb } from '@/boot/pocketbase';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

// 1. Capture invite from URL and persist in SessionStorage to survive OAuth redirects
let inviteGroupId = route.query.invite as string | undefined;
if (inviteGroupId) {
  sessionStorage.setItem('pending_invite_group', inviteGroupId);
} else {
  inviteGroupId = sessionStorage.getItem('pending_invite_group') || undefined;
}

const tabMode = ref<'login' | 'register'>('login');
const email = ref('');
const password = ref('');
const username = ref('');
const loading = ref(false);

async function autoJoinGroupOnLogin() {
  const targetGroupId = sessionStorage.getItem('pending_invite_group') || inviteGroupId;
  if (!targetGroupId || !authStore.user?.id) return;

  try {
    // Check if the user is already a member of this group
    const checkList = await pb.collection('group_members_id').getList(1, 1, {
      filter: `prediction_group = "${targetGroupId}" && user = "${authStore.user.id}"`,
    });

    if (checkList.items.length === 0) {
      await pb.collection('group_members_id').create({
        prediction_group: targetGroupId,
        user: authStore.user.id,
        total_points: 0,
        rank: 1,
      });
      Notify.create({
        type: 'positive',
        message: '¡Te has unido al grupo de predicción de forma automática!',
      });
    }

    // Clean up storage so it does not trigger redundant checks on future logins
    sessionStorage.removeItem('pending_invite_group');
  } catch (err: unknown) {
    console.error('La unión automática al grupo falló:', err);
  }
}

async function handleLogin() {
  loading.value = true;
  try {
    await authStore.loginEmail(email.value, password.value);
    Notify.create({ type: 'positive', message: '¡Bienvenido nuevamente!' });
    await autoJoinGroupOnLogin(); // Safe persistent join
    void router.push('/app/matches');
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Credenciales inválidas';
    Notify.create({ type: 'negative', message: msg });
  } finally {
    loading.value = false;
  }
}

async function handleRegister() {
  loading.value = true;
  try {
    await authStore.register(email.value, password.value, username.value);
    Notify.create({ type: 'positive', message: '¡Cuenta registrada exitosamente!' });
    await autoJoinGroupOnLogin(); // Safe persistent join
    void router.push('/app/matches');
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Error en el registro';
    Notify.create({ type: 'negative', message: msg });
  } finally {
    loading.value = false;
  }
}

async function handleGoogleSignIn() {
  try {
    await authStore.loginGoogle();
    Notify.create({ type: 'positive', message: '¡Autenticación con Google completada!' });
    await autoJoinGroupOnLogin(); // Safe persistent join
    void router.push('/app/matches');
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Error en Google OAuth';
    Notify.create({ type: 'negative', message: msg });
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
.google-btn-transition {
  transition:
    transform 0.1s ease-in-out,
    box-shadow 0.1s ease-in-out;
}
.google-btn-transition:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08) !important;
}
</style>
