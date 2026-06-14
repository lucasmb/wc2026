<template>
  <q-page class="q-pa-md">
    <div class="text-h5 text-weight-bold text-primary q-mb-md">
      Resultados y Clasificación
    </div>

    <q-select
      v-model="selectedGroupId"
      :options="groupOptions"
      option-value="id"
      option-label="name"
      emit-value
      map-options
      outlined
      dense
      label="Seleccionar Grupo"
      class="q-mb-md"
      style="max-width: 300px"
    />

    <div v-if="loading" class="row justify-center q-my-xl">
      <q-spinner-dots color="primary" size="40px" />
    </div>

    <div v-else>
      <!-- Leaderboard Section -->
      <q-card
        flat
        bordered
        :dark="$q.dark.isActive"
        :class="$q.dark.isActive ? 'bg-grey-10 border-grey-8' : 'bg-white'"
        class="rounded-borders shadow-1 q-mb-lg"
      >
        <q-card-section class="q-pb-none">
          <div class="text-h6 text-weight-bold text-primary text-center">🏆 Clasificación</div>
        </q-card-section>

        <q-card-section class="q-pt-md">
          <q-list separator>
            <q-item
              v-for="user in leaderboard"
              :key="user.userId"
              class="q-py-md q-my-sm rounded-borders list-row-animation cursor-pointer"
              :class="getRankingClass(user.rank)"
              @click="toggleUserPredictions(user.userId)"
            >
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
                    :class="getTextClass(user.rank)"
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
                      :class="getTextClass(user.rank)"
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

              <q-item-section side>
                <div class="row items-center q-gutter-x-sm">
                  <q-chip
                    dense
                    :color="user.totalPoints > 0 && user.rank <= 3 ? 'black' : 'primary'"
                    :text-color="user.totalPoints > 0 && user.rank <= 3 ? 'white' : 'white'"
                    class="text-weight-bolder text-subtitle2 q-px-md shadow-2"
                  >
                    {{ user.totalPoints }} Pts
                  </q-chip>
                  <q-icon
                    :name="expandedUserId === user.userId ? 'expand_less' : 'expand_more'"
                    :color="user.totalPoints > 0 && user.rank <= 3 ? 'black' : 'grey-7'"
                  />
                </div>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>

      <!-- User Predictions Expansion -->
      <q-slide-transition v-if="expandedUserId">
        <div v-show="expandedUserId" class="q-mb-lg">
          <q-card
            flat
            bordered
            :class="$q.dark.isActive ? 'bg-grey-10' : 'bg-white'"
            class="rounded-borders"
          >
            <q-card-section>
              <div class="text-h6 text-weight-bold text-primary q-mb-md">
                Predicciones de {{ getExpandedUserName() }}
              </div>

              <div v-if="finishedMatches.length === 0" class="text-center q-my-md text-grey-6">
                No hay partidos finalizados aún.
              </div>

              <q-list v-else separator>
                <q-item
                  v-for="match in finishedMatches"
                  :key="match.id"
                  class="q-py-sm"
                >
                  <q-item-section>
                    <q-item-label class="text-weight-medium">
                      {{ match.expand?.home_team?.name || 'Local' }} vs {{ match.expand?.away_team?.name || 'Visitante' }}
                    </q-item-label>
                    <q-item-label caption>
                      {{ getMatchCaption(match) }} • Resultado: {{ match.score_home }} - {{ match.score_away }}
                    </q-item-label>
                  </q-item-section>

                  <q-item-section side>
                    <div class="row items-center q-gutter-x-md">
                      <div class="text-center">
                        <div class="text-caption text-grey-6">Predicción</div>
                        <div class="text-weight-bold">
                          {{ getUserPredictionForMatch(expandedUserId, match.id)?.predicted_home ?? '-' }} -
                          {{ getUserPredictionForMatch(expandedUserId, match.id)?.predicted_away ?? '-' }}
                        </div>
                      </div>
                      <q-badge
                        :color="getPointsColor(getUserPredictionForMatch(expandedUserId, match.id)?.points_awarded ?? 0)"
                        class="q-px-md text-weight-bold"
                      >
                        {{ getUserPredictionForMatch(expandedUserId, match.id)?.points_awarded ?? 0 }} pts
                      </q-badge>
                    </div>
                  </q-item-section>
                </q-item>
              </q-list>
            </q-card-section>
          </q-card>
        </div>
      </q-slide-transition>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { pb, PB_URL } from '@/boot/pocketbase';
import { useAuthStore } from '@/stores/auth';
import { useTournamentStore } from '@/stores/tournament';
import type { Match, PredictionGroup, Prediction, LeaderboardUser } from '@/types';

interface UserPrediction extends Prediction {
  username: string;
  avatarUrl: string;
}

const authStore = useAuthStore();
const tournamentStore = useTournamentStore();

