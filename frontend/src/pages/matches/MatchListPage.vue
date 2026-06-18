<template>
  <q-page class="q-pa-md q-pb-xl">
    <!-- Header Block -->
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h6 text-weight-bold text-primary">Predicciones de Partidos</div>

      <!-- Selectors & Slider Row -->
      <div class="row items-center q-gutter-x-sm">
        <!-- Slider Navigation Controls -->
        <div
          class="row items-center border rounded-borders q-px-xs"
          :class="$q.dark.isActive ? 'bg-grey-10' : 'bg-white'"
        >
          <q-btn
            flat
            dense
            round
            icon="chevron_left"
            color="primary"
            :disabled="isFirstGroup"
            @click="prevGroup"
          />
          <q-select
            v-model="selectedGroupStage"
            :options="groupStageOptions"
            :dark="$q.dark.isActive"
            borderless
            dense
            emit-value
            map-options
            input-class="text-weight-bold text-center"
            style="width: 100px"
          />
          <q-btn
            flat
            dense
            round
            icon="chevron_right"
            color="primary"
            :disabled="isLastGroup"
            @click="nextGroup"
          />
        </div>

        <!-- User's Prediction Groups Dropdown -->
        <q-select
          v-if="myGroups.length > 0"
          v-model="activeGroupId"
          :options="myGroups"
          :dark="$q.dark.isActive"
          option-value="id"
          option-label="name"
          emit-value
          map-options
          label="Predicting For"
          outlined
          dense
          style="width: 170px"
        />
      </div>
    </div>

    <!-- No Groups Warning -->
    <div v-if="!loadingGroups && myGroups.length === 0" class="text-center q-my-xl">
      <q-card flat bordered class="q-pa-xl rounded-borders max-width-card mx-auto">
        <q-icon name="group" size="64px" color="grey-4" />
        <div class="text-h6 text-weight-bold text-primary q-mt-md">
          Unite a un grupo para empezar
        </div>
        <div class="text-subtitle2 text-grey-6 q-mt-xs">
          Las predicciones se guardan dentro de un grupo. Crea o unete a un grupo para empezar!
        </div>
        <q-btn
          label="Go to Groups"
          color="primary"
          unelevated
          to="/app/groups"
          class="q-mt-lg q-px-xl"
        />
      </q-card>
    </div>

    <!-- Loading Spinner -->
    <div
      v-else-if="loadingMatches || loadingGroups || loadingPredictions"
      class="row justify-center q-my-xl"
    >
      <q-spinner-dots color="primary" size="40px" />
    </div>

    <!-- Render Match Cards Grid -->
    <div v-else class="row q-col-gutter-md q-mb-xl">
      <div v-if="filteredMatches.length === 0" class="col-12 text-center q-my-xl text-grey-6">
        No hay partidos configurados para esta seleccion
      </div>
      <div v-for="match in filteredMatches" :key="match.id" class="col-12 col-md-6">
        <div class="text-caption text-grey-6 q-mb-xs">
          Partido #{{ match.match_number }} • {{ formatLocalTime(match.kickoff) }} •
          {{ match.venue }}
        </div>
        <!-- Child Form emitting score changes back to parent -->
        <PredictionForm
          v-if="activeGroupId"
          :match="match"
          :groupId="activeGroupId"
          :prediction="predictionsMap[match.id]"
          @update="handlePredictionUpdate(match.id, $event)"
        />
        <!-- Finished game score display -->
        <div
          v-if="match.status === 'finished'"
          class="q-mt-xs q-pa-sm rounded-borders text-center"
          :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-grey-2'"
        >
          <div class="text-caption text-weight-bold q-mb-xs">Resultado Final</div>
          <div class="row items-center justify-center q-gutter-x-md">
            <span class="text-weight-bold">{{ match.expand?.home_team?.name }}</span>
            <span class="text-h6 text-weight-bolder text-primary">
              {{ match.score_home }} - {{ match.score_away }}
            </span>
            <span class="text-weight-bold">{{ match.expand?.away_team?.name }}</span>
          </div>
          <q-badge color="grey-7" class="q-mt-xs">Finalizado</q-badge>
        </div>
        <!-- Collapsible predictions for locked/live/finished matches -->
        <q-card
          v-if="activeGroupId && isMatchLockedOrActive(match) && matchPredictions[match.id]?.length"
          flat
          bordered
          class="q-mt-xs rounded-borders"
          :class="$q.dark.isActive ? 'bg-grey-10' : 'bg-grey-1'"
        >
          <q-expansion-item
            dense
            :label="`Predicciones del grupo (${matchPredictions[match.id]?.length ?? 0})`"
            :header-class="$q.dark.isActive ? 'text-grey-4' : 'text-grey-8'"
          >
            <q-separator />
            <q-list dense>
              <q-item
                v-for="pred in matchPredictions[match.id]"
                :key="pred.id"
                class="q-py-xs"
              >
                <q-item-section avatar>
                  <q-avatar size="28px" color="primary" text-color="white">
                    <img v-if="pred.userAvatar" :src="pred.userAvatar" :alt="pred.userName" />
                    <span v-else class="text-caption">{{ pred.userName.charAt(0).toUpperCase() }}</span>
                  </q-avatar>
                </q-item-section>
                <q-item-section>
                  <q-item-label class="text-caption text-weight-medium">{{ pred.userName }}</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <div class="row items-center q-gutter-x-sm">
                    <span class="text-caption text-weight-bold">
                      {{ pred.predicted_home }} - {{ pred.predicted_away }}
                    </span>
                    <q-badge
                      v-if="match.status === 'finished'"
                      :color="pred.points_awarded > 0 ? 'positive' : 'grey'"
                      dense
                    >
                      {{ pred.points_awarded }} pts
                    </q-badge>
                  </div>
                </q-item-section>
              </q-item>
            </q-list>
          </q-expansion-item>
        </q-card>
      </div>

      <!-- Unified Action Block at the bottom of the scroll flow (Mobile-safe, non-overlapping) -->
      <div
        v-if="myGroups.length > 0"
        class="col-12 q-py-md q-px-md border-top q-mt-lg rounded-borders shadow-1"
        :class="$q.dark.isActive ? 'bg-grey-10 text-white border-grey-8' : 'bg-white text-grey-7'"
      >
        <div class="max-width-footer mx-auto row justify-between items-center q-col-gutter-sm">
          <span
            class="text-caption text-weight-medium col-12 col-sm-auto text-center text-sm-left"
            :class="$q.dark.isActive ? 'text-grey-4' : 'text-grey-8'"
          >
            Modifica tus predicciones y guarda a continuación.
          </span>
          <div class="row q-gutter-x-sm col-12 col-sm-auto justify-center q-mt-sm q-mt-sm-none">
            <q-btn
              flat
              label="Siguiente Grupo"
              color="secondary"
              icon-right="chevron_right"
              :disabled="isLastGroup"
              @click="nextGroup"
            />
            <q-btn
              label="Guardar Grupo"
              color="primary"
              icon="save"
              unelevated
              :loading="savingAll"
              :disabled="isSaveDisabled"
              @click="saveAllGroupPredictions"
            />
          </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { pb, PB_URL, getFileUrl } from '@/boot/pocketbase';
