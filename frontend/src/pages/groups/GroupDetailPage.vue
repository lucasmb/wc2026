<template>
  <q-page class="q-pa-md">
    <!-- Header Block -->
    <div v-if="group" class="row items-center justify-between q-mb-lg">
      <div>
        <div class="text-h5 text-weight-bold text-primary">{{ group.name }}</div>
        <div class="text-caption text-grey-6 q-mt-xs row items-center q-gutter-x-sm">
          <span>Código de Grupo:</span>
          <span class="text-weight-bolder text-secondary text-subtitle2">
            {{ group.invite_code }}
          </span>
          <q-btn
            flat
            dense
            size="md"
            color="primary"
            icon="share"
            label="Invitar amigos"
            class="q-ml-sm"
            @click="copyInviteLink"
          />
        </div>
      </div>
      <q-btn flat icon="arrow_back" label="Volver a Grupos" color="grey-7" to="/app/groups" />
    </div>

    <!-- Leaderboard Loader -->
    <div v-if="loading" class="row justify-center q-my-xl">
      <q-spinner-dots color="primary" size="40px" />
    </div>

    <div v-else class="row q-col-gutter-md justify-center">
      <!-- Leaderboard Column -->
      <div class="col-12 col-md-8">
        <q-card
          flat
          bordered
          :dark="$q.dark.isActive"
          :class="$q.dark.isActive ? 'bg-grey-10 border-grey-8' : 'bg-white'"
          class="rounded-borders shadow-1"
        >
          <q-card-section class="q-pb-none">
            <div class="text-h6 text-weight-bold text-primary text-center">🏆 Leaderboard</div>
          </q-card-section>

          <q-card-section class="q-pt-md">
            <q-list separator>
              <q-item
                v-for="user in leaderboard"
                :key="user.userId"
                class="q-py-md q-my-sm rounded-borders list-row-animation"
                :class="
                  user.totalPoints > 0
                    ? getRankingClass(user.rank)
                    : $q.dark.isActive
                      ? 'bg-grey-9 border'
                      : 'bg-grey-1 border'
                "
              >
                <!-- Rank Medal Icon / Badge Slot -->
                <q-item-section avatar class="row items-center justify-center">
                  <div class="row items-center q-gutter-x-sm">
                    <q-icon
                      v-if="user.totalPoints > 0 && user.rank <= 3"
                      name="emoji_events"
                      size="28px"
                      :color="getMedalColor(user.rank)"
                    />
                    <span
                      class="text-subtitle1 text-weight-bolder"
                      :class="
                        user.totalPoints > 0 && user.rank <= 3
                          ? 'text-black'
                          : $q.dark.isActive
                            ? 'text-grey-4'
                            : 'text-grey-8'
                      "
                    >
                      #{{ user.rank }}
                    </span>
                  </div>
                </q-item-section>

              <q-item-section>
                <div class="row items-center q-gutter-x-md">
                  <q-avatar
                    size="40px"
                    :color="user.totalPoints > 0 && user.rank <= 3 ? 'white' : 'primary'"
                    :text-color="user.totalPoints > 0 && user.rank <= 3 ? 'black' : 'white'"
                    class="shadow-1"
                  >
                    <img v-if="user.avatarUrl || user.avatar" :src="user.avatarUrl || user.avatar" :alt="user.username" />
                    <span v-else>{{ user.username.charAt(0).toUpperCase() }}</span>
                  </q-avatar>
                    <div>
                      <q-item-label
                        class="text-subtitle2 text-weight-bold"
                        :class="
                          user.totalPoints > 0 && user.rank <= 3
                            ? 'text-black'
                            : $q.dark.isActive
                              ? 'text-white'
                              : 'text-grey-9'
                        "
                      >
                        {{ user.username }}
                      </q-item-label>
                      <q-item-label
                        caption
                        v-if="user.totalPoints > 0 && user.rank === 1"
                        class="text-weight-bold text-amber-10"
                      >
                        Líder
                      </q-item-label>
                    </div>
                  </div>
                </q-item-section>

                <!-- Points Highlight Badge -->
                <q-item-section side>
                  <q-chip
                    dense
                    :color="user.totalPoints > 0 && user.rank <= 3 ? 'black' : 'primary'"
                    :text-color="user.totalPoints > 0 && user.rank <= 3 ? 'white' : 'white'"
                    class="text-weight-bolder text-subtitle2 q-px-md shadow-2"
                  >
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
import { pb, PB_URL } from '@/boot/pocketbase';
import { Notify, copyToClipboard } from 'quasar'; // Import copyToClipboard
import type { PredictionGroup, LeaderboardUser } from '@/types';

const route = useRoute();
const groupId = route.params.id as string;

const group = ref<PredictionGroup | null>(null);
const leaderboard = ref<LeaderboardUser[]>([]);
const loading = ref(true);

async function copyInviteLink() {
  // Generates cross-platform hash-based URL
  const inviteUrl = `${window.location.origin}/#/?invite=${groupId}`;
  try {
    await copyToClipboard(inviteUrl);
    Notify.create({
      type: 'positive',
      message: '¡Enlace de invitación copiado al portapapeles!',
    });
  } catch (err: unknown) {
    console.error('Failed copying invite link:', err);
  }
}

async function fetchGroupDetails() {
  try {
    const rawGroup = await pb.collection('prediction_groups_id').getOne(groupId);
    group.value = rawGroup as unknown as PredictionGroup;
  } catch (err: unknown) {
    console.error('Failed fetching group details:', err);
  }
}

async function fetchLeaderboard() {
  try {
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

// Visual stylings based on ranking placements (Podium layout)
function getRankingClass(rank: number): string {
  switch (rank) {
    case 1:
      return 'bg-amber-4 shadow-1'; // Gold
    case 2:
      return 'bg-blue-grey-2 shadow-1'; // Silver
    case 3:
      return 'bg-orange-3 shadow-1'; // Bronze
    default:
      return 'bg-transparent';
  }
}

function getMedalColor(rank: number): string {
  switch (rank) {
    case 1:
      return 'warning';
    case 2:
      return 'grey-6';
    case 3:
      return 'orange-9';
    default:
      return 'grey-7';
  }
}

onMounted(async () => {
  await fetchGroupDetails();
  await fetchLeaderboard();
});
</script>

<style scoped>
.rounded-borders {
  border-radius: 12px;
}
.border-grey-8 {
  border-color: rgba(255, 255, 255, 0.08) !important;
}
.list-row-animation {
  transition:
    transform 0.2s ease-in-out,
    box-shadow 0.2s ease-in-out;
}
.list-row-animation:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}
</style>
