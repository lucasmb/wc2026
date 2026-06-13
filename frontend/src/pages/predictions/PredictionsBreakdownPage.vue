<template>
  <q-page class="q-pa-md">
    <div class="text-h5 text-weight-bold text-primary q-mb-md">
      Predicciones de Partidos Finalizados
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

    <div v-else-if="finishedMatches.length === 0" class="text-center q-my-xl text-grey-6">
      No hay partidos finalizados aún.
    </div>

    <q-list v-else separator>
      <q-expansion-item
        v-for="match in finishedMatches"
        :key="match.id"
        :label="getMatchLabel(match)"
        :caption="getMatchCaption(match)"
        :header-class="$q.dark.isActive ? 'bg-grey-9' : 'bg-grey-2'"
        class="q-mb-sm rounded-borders"
      >
        <q-card>
          <q-card-section>
            <div class="text-h6 text-center q-mb-md">
              <span class="text-weight-bold">{{ match.expand?.home_team?.name || 'Local' }}</span>
              <span class="q-mx-md text-primary">{{ match.score_home }} - {{ match.score_away }}</span>
              <span class="text-weight-bold">{{ match.expand?.away_team?.name || 'Visitante' }}</span>
            </div>

            <q-table
              :rows="getPredictionsForMatch(match.id)"
              :columns="predictionColumns"
              flat
              bordered
              dense
              hide-bottom
              :class="$q.dark.isActive ? 'bg-grey-10' : 'bg-white'"
            >
              <template #body-cell-user="props">
                <q-td :props="props">
                  <div class="row items-center q-gutter-x-sm">
                    <q-avatar size="32px" color="primary" text-color="white">
                      <img v-if="props.row.avatarUrl" :src="props.row.avatarUrl" :alt="props.row.username" />
                      <span v-else>{{ props.row.username.charAt(0).toUpperCase() }}</span>
                    </q-avatar>
                    <span class="text-weight-medium">{{ props.row.username }}</span>
                  </div>
                </q-td>
              </template>
              <template #body-cell-prediction="props">
                <q-td :props="props" class="text-center">
                  <span class="text-weight-bold">{{ props.row.predicted_home }} - {{ props.row.predicted_away }}</span>
                </q-td>
              </template>
              <template #body-cell-points="props">
                <q-td :props="props" class="text-center">
                  <q-badge :color="getPointsColor(props.row.points_awarded)" class="q-px-md">
                    {{ props.row.points_awarded }} pts
                  </q-badge>
                </q-td>
              </template>
            </q-table>
          </q-card-section>
        </q-card>
      </q-expansion-item>
    </q-list>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { pb } from '@/boot/pocketbase';
import { useAuthStore } from '@/stores/auth';
import { useTournamentStore } from '@/stores/tournament';
import type { Match, PredictionGroup, Prediction } from '@/types';

interface UserPrediction extends Prediction {
  username: string;
  avatarUrl: string;
}

const authStore = useAuthStore();
const tournamentStore = useTournamentStore();

const loading = ref(true);
const selectedGroupId = ref<string>('');
const myGroups = ref<PredictionGroup[]>([]);
const allPredictions = ref<UserPrediction[]>([]);

const groupOptions = computed(() => myGroups.value);

const finishedMatches = computed(() => {
  return tournamentStore.matches
    .filter((m) => m.status === 'finished')
    .sort((a, b) => new Date(a.kickoff).getTime() - new Date(b.kickoff).getTime());
});

const predictionColumns = [
  { name: 'user', label: 'Usuario', field: 'username', align: 'left' as const },
  { name: 'prediction', label: 'Predicción', field: 'prediction', align: 'center' as const },
  { name: 'points', label: 'Puntos', field: 'points_awarded', align: 'center' as const, sortable: true },
];

function getMatchLabel(match: Match): string {
  const home = match.expand?.home_team?.name || 'Local';
  const away = match.expand?.away_team?.name || 'Visitante';
  return `${home} ${match.score_home} - ${match.score_away} ${away}`;
}

function getMatchCaption(match: Match): string {
  const phase = match.phase === 'group' ? `Grupo ${match.group_code}` : match.phase.toUpperCase();
  return `${phase} • Partido #${match.match_number}`;
}

function getPredictionsForMatch(matchId: string): UserPrediction[] {
  return allPredictions.value
    .filter((p) => p.match === matchId)
    .sort((a, b) => b.points_awarded - a.points_awarded);
}

function getPointsColor(points: number): string {
  if (points >= 8) return 'positive';
  if (points >= 3) return 'warning';
  if (points > 0) return 'info';
  return 'grey';
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

async function fetchPredictions() {
  if (!selectedGroupId.value) return;
  loading.value = true;
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
  } finally {
    loading.value = false;
  }
}

watch(selectedGroupId, () => {
  void fetchPredictions();
});

onMounted(async () => {
  await tournamentStore.fetchMatches();
  await fetchUserGroups();
});
</script>

<style scoped>
.rounded-borders {
  border-radius: 8px;
}
</style>