import { Notify } from 'quasar';

import { useAuthStore } from '@/stores/auth';
import { useTournamentStore } from '@/stores/tournament';
import PredictionForm from '@/components/match/PredictionForm.vue';
import type { PredictionGroup, Prediction, Match } from '@/types';

interface MatchPrediction {
  id: string;
  userId: string;
  predicted_home: number;
  predicted_away: number;
  points_awarded: number;
  userName: string;
  userAvatar: string;
}

const tournamentStore = useTournamentStore();
const authStore = useAuthStore();

const loadingMatches = ref(true);
const loadingGroups = ref(true);
const loadingPredictions = ref(false);
const savingAll = ref(false);

const selectedGroupStage = ref('A');
const activeGroupId = ref<string>('');
const myGroups = ref<PredictionGroup[]>([]);

const predictionsMap = ref<Record<string, Prediction>>({});
const matchPredictions = ref<Record<string, MatchPrediction[]>>({});
const currentScores = ref<Record<string, { predicted_home: number; predicted_away: number }>>({});
const dirtyMatches = ref<Record<string, boolean>>({});

// Expanded to include both Groups A-L and the Knockout Single-Elimination phases
const groupStageOptions = [
  { label: 'Group A', value: 'A' },
  { label: 'Group B', value: 'B' },
  { label: 'Group C', value: 'C' },
  { label: 'Group D', value: 'D' },
  { label: 'Group E', value: 'E' },
  { label: 'Group F', value: 'F' },
  { label: 'Group G', value: 'G' },
  { label: 'Group H', value: 'H' },
  { label: 'Group I', value: 'I' },
  { label: 'Group J', value: 'J' },
  { label: 'Group K', value: 'K' },
  { label: 'Group L', value: 'L' },
  // KNOCKOUT PHASES
  { label: 'Round of 32', value: 'r32' },
  { label: 'Round of 16', value: 'r16' },
  { label: 'Quarter-Finals', value: 'qf' },
  { label: 'Semi-Finals', value: 'sf' },
  { label: 'Finals', value: 'final' },
];