const loading = ref(true);
const selectedGroupId = ref<string>('');
const myGroups = ref<PredictionGroup[]>([]);
const leaderboard = ref<LeaderboardUser[]>([]);
const allPredictions = ref<UserPrediction[]>([]);
const expandedUserId = ref<string | null>(null);

const groupOptions = computed(() => myGroups.value);

const finishedMatches = computed(() => {
  return tournamentStore.matches
    .filter((m) => m.status === 'finished')
    .sort((a, b) => new Date(a.kickoff).getTime() - new Date(b.kickoff).getTime());
});

function getMatchCaption(match: Match): string {
  const phase = match.phase === 'group' ? `Grupo ${match.group_code}` : match.phase.toUpperCase();
  return `${phase} • Partido #${match.match_number}`;
}

function getUserPredictionForMatch(userId: string, matchId: string): UserPrediction | undefined {
  return allPredictions.value.find((p) => p.user === userId && p.match === matchId);
}

function getPointsColor(points: number): string {
  if (points >= 8) return 'positive';
  if (points >= 3) return 'warning';
  if (points > 0) return 'info';
  return 'grey';
}

function getRankingClass(rank: number): string {
  if (rank === 1) return 'rank-1';
  if (rank === 2) return 'rank-2';
  if (rank === 3) return 'rank-3';
  return $q_dark() ? 'bg-grey-9' : 'bg-grey-1';
}

function getTextClass(rank: number): string {
  if (rank <= 3) return 'text-black';
  return $q_dark() ? 'text-grey-4' : 'text-grey-8';
}

function getMedalColor(rank: number): string {
  if (rank === 1) return 'warning';
  if (rank === 2) return 'grey-6';
  if (rank === 3) return 'orange-9';
  return 'grey-7';
}

function $q_dark(): boolean {
  return document.body.classList.contains('body--dark');
}

function toggleUserPredictions(userId: string) {
  expandedUserId.value = expandedUserId.value === userId ? null : userId;
}

function getExpandedUserName(): string {
  const user = leaderboard.value.find((u) => u.userId === expandedUserId.value);
  return user?.username || 'Usuario';
}

async function fetchUserGroups() {
  if (!authStore.user?.id) return;
  try {
    const memberships = await pb.collection('group_members_id').getFullList({
      filter: `user = "${authStore.user.id}"`,
      expand: 'prediction_group',
    });

    myGroups.value = memberships
      .map((m) => m.expand?.prediction_group as unknown as PredictionGroup)
      .filter(Boolean);

    if (myGroups.value.length > 0 && myGroups.value[0]) {
      selectedGroupId.value = String(myGroups.value[0].id);
    }
  } catch (err: unknown) {
    console.error('Failed to load user prediction groups:', err);
  }
}

async function fetchLeaderboard() {
  if (!selectedGroupId.value) return;
  try {
    const response = await fetch(`${PB_URL}/api/wc/leaderboard/${selectedGroupId.value}`, {
      headers: {
        Authorization: `Bearer ${pb.authStore.token}`,
      },
    });
    if (response.ok) {
      leaderboard.value = await response.json();
    }
  } catch (err: unknown) {
    console.error('Failed fetching group leaderboard:', err);
  }
}

async function fetchPredictions() {
  if (!selectedGroupId.value) return;
  try {
    const members = await pb.collection('group_members_id').getFullList({
      filter: `prediction_group = "${selectedGroupId.value}"`,
      expand: 'user',
    });

    const userMap = new Map<string, { username: string; avatarUrl: string }>();
    
    members.forEach((m) => {
      const user = m.expand?.user as { id: string; username: string; avatar_url?: string } | undefined;
      if (user) {
        userMap.set(user.id, {
          username: user.username || 'Unknown',
          avatarUrl: user.avatar_url || '',
        });
      }
    });

    const predictions = await pb.collection('predictions_id').getFullList({
      filter: `prediction_group = "${selectedGroupId.value}"`,
    });

    allPredictions.value = predictions.map((p) => {
      const userInfo = userMap.get(p.user as string) || { username: 'Unknown', avatarUrl: '' };
      return {
        ...p,
        username: userInfo.username,
        avatarUrl: userInfo.avatarUrl,
      } as unknown as UserPrediction;
    });
  } catch (err: unknown) {
    console.error('Failed to fetch predictions:', err);
  }
}

watch(selectedGroupId, async () => {
  loading.value = true;
  expandedUserId.value = null;
  await Promise.all([fetchLeaderboard(), fetchPredictions()]);
  loading.value = false;
});

onMounted(async () => {
  await tournamentStore.fetchMatches();
  await fetchUserGroups();
  loading.value = false;
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
.cursor-pointer {
  cursor: pointer;
}
.rank-1 {
  background-color: #fbbf24 !important;
}
.rank-2 {
  background-color: #cbd5e1 !important;
}
.rank-3 {
  background-color: #fdba74 !important;
}
</style>
