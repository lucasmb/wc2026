<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-lg">
      <div class="text-h6 text-weight-bold text-primary">Grupos de Predicción</div>
      <div class="q-gutter-x-sm">
        <q-btn
          label="Unirse con Código"
          color="secondary"
          outline
          dense
          class="q-px-md"
          @click="showJoinDialog"
        />

        <!-- Creation button disabled dynamically after group phase ends -->
        <q-btn
          label="Crear Grupo"
          color="primary"
          unelevated
          dense
          class="q-px-md"
          :disabled="isGroupCreationDisabled"
          @click="showCreateDialog"
        >
          <q-tooltip v-if="isGroupCreationDisabled" class="bg-red text-white">
            La creación de grupos está cerrada porque la Fase de Grupos ya finalizó.
          </q-tooltip>
        </q-btn>
      </div>
    </div>

    <!-- Active Groups Grid -->
    <div v-if="loading" class="row justify-center q-my-xl">
      <q-spinner-dots color="primary" size="40px" />
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-if="groups.length === 0" class="col-12 text-center q-my-xl text-grey-6">
        No eres miembro de ningún grupo de predicción todavía. ¡Únete o crea uno arriba!
      </div>
      <div v-for="group in groups" :key="group.id" class="col-12 col-sm-6 col-md-4">
        <q-card
          flat
          bordered
          class="rounded-borders hover-shadow cursor-pointer"
          @click="goToGroup(group.id)"
        >
          <q-card-section>
            <div class="text-subtitle1 text-weight-bold text-primary text-ellipsis">
              {{ group.name }}
            </div>
            <div class="text-caption text-grey-6 q-mt-xs">
              Creador: {{ group.expand?.owner?.username || 'Sistema' }}
            </div>
          </q-card-section>
          <q-card-actions align="right">
            <q-btn flat label="Ver Posiciones" color="primary" dense />
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <!-- Join Group Dialog -->
    <q-dialog v-model="joinDialog" persistent>
      <q-card style="min-width: 320px">
        <q-card-section class="q-pb-none">
          <div class="text-h6 text-weight-bold">Unirse a Grupo</div>
        </q-card-section>
        <q-card-section class="q-py-md">
          <q-input
            v-model="inviteCode"
            outlined
            label="Código de Invitación"
            placeholder="Ingresar código de 8 dígitos"
            dense
          />
        </q-card-section>
        <q-card-actions align="right" class="q-pt-none q-pb-md q-px-md">
          <q-btn flat label="Cancelar" color="grey-6" v-close-popup />
          <q-btn label="Unirse" color="primary" unelevated :loading="joining" @click="joinGroup" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Create Group Dialog -->
    <q-dialog v-model="createDialog" persistent>
      <q-card style="min-width: 320px">
        <q-card-section class="q-pb-none">
          <div class="text-h6 text-weight-bold">Crear Grupo</div>
        </q-card-section>
        <q-card-section class="q-py-md">
          <q-input
            v-model="newGroupName"
            outlined
            label="Nombre de Grupo"
            placeholder="Ej. Amigos del Trabajo"
            dense
          />
        </q-card-section>
        <q-card-actions align="right" class="q-pt-none q-pb-md q-px-md">
          <q-btn flat label="Cancelar" color="grey-6" v-close-popup />
          <q-btn
            label="Crear"
            color="primary"
            unelevated
            :loading="creating"
            @click="createGroup"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { Notify } from 'quasar';
import { pb } from '@/boot/pocketbase';
import { useAuthStore } from '@/stores/auth';
import type { PredictionGroup } from '@/types';

const router = useRouter();
const authStore = useAuthStore();

const groups = ref<PredictionGroup[]>([]);
const loading = ref(true);

const joinDialog = ref(false);
const inviteCode = ref('');
const joining = ref(false);

const createDialog = ref(false);
const newGroupName = ref('');
const creating = ref(false);

// Active tournament phase state
const currentPhase = ref('group');
const isGroupCreationDisabled = computed(() => currentPhase.value !== 'group');

async function fetchTournamentSettings() {
  try {
    const settings = await pb.collection('settings_id').getOne('tournsettings26');
    currentPhase.value = settings.current_phase || 'group';
  } catch (err: unknown) {
    console.error('Failed fetching tournament settings:', err);
  }
}

async function fetchMyGroups() {
  if (!authStore.user?.id) return;
  try {
    const memberships = await pb.collection('group_members_id').getFullList({
      filter: `user = "${authStore.user.id}"`,
      expand: 'prediction_group.owner',
    });

    groups.value = memberships
      .map((m) => {
        const g = m.expand?.prediction_group as unknown as PredictionGroup;
        return g;
      })
      .filter(Boolean);
  } catch (err: unknown) {
    console.error('Failed fetching groups:', err);
  } finally {
    loading.value = false;
  }
}

function showJoinDialog() {
  inviteCode.value = '';
  joinDialog.value = true;
}

function showCreateDialog() {
  newGroupName.value = '';
  createDialog.value = true;
}

async function joinGroup() {
  if (!inviteCode.value || !authStore.user?.id) return;
  joining.value = true;
  try {
    const groupList = await pb.collection('prediction_groups_id').getList(1, 1, {
      filter: `invite_code = "${inviteCode.value}"`,
    });

    const targetGroup = groupList.items[0];
    if (!targetGroup) {
      throw new Error('Código de invitación inválido o vencido.');
    }

    await pb.collection('group_members_id').create({
      prediction_group: targetGroup.id,
      user: authStore.user.id,
      total_points: 0,
      rank: 1,
    });

    joinDialog.value = false;
    Notify.create({ type: 'positive', message: '¡Te has unido al grupo exitosamente!' });
    void fetchMyGroups();
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Por favor intente nuevamente';
    Notify.create({ type: 'negative', message: message });
  } finally {
    joining.value = false;
  }
}

function generateInviteToken(): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
  let result = '';
  for (let i = 0; i < 8; i++) {
    result += chars.charAt(Math.floor(randValue() * chars.length));
  }
  return result;
}

function randValue(): number {
  return Math.random();
}

async function createGroup() {
  if (!newGroupName.value || !authStore.user?.id) return;
  creating.value = true;
  try {
    const inviteToken = generateInviteToken();

    const createdGroup = await pb.collection('prediction_groups_id').create({
      name: newGroupName.value,
      owner: authStore.user.id,
      invite_code: inviteToken,
      is_public: false,
    });

    await pb.collection('group_members_id').create({
      prediction_group: createdGroup.id,
      user: authStore.user.id,
      total_points: 0,
      rank: 1,
    });

    createDialog.value = false;
    Notify.create({ type: 'positive', message: '¡Grupo creado exitosamente!' });
    void fetchMyGroups();
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Por favor intente nuevamente';
    Notify.create({ type: 'negative', message: message });
  } finally {
    creating.value = false;
  }
}

function goToGroup(id: string) {
  void router.push(`/app/groups/${id}`);
}

onMounted(() => {
  void fetchTournamentSettings();
  void fetchMyGroups();
});
</script>

<style scoped>
.rounded-borders {
  border-radius: 8px;
}
.hover-shadow {
  transition: box-shadow 0.2s ease-in-out;
}
.hover-shadow:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}
</style>