const groupStages = [
  'A',
  'B',
  'C',
  'D',
  'E',
  'F',
  'G',
  'H',
  'I',
  'J',
  'K',
  'L',
  'r32',
  'r16',
  'qf',
  'sf',
  'final',
];

const isFirstGroup = computed(() => selectedGroupStage.value === 'A');
const isLastGroup = computed(() => selectedGroupStage.value === 'final');

// Automatically disables the group save action if all matches are locked or finished
const isSaveDisabled = computed(() => {
  if (filteredMatches.value.length === 0) return true;
  return filteredMatches.value.every((m) => {
    if (m.status === 'finished') return true;
    if (!m.kickoff) return false;
    // Standardize parsing for absolute browser engine compatibility
    const isoString = m.kickoff.replace(' ', 'T');
    return new Date(isoString).getTime() <= Date.now();
  });
});
const filteredMatches = computed(() => {
  const selected = selectedGroupStage.value;
  // If a knockout stage is selected, match by the 'phase' attribute directly
  if (['r32', 'r16', 'qf', 'sf', 'final'].includes(selected)) {
    return tournamentStore.matches.filter((m) => m.phase === selected);
  }
  // Otherwise, match group stage letters
  return tournamentStore.matches.filter((m) => m.phase === 'group' && m.group_code === selected);
});

function formatLocalTime(isoString: string) {
  return new Intl.DateTimeFormat(navigator.language, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(isoString));
}

function isMatchLockedOrActive(match: Match): boolean {
  if (match.status === 'finished' || match.status === 'live') return true;
  if (!match.kickoff) return false;
  const isoString = match.kickoff.replace(' ', 'T');
  return new Date(isoString).getTime() <= Date.now();
}

function prevGroup() {
  const idx = groupStages.indexOf(selectedGroupStage.value);
  if (idx > 0) {
    const prevGroupVal = groupStages[idx - 1];
    if (prevGroupVal) {
      selectedGroupStage.value = prevGroupVal;
    }
  }
}

function nextGroup() {
  const idx = groupStages.indexOf(selectedGroupStage.value);
  if (idx < groupStages.length - 1) {
    const nextGroupVal = groupStages[idx + 1];
    if (nextGroupVal) {
      selectedGroupStage.value = nextGroupVal;
    }
  }
}

function handlePredictionUpdate(
  matchId: string,
  formValues: { predicted_home: number; predicted_away: number },
) {
  currentScores.value[matchId] = formValues;

  const existing = predictionsMap.value[matchId];
  const isDifferent =
    !existing ||
    existing.predicted_home !== formValues.predicted_home ||
    existing.predicted_away !== formValues.predicted_away;

  if (isDifferent) {
    dirtyMatches.value[matchId] = true;
  } else {
    delete dirtyMatches.value[matchId];
  }
}

async function fetchPredictions() {
  if (!authStore.user?.id || !activeGroupId.value) return;
  loadingPredictions.value = true;
  try {
    const list = await pb.collection('predictions_id').getFullList({
      filter: `user = "${authStore.user.id}" && prediction_group = "${activeGroupId.value}"`,
    });

    const map: Record<string, Prediction> = {};
    list.forEach((p) => {
      let mId = '';
      if (typeof p.match === 'string') {
        mId = p.match;
      } else if (Array.isArray(p.match) && p.match.length > 0) {
        mId = String(p.match[0]);
      } else if (p.match && typeof p.match === 'object') {
        const matchObj = p.match as { id?: string };
        mId = String(matchObj.id || '');
      }

      if (mId) {
        map[mId] = p as unknown as Prediction;
      }
    });
    predictionsMap.value = map;
    currentScores.value = {};
    dirtyMatches.value = {};

    await fetchGroupPredictions();
  } catch (err: unknown) {
    console.error('Failed fetching predictions batch:', err);
  } finally {
    loadingPredictions.value = false;
  }
}

