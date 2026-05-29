<template>
  <!-- Dynamic page background class -->
  <q-page
    class="q-pa-md row justify-center items-center"
    :class="$q.dark.isActive ? 'bg-grey-11' : 'bg-grey-1'"
  >
    <div class="col-12 col-sm-8 col-md-5">
      <!-- Dynamic card theme and background classes -->
      <q-card
        flat
        bordered
        :dark="$q.dark.isActive"
        :class="$q.dark.isActive ? 'bg-grey-10' : 'bg-white'"
        class="rounded-borders q-pa-md shadow-1"
      >
        <q-card-section class="text-center">
          <q-avatar size="80px" color="primary" text-color="white" class="q-mb-md">
            {{ username.charAt(0).toUpperCase() }}
          </q-avatar>
          <div class="text-h6 text-weight-bold text-primary">Edicion de Perfil</div>
          <div class="text-caption text-grey-6">Actualiza tus datos personales</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit.prevent="saveProfile" class="q-gutter-y-md">
            <q-input
              v-model="username"
              outlined
              label="Username"
              dense
              :rules="[(val) => !!val || 'Username is required']"
              hide-bottom-space
            />

            <q-input
              v-model="avatarUrl"
              outlined
              label="Avatar URL (Optional)"
              dense
              hide-bottom-space
            />

            <q-btn
              label="Guardar cambios"
              type="submit"
              color="primary"
              unelevated
              class="full-width q-py-sm"
              :loading="saving"
            />
          </q-form>
        </q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { Notify } from 'quasar';
import { pb } from '@/boot/pocketbase';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();

const username = ref('');
const avatarUrl = ref('');
const saving = ref(false);

watch(
  () => authStore.user,
  (val) => {
    if (val) {
      username.value = val.username || '';
      avatarUrl.value = val.avatar_url || '';
    }
  },
  { immediate: true },
);

async function saveProfile() {
  if (!authStore.user?.id || !username.value) return;
  saving.value = true;
  try {
    const updated = await pb.collection('users').update(authStore.user.id, {
      username: username.value,
      avatar_url: avatarUrl.value,
    });

    authStore.user = updated;

    Notify.create({
      type: 'positive',
      message: 'Profile successfully updated!',
    });
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Error updating profile data';
    Notify.create({
      type: 'negative',
      message: msg,
    });
  } finally {
    saving.value = false;
  }
}
</script>
