<template>
  <q-card flat bordered class="rounded-borders q-pa-sm bg-white">
    <div class="row items-center justify-around q-col-gutter-sm">
      <!-- Home Side -->
      <div class="col-4 text-center">
        <q-img
          v-if="match.expand?.home_team?.flag_url"
          :src="match.expand.home_team.flag_url"
          style="
            width: 48px;
            height: 32px;
            border-radius: 4px;
            border: 1px solid rgba(0, 0, 0, 0.1);
          "
        />
        <!-- Dynamic Contrast Text Color in Dark Mode -->
        <div
          class="text-subtitle2 text-weight-bold q-mt-xs text-ellipsis"
          :class="$q.dark.isActive ? 'text-white' : 'text-grey-9'"
        >
          {{ match.expand?.home_team?.name }}
        </div>
      </div>

      <!-- Score Inputs -->
      <div class="col-4 text-center">
        <div v-if="isLocked" class="row items-center justify-center q-gutter-x-xs no-wrap">
          <span class="text-h5 text-weight-bold text-grey-7">{{ form.predicted_home }}</span>
          <span class="text-grey-5">:</span>
          <span class="text-h5 text-weight-bold text-grey-7">{{ form.predicted_away }}</span>
        </div>
        <div v-else class="row items-center justify-center q-gutter-y-xs q-gutter-x-none">
          <!-- no-wrap prevents horizontal stack breaking under all mobile widths -->
          <div class="row items-center justify-center no-wrap q-gutter-x-xs">
            <q-input
              v-model.number="form.predicted_home"
              type="number"
              min="0"
              inputmode="numeric"
              pattern="[0-9]*"
              outlined
              dense
              input-class="text-center text-weight-bold text-subtitle1 q-px-xs"
              style="width: 65px"
              hide-bottom-space
              @keydown="blockInvalidChars"
              @update:model-value="emitFormChange"
            />
            <span class="text-grey-7 text-weight-bold">:</span>
            <q-input
              v-model.number="form.predicted_away"
              type="number"
              min="0"
              inputmode="numeric"
              pattern="[0-9]*"
              outlined
              dense
              input-class="text-center text-weight-bold text-subtitle1 q-px-xs"
              style="width: 65px"
              hide-bottom-space
              @keydown="blockInvalidChars"
              @update:model-value="emitFormChange"
            />
          </div>
        </div>
      </div>

      <!-- Away Side -->
      <div class="col-4 text-center">
        <q-img
          v-if="match.expand?.away_team?.flag_url"
          :src="match.expand.away_team.flag_url"
          style="
            width: 48px;
            height: 32px;
            border-radius: 4px;
            border: 1px solid rgba(0, 0, 0, 0.1);
          "
        />
        <!-- Dynamic Contrast Text Color in Dark Mode -->
        <div
          class="text-subtitle2 text-weight-bold q-mt-xs text-ellipsis"
          :class="$q.dark.isActive ? 'text-white' : 'text-grey-9'"
        >
          {{ match.expand?.away_team?.name }}
        </div>
      </div>
    </div>

    <q-separator class="q-my-md" />

    <!-- Details footer row inside match card -->
    <div class="row items-center justify-between">
      <div>
        <q-chip v-if="isLocked" dense color="grey-3" text-color="grey-8" icon="lock">
          Predictions Locked
        </q-chip>
        <q-chip v-else dense color="green-1" text-color="green-8" icon="schedule"> Open </q-chip>

        <!-- High-contrast dynamic points awarded chip -->
        <q-chip
          v-if="match.status === 'finished'"
          dense
          :color="pointsAwarded > 0 ? 'amber-5' : 'grey-7'"
          :text-color="pointsAwarded > 0 ? 'black' : 'white'"
          icon="emoji_events"
          class="text-weight-bold q-px-sm"
        >
          {{ pointsAwarded }} Pts
        </q-chip>
      </div>
    </div>
  </q-card>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import type { Match, Prediction } from '@/types';

const props = defineProps<{
  match: Match;
  groupId: string;
  prediction: Prediction | undefined;
}>();

const emit = defineEmits<{
  (e: 'update', value: { predicted_home: number; predicted_away: number }): void;
}>();

const predictionId = ref<string | null>(null);
const pointsAwarded = ref(0);

const form = ref({
  predicted_home: 0,
  predicted_away: 0,
});

const isLocked = computed(() => {
  if (props.match.status === 'finished') return true;
  if (!props.match.kickoff) return false;
  const isoString = props.match.kickoff.replace(' ', 'T');
  return new Date(isoString).getTime() <= Date.now();
});

// Intercept keys to strictly allow digits only and prevent exponential expressions
function blockInvalidChars(event: KeyboardEvent) {
  const allowedKeys = ['Backspace', 'Tab', 'ArrowLeft', 'ArrowRight', 'Delete', 'Enter'];
  if (allowedKeys.includes(event.key)) {
    return;
  }
  if (!/^[0-9]$/.test(event.key)) {
    event.preventDefault();
  }
}

// Reactively load and assign incoming predictions directly from parent cache
watch(
  () => props.prediction,
  (newPred) => {
    if (newPred) {
      predictionId.value = newPred.id;
      form.value.predicted_home = newPred.predicted_home;
      form.value.predicted_away = newPred.predicted_away;
      pointsAwarded.value = newPred.points_awarded || 0;
    } else {
      predictionId.value = null;
      form.value.predicted_home = 0;
      form.value.predicted_away = 0;
      pointsAwarded.value = 0;
    }
    emitFormChange();
  },
  { immediate: true },
);

function emitFormChange() {
  emit('update', {
    predicted_home: form.value.predicted_home,
    predicted_away: form.value.predicted_away,
  });
}
</script>

<style scoped>
.text-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
