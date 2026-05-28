<template>
  <q-page class="q-pa-md">
    <!-- Header Block -->
    <div v-if="group" class="row items-center justify-between q-mb-lg">
      <div>
        <div class="text-h6 text-weight-bold text-primary">{{ group.name }}</div>
        <div class="text-caption text-grey-6 q-mt-xs">
          Invite Token:
          <span class="text-weight-bold text-secondary text-subtitle2 q-ml-xs">
            {{ group.invite_code }}
          </span>
        </div>
      </div>
      <q-btn flat icon="arrow_back" label="Back to Groups" color="grey-7" to="/app/groups" />
    </div>

    <!-- Leaderboard Standings Loader -->
    <div v-if="loading" class="row justify-center q-my-xl">
      <q-spinner-dots color="primary" size="40px" />
    </div>

    <div v-else class="row q-col-gutter-md">
      <!-- Leaderboard Column -->
      <div class="col-12">
        <q-card flat bordered class="rounded-borders bg-white">
          <q-card-section class="q-pb-none">
            <div class="text-subtitle1 text-weight-bold text-primary">Standings Leaderboard</div>
          </q-card-section>

          <q-card-section>
            <q-list separator>
              <q-item v-for="user in leaderboard" :key="user.userId" class="q-py-md">
                <q-item-section avatar class="row items-center justify-center">
                  <span class="text-subtitle1 text-weight-bold q-mr-sm text-grey-8">
                    #{{ user.rank }}
                  </span>
                  <q-avatar size="36px" color="blue-1" text-color="primary">
                    {{ user.username.charAt(0).toUpperCase() }}
                  </q-avatar>
                </q-item-section>

                <q-item-section>
                  <q-item-label class="text-weight-bold text-grey-9">
                    {{ user.username }}
                  </q-item-label>
                </q-item-section>

                <q-item-section side>
                  <q-chip dense color="primary" text-color="white" class="text-weight-bold q-px-md">
                    {{ user.totalPoints }} Pts
                  </q-chip>
                </q-item-section>
              </q-item>
            </q-list>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { pb, PB_URL } from 'src/boot/pocketbase';
import type { PredictionGroup, LeaderboardUser } from 'src/types';

const route = useRoute();
const groupId = route.params.id as string;

const group = ref<PredictionGroup | null>(null);
const leaderboard = ref<LeaderboardUser[]>([]);
const loading = ref(true);

async function fetchGroupDetails() {
  try {
    const rawGroup = await pb.collection('prediction_groups_id').getOne(groupId);
    group.value = rawGroup as unknown as PredictionGroup;
  } catch (err: unknown) {
    console.error('Failed fetching group detail metrics:', err);
  }
}

async function fetchLeaderboard() {
  try {
    // Call the custom backend endpoint registered in Step 3
    const response = await fetch(`${PB_URL}/api/wc/leaderboard/${groupId}`, {
      headers: {
        Authorization: `Bearer ${pb.authStore.token}`,
      },
    });
    if (response.ok) {
      leaderboard.value = await response.json();
    }
  } catch (err: unknown) {
    console.error('Failed fetching group leaderboard:', err);
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await fetchGroupDetails();
  await fetchLeaderboard();
});
</script>
