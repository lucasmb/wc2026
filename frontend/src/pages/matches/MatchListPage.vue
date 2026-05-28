<template>
  <q-page class="q-pa-md">
    <!-- Header Block -->
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h6 text-weight-bold text-primary">Match Predictions</div>

      <!-- Selectors Row -->
      <div class="row q-gutter-x-sm">
        <!-- Group Stage Selector -->
        <q-select
          v-model="selectedGroupStage"
          :options="groupStageOptions"
          label="Stage/Group"
          outlined
          dense
          emit-value
          map-options
          style="width: 130px"
        />

        <!-- User's Prediction Groups Dropdown -->
        <q-select
          v-if="myGroups.length > 0"
          v-model="activeGroupId"
          :options="myGroups"
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

    <!-- No Groups Warning Call-to-Action -->
    <div v-if="!loadingGroups && myGroups.length === 0" class="text-center q-my-xl">
      <q-card flat bordered class="q-pa-xl rounded-borders bg-white max-width-card mx-auto">
        <q-icon name="group" size="64px" color="grey-4" />
        <div class="text-h6 text-weight-bold text-primary q-mt-md">Join a Group to Predict</div>
        <div class="text-subtitle2 text-grey-6 q-mt-xs">
          Predictions are recorded within private/public prediction groups. Create or join a group
          to start!
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

    <!-- Group Phase Match Loading Indicators -->
    <div
      v-else-if="loadingMatches || loadingGroups || loadingPredictions"
      class="row justify-center q-my-xl"
    >
      <q-spinner-dots color="primary" size="40px" />
    </div>

    <!-- Render Match Cards Grid -->
    <div v-else class="row q-col-gutter-md">
      <div v-if="filteredMatches.length === 0" class="col-12 text-center q-my-xl text-grey-6">
        No matches configured for this selection
      </div>
      <div v-for="match in filteredMatches" :key="match.id" class="col-12 col-md-6">
        <div class="text-caption text-grey-6 q-mb-xs">
          Match #{{ match.match_number }} • {{ formatLocalTime(match.kickoff) }} • {{ match.venue }}
        </div>
        <!-- Pass the matched prediction down reactively and listen for changes -->
        <PredictionForm
          v-if="activeGroupId"
          :match="match"
          :groupId="activeGroupId"
          :prediction="predictionsMap[match.id]"
          @saved="handlePredictionSaved"
        />
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { pb } from '@/boot/pocketbase';
import { useAuthStore } from '@/stores/auth';
import { useTournamentStore } from '@/stores/tournament';
import PredictionForm from '@/components/match/PredictionForm.vue';
import type { PredictionGroup, Prediction } from '@/types';

const tournamentStore = useTournamentStore();
const authStore = useAuthStore();

const loadingMatches = ref(true);
const loadingGroups = ref(true);
const loadingPredictions = ref(false);

const selectedGroupStage = ref('A');
const activeGroupId = ref<string>('');
const myGroups = ref<PredictionGroup[]>([]);

// Mapping dictionary to cache and distribute loaded predictions instantly
const predictionsMap = ref<Record<string, Prediction>>({});

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
];

const filteredMatches = computed(() => {
  return tournamentStore.matches.filter(
    (m) => m.phase === 'group' && m.group_code === selectedGroupStage.value,
  );
});

function formatLocalTime(isoString: string) {
  return new Intl.DateTimeFormat(navigator.language, {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(new Date(isoString));
}

// 1 Batch fetch all predictions the user has placed in this active group
async function fetchPredictions() {
  if (!authStore.user?.id || !activeGroupId.value) return;
  loadingPredictions.value = true;
  try {
    const list = await pb.collection('predictions_id').getFullList({
      filter: `user = "${authStore.user.id}" && prediction_group = "${activeGroupId.value}"`,
    });

    const map: Record<string, Prediction> = {};
    list.forEach((p) => {
      map[p.match] = p as unknown as Prediction;
    });
    predictionsMap.value = map;
  } catch (err: unknown) {
    console.error('Failed fetching predictions batch:', err);
  } finally {
    loadingPredictions.value = false;
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

// Keep the local batch-cache updated instantly on any child writes
function handlePredictionSaved(updatedPrediction: Prediction) {
  predictionsMap.value[updatedPrediction.match] = updatedPrediction;
}

// Trigger batch fetch whenever the active group or user changes
watch(
  [() => activeGroupId.value, () => authStore.user?.id],
  async ([groupId, userId]) => {
    if (groupId && userId) {
      await fetchPredictions();
    }
  },
  { immediate: true },
);

watch(
  () => authStore.user?.id,
  async (newUserId) => {
    if (newUserId) {
      await fetchUserGroups();
    }
  },
  { immediate: true },
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
.mx-auto {
  margin-left: auto;
  margin-right: auto;
}
</style>