async function fetchGroupPredictions() {
  if (!activeGroupId.value) return;
  try {
    const response = await fetch(`${PB_URL}/api/wc/match-predictions/${activeGroupId.value}`, {
      headers: {
        Authorization: `Bearer ${pb.authStore.token}`,
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }

    const allPredictions = await response.json();

    const grouped: Record<string, MatchPrediction[]> = {};
    allPredictions.forEach((p: { match: string; id: string; userId: string; predicted_home: number; predicted_away: number; points_awarded: number; userName: string; userAvatarUrl: string; userAvatar: string }) => {
      const matchId = p.match;
      if (matchId) {
        if (!grouped[matchId]) {
          grouped[matchId] = [];
        }
        grouped[matchId].push({
          id: p.id,
          userId: p.userId,
          predicted_home: p.predicted_home,
          predicted_away: p.predicted_away,
          points_awarded: p.points_awarded || 0,
          userName: p.userName,
          userAvatar: p.userAvatarUrl || getFileUrl(p.userId, p.userAvatar) || '',
        });
      }
    });

    Object.values(grouped).forEach((preds) => {
      preds.sort((a, b) => b.points_awarded - a.points_awarded);
    });

    matchPredictions.value = grouped;
  } catch (err: unknown) {
    console.error('Failed fetching group predictions:', err);
  }
}

async function saveAllGroupPredictions() {
  if (!authStore.user?.id || !activeGroupId.value) return;

  const userId = authStore.user.id;
  savingAll.value = true;
  try {
    for (const match of filteredMatches.value) {
      if (!dirtyMatches.value[match.id]) continue;

      const scoreForm = currentScores.value[match.id];
      if (!scoreForm) continue;

      const isLocked =
        match.status === 'finished' ||
        (match.kickoff && new Date(match.kickoff.replace(' ', 'T')).getTime() <= Date.now());
      if (isLocked) continue;

      const existing = predictionsMap.value[match.id];
      const data = {
        user: userId,
        match: match.id,
        prediction_group: activeGroupId.value,
        predicted_home: Math.max(0, scoreForm.predicted_home),
        predicted_away: Math.max(0, scoreForm.predicted_away),
      };

      if (existing?.id) {
        await pb.collection('predictions_id').update(existing.id, data);
      } else {
        await pb.collection('predictions_id').create(data);
      }
    }

    Notify.create({
      type: 'positive',
      message: `Predictions saved successfully for ${selectedGroupStage.value}!`,
    });

    await fetchPredictions();

    if (!isLastGroup.value) {
      nextGroup();
    } else {
      Notify.create({
        type: 'info',
        message: 'All predictions completed!',
      });
    }
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Error executing batch write';
    Notify.create({
      type: 'negative',
      message: `Failed saving predictions: ${message}`,
    });
  } finally {
    savingAll.value = false;
  }
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
      activeGroupId.value = String(myGroups.value[0].id);
    }
  } catch (err: unknown) {
    console.error('Failed to load user prediction groups:', err);
  } finally {
    loadingGroups.value = false;
  }
}

watch(
  [() => authStore.user?.id, () => selectedGroupStage.value],
  async ([userId]) => {
    if (!userId) return;

    if (myGroups.value.length === 0) {
      await fetchUserGroups();
    } else {
      await fetchPredictions();
    }
  },
  { immediate: true },
);

watch(
  () => activeGroupId.value,
  async (groupId) => {
    if (groupId) {
      await fetchPredictions();
    }
  },
);

onMounted(async () => {
  try {
    await tournamentStore.fetchMatches();
    tournamentStore.subscribeToMatches();
  } finally {
    loadingMatches.value = false;
  }
});

onUnmounted(() => {
  tournamentStore.unsubscribeFromMatches();
});
</script>

<style scoped>
.max-width-card {
  max-width: 500px;
}
.max-width-footer {
  max-width: 960px;
}
.border {
  border: 1px solid rgba(0, 0, 0, 0.12);
}
.mx-auto {
  margin-left: auto;
  margin-right: auto;
}
</style>
