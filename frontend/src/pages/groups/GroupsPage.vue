<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-lg">
      <div class="text-h6 text-weight-bold text-primary">Prediction Groups</div>
      <div class="q-gutter-x-sm">
        <q-btn
          label="Join with Code"
          color="secondary"
          outline
          dense
          class="q-px-md"
          @click="showJoinDialog"
        />
        <q-btn
          label="Create Group"
          color="primary"
          unelevated
          dense
          class="q-px-md"
          @click="showCreateDialog"
        />
      </div>
    </div>

    <!-- Active Groups Grid -->
    <div v-if="loading" class="row justify-center q-my-xl">
      <q-spinner-dots color="primary" size="40px" />
    </div>

    <div v-else class="row q-col-gutter-md">
      <div v-if="groups.length === 0" class="col-12 text-center q-my-xl text-grey-6">
        You are not a member of any prediction group yet. Click above to join or create one!
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
              Owner: {{ group.expand?.owner?.username || 'System' }}
            </div>
          </q-card-section>
          <q-card-actions align="right">
            <q-btn flat label="View Leaderboard" color="primary" dense />
          </q-card-actions>
        </q-card>
      </div>
    </div>

    <!-- Join Group Dialog -->
    <q-dialog v-model="joinDialog" persistent>
      <q-card style="min-width: 320px">
        <q-card-section class="q-pb-none">
          <div class="text-h6 text-weight-bold">Join Group</div>
        </q-card-section>
        <q-card-section class="q-py-md">
          <q-input
            v-model="inviteCode"
            outlined
            label="Invite Code"
            placeholder="Enter 8-digit token"
            dense
          />
        </q-card-section>
        <q-card-actions align="right" class="q-pt-none q-pb-md q-px-md">
          <q-btn flat label="Cancel" color="grey-6" v-close-popup />
          <q-btn label="Join" color="primary" unelevated :loading="joining" @click="joinGroup" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- Create Group Dialog -->
    <q-dialog v-model="createDialog" persistent>
      <q-card style="min-width: 320px">
        <q-card-section class="q-pb-none">
          <div class="text-h6 text-weight-bold">Create Group</div>
        </q-card-section>
        <q-card-section class="q-py-md">
          <q-input
            v-model="newGroupName"
            outlined
            label="Group Name"
            placeholder="e.g. Work Friends"
            dense
          />
        </q-card-section>
        <q-card-actions align="right" class="q-pt-none q-pb-md q-px-md">
          <q-btn flat label="Cancel" color="grey-6" v-close-popup />
          <q-btn
            label="Create"
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
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useQuasar } from 'quasar';
import { pb } from 'src/boot/pocketbase';
import { useAuthStore } from 'src/stores/auth';
import type { PredictionGroup } from 'src/types';

const router = useRouter();
const $q = useQuasar();
const authStore = useAuthStore();

const groups = ref<PredictionGroup[]>([]);
const loading = ref(true);

const joinDialog = ref(false);
const inviteCode = ref('');
const joining = ref(false);

const createDialog = ref(false);
const newGroupName = ref('');
const creating = ref(false);

async function fetchMyGroups() {
  if (!authStore.user?.id) return;
  try {
    // 1. Get all memberships for the logged in user
    const memberships = await pb.collection('group_members_id').getFullList({
      filter: `user = "${authStore.user.id}"`,
      expand: 'prediction_group.owner',
    });

    // 2. Map and extract group records
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
    // Resolve Group by invite code
    const groupList = await pb.collection('prediction_groups_id').getList(1, 1, {
      filter: `invite_code = "${inviteCode.value}"`,
    });

    const targetGroup = groupList.items[0];
    // TypeScript safe guard: completely narrows the type
    if (!targetGroup) {
      throw new Error('Invalid or expired invite code.');
    }

    // Create membership record
    await pb.collection('group_members_id').create({
      prediction_group: targetGroup.id,
      user: authStore.user.id,
      total_points: 0,
      rank: 1,
    });

    joinDialog.value = false;
    $q.notify({ type: 'positive', message: 'Successfully joined group!' });

    // Prefix with void to satisfy @typescript-eslint/no-floating-promises
    void fetchMyGroups();
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Please try again';
    $q.notify({ type: 'negative', message: message });
  } finally {
    joining.value = false;
  }
}

// Generate random alphanumeric 8-character token
function generateInviteToken(): string {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
  let result = '';
  for (let i = 0; i < 8; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

async function createGroup() {
  if (!newGroupName.value || !authStore.user?.id) return;
  creating.value = true;
  try {
    const inviteToken = generateInviteToken();

    // Create prediction group
    const createdGroup = await pb.collection('prediction_groups_id').create({
      name: newGroupName.value,
      owner: authStore.user.id,
      invite_code: inviteToken,
      is_public: false,
    });

    // Automatically register creator as a member
    await pb.collection('group_members_id').create({
      prediction_group: createdGroup.id,
      user: authStore.user.id,
      total_points: 0,
      rank: 1,
    });

    createDialog.value = false;
    $q.notify({ type: 'positive', message: 'Group successfully created!' });

    // Prefix with void to satisfy floating promise rules
    void fetchMyGroups();
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Please try again';
    $q.notify({ type: 'negative', message: message });
  } finally {
    creating.value = false;
  }
}

function goToGroup(id: string) {
  // Prefix router push with void to satisfy floating promise rules
  void router.push(`/app/groups/${id}`);
}

onMounted(() => {
  void fetchMyGroups();
});
</script>

<style scoped>
.hover-shadow {
  transition: box-shadow 0.2s ease-in-out;
}
.hover-shadow:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}
</style>
